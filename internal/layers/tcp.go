package layers

import (
	"encoding/binary"
	"fmt"
	"strconv"

	lp "lazypacket"
)

//		0                   1                   2                   3
//		0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
//		+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//		|          Source Port          |       Destination Port        |
//		+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//		|                        Sequence Number                        |
//		+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//		|                    Acknowledgment Number                      |
//		+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//		|  Data |       |C|E|U|A|P|R|S|F|                               |
//		| Offset| Rsrvd |W|C|R|C|S|S|Y|I|            Window             |
//		|       |       |R|E|G|K|H|T|N|N|                               |
//		+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//		|           Checksum            |         Urgent Pointer        |
//		+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//		|                           [Options]                           |
//		+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//		|                                                               :
//		:                             Data                              :
//		:                                                               |
//		+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

type TCPPort uint16

var tcpPortNames = map[TCPPort]string{
	20:   "ftp-data",
	21:   "ftp",
	22:   "ssh",
	23:   "telnet",
	25:   "smtp",
	53:   "dns",
	80:   "http",
	110:  "pop3",
	143:  "imap",
	443:  "https",
	3306: "mysql",
	5432: "postgresql",
}

func (tp TCPPort) String() string {
	if name, ok := tcpPortNames[tp]; ok {
		return fmt.Sprintf("%d (%s)", uint16(tp), name)
	}
	return strconv.Itoa(int(tp))
}

type TCPOptionKind uint8

const (
	TCPOptionKindEndList       TCPOptionKind = 0
	TCPOptionKindNop           TCPOptionKind = 1
	TCPOptionKindMSS           TCPOptionKind = 2
	TCPOptionKindWindowScale   TCPOptionKind = 3
	TCPOptionKindSACKPermitted TCPOptionKind = 4
	TCPOptionKindSACK          TCPOptionKind = 5
	TCPOptionKindTimestamp     TCPOptionKind = 8
)

func (tok TCPOptionKind) String() string {
	switch tok {
	case TCPOptionKindEndList:
		return "EndList"
	case TCPOptionKindNop:
		return "NOP"
	case TCPOptionKindMSS:
		return "MSS"
	case TCPOptionKindWindowScale:
		return "WindowScale"
	case TCPOptionKindSACKPermitted:
		return "SACKPermitted"
	case TCPOptionKindSACK:
		return "SACK"
	case TCPOptionKindTimestamp:
		return "Timestamp"
	default:
		return "Unsupported"
	}
}

type TCPOption struct {
	OptionType   TCPOptionKind
	OptionLength uint8
	OptionData   []byte
}

type TCP struct {
	lp.BaseLayer
	SrcPort, DstPort                       TCPPort
	SeqNum                                 uint32
	AckNum                                 uint32
	DataOffset                             uint8
	CWR, ECE, URG, ACK, PSH, RST, SYN, FIN bool
	Window                                 uint16
	Checksum                               uint16
	UrgentPointer                          uint16
	Options                                []TCPOption
	Padding                                []byte
}

func (t *TCP) LayerType() lp.LayerType {
	return lp.LayerTypeTCP
}

func (t *TCP) DecodeFromBytes(data []byte, p lp.PacketBuilder) error {
	if len(data) < 20 {
		p.SetTruncated()
		return fmt.Errorf("[TCP] tcp header too small, only %d bytes minimum 20 bytes", len(data))
	}

	t.SrcPort = TCPPort(binary.BigEndian.Uint16(data[0:2]))
	t.DstPort = TCPPort(binary.BigEndian.Uint16(data[2:4]))
	t.SeqNum = binary.BigEndian.Uint32(data[4:8])
	t.AckNum = binary.BigEndian.Uint32(data[8:12])
	t.DataOffset = data[12] >> 4

	t.CWR = data[13]&0x80 != 0
	t.ECE = data[13]&0x40 != 0
	t.URG = data[13]&0x20 != 0
	t.ACK = data[13]&0x10 != 0
	t.PSH = data[13]&0x08 != 0
	t.RST = data[13]&0x04 != 0
	t.SYN = data[13]&0x02 != 0
	t.FIN = data[13]&0x01 != 0

	t.Window = binary.BigEndian.Uint16(data[14:16])
	t.Checksum = binary.BigEndian.Uint16(data[16:18])
	t.UrgentPointer = binary.BigEndian.Uint16(data[18:20])

	if t.DataOffset < 5 {
		return fmt.Errorf("[TCP Error] invalid TCP DOffset %d < 5", t.DataOffset)
	}

	dataStart := int(t.DataOffset) * 4
	if dataStart > len(data) {
		p.SetTruncated()
		t.Contents = data
		t.Payload = nil
		return fmt.Errorf("[TCP Error] truncated header: need %d bytes, got %d", dataStart, len(data))
	}

	t.Contents = data[:dataStart]
	t.Payload = data[dataStart:]

	optionsData := data[20:dataStart]
OPTIONS:
	for len(optionsData) > 0 {
		t.Options = append(t.Options, TCPOption{OptionType: TCPOptionKind(optionsData[0])})
		opt := &t.Options[len(t.Options)-1] // get the latest option
		switch opt.OptionType {
		case TCPOptionKindEndList:
			opt.OptionLength = 1
			t.Padding = optionsData[1:]
			break OPTIONS
		case TCPOptionKindNop:
			opt.OptionLength = 1
		default:
			if len(optionsData) < 2 {
				return fmt.Errorf("[TCP Error] truncated options, missing length byte")
			}

			length := optionsData[1]
			if length < 2 || int(length) > len(optionsData) {
				return fmt.Errorf("[TCP Error] invalid option length %d for kind %d", length, opt.OptionType)
			}

			opt.OptionLength = length
			opt.OptionData = optionsData[2:length]
		}

		optionsData = optionsData[opt.OptionLength:] // truncate to move to next
		// this is incomplete but should be sufficient for now
	}

	return nil
}
