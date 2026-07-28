package capture

import (
	"fmt"
	"net"
	"syscall"
	"time"
)

type Handle struct {
	fd    int
	iface *net.Interface
}

func htons(i uint16) uint16 {
	return (i<<8)&0xff00 | i>>8
}

func Open(ifaceName string) (*Handle, error) {
	iface, err := net.InterfaceByName(ifaceName)
	if err != nil {
		return nil, fmt.Errorf("interface look up failed: %w", err)
	}

	fd, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW, int(htons(syscall.ETH_P_ALL)))
	if err != nil {
		return nil, fmt.Errorf("'Socket()' failed: %w", err)
	}

	addr := syscall.SockaddrLinklayer{
		Protocol: htons(syscall.ETH_P_ALL),
		Ifindex:  iface.Index,
	}

	// bind
	err = syscall.Bind(fd, &addr)
	if err != nil {
		syscall.Close(fd)
		return nil, fmt.Errorf("'Bind()' failed: %w", err)
	}

	return &Handle{fd, iface}, nil
}

func (h *Handle) Close() error {
	return syscall.Close(h.fd)
}

type Frame struct {
	Data      []byte
	Timestamp time.Time
}

func (h *Handle) ReadFrame() (Frame, error) {
	buf := make([]byte, 65536)
	n, _, err := syscall.Recvfrom(h.fd, buf, 0)
	if err != nil {
		return Frame{}, fmt.Errorf("read frame: %w", err)
	}

	frame := make([]byte, n)
	copy(frame, buf[:n])

	return Frame{Data: frame, Timestamp: time.Now()}, nil
}
