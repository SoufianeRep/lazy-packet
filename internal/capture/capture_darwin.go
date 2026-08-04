package capture

import (
	"fmt"
	"net"
	"syscall"
	"time"
	"unsafe"
)

type Handle struct {
	fd    int
	iface *net.Interface
	buf   []byte
	rest  []byte
}

// ifreq mirrors struct ifreq from <net/if.h>: a 16-byte interface name
// followed by a 16-byte union. BIOCSETIF only reads the name, so the
// union half is left zeroed.
type ifreq struct {
	Name [syscall.IFNAMSIZ]byte
	_    [16]byte
}

func Open(ifaceName string) (*Handle, error) {
	iface, err := net.InterfaceByName(ifaceName)
	if err != nil {
		return nil, fmt.Errorf("interface look up failed: %w", err)
	}

	fd, err := openBPFDevice()
	if err != nil {
		return nil, err
	}

	var req ifreq
	copy(req.Name[:], ifaceName)
	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), syscall.BIOCSETIF, uintptr(unsafe.Pointer(&req))); errno != 0 {
		syscall.Close(fd)
		return nil, fmt.Errorf("'BIOCSETIF' failed: %w", errno)
	}

	// deliver each packet as soon as it arrives instead of waiting for the buffer to fill
	immediate := uint32(1)
	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), syscall.BIOCIMMEDIATE, uintptr(unsafe.Pointer(&immediate))); errno != 0 {
		syscall.Close(fd)
		return nil, fmt.Errorf("'BIOCIMMEDIATE' failed: %w", errno)
	}

	var bufLen uint32
	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), syscall.BIOCGBLEN, uintptr(unsafe.Pointer(&bufLen))); errno != 0 {
		syscall.Close(fd)
		return nil, fmt.Errorf("'BIOCGBLEN' failed: %w", errno)
	}

	return &Handle{fd: fd, iface: iface, buf: make([]byte, bufLen)}, nil
}

// openBPFDevice finds the first /dev/bpfN device not already claimed by another process.
// BPF devices are exclusive: only one open file descriptor per device at a time.
func openBPFDevice() (int, error) {
	for i := 0; i < 256; i++ {
		path := fmt.Sprintf("/dev/bpf%d", i)
		fd, err := syscall.Open(path, syscall.O_RDWR, 0)
		if err == nil {
			return fd, nil
		}
		if err != syscall.EBUSY {
			return -1, fmt.Errorf("open %s failed: %w", path, err)
		}
	}
	return -1, fmt.Errorf("no free /dev/bpf* device found")
}

func (h *Handle) Close() error {
	return syscall.Close(h.fd)
}

func bpfWordAlign(n int) int {
	const align = syscall.BPF_ALIGNMENT
	return (n + align - 1) &^ (align - 1)
}

func (h *Handle) ReadFrame() ([]byte, time.Time, error) {
	if len(h.rest) == 0 {
		n, err := syscall.Read(h.fd, h.buf)
		if err != nil {
			return nil, time.Time{}, fmt.Errorf("read frame: %w", err)
		}
		h.rest = h.buf[:n]
	}

	if len(h.rest) < syscall.SizeofBpfHdr {
		h.rest = nil
		return nil, time.Time{}, fmt.Errorf("read frame: short bpf header")
	}

	hdr := (*syscall.BpfHdr)(unsafe.Pointer(&h.rest[0]))
	hdrLen := int(hdr.Hdrlen)
	capLen := int(hdr.Caplen)

	if hdrLen+capLen > len(h.rest) {
		h.rest = nil
		return nil, time.Time{}, fmt.Errorf("read frame: truncated bpf packet")
	}

	data := make([]byte, capLen)
	copy(data, h.rest[hdrLen:hdrLen+capLen])
	ts := time.Unix(int64(hdr.Tstamp.Sec), int64(hdr.Tstamp.Usec)*1000)

	advance := bpfWordAlign(hdrLen + capLen)
	if advance >= len(h.rest) {
		h.rest = nil
	} else {
		h.rest = h.rest[advance:]
	}

	return data, ts, nil
}
