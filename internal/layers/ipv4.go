package layers

import (
	"encoding/binary"
	"fmt"
	lp "lazypacket"
	"net"
	"strings"
)

type IPv4Flag uint8

const (
	IPv4EvilBit       IPv4Flag = 1 << 2 // https://datatracker.ietf.org/doc/html/rfc3514
	IPv4DontFragment  IPv4Flag = 1 << 1
	IPv4MoreFragments IPv4Flag = 1 << 0
)

func (ipf IPv4Flag) String() string {
	var s []string
	if ipf&IPv4EvilBit != 0 {
		s = append(s, "Evil")
	}

	if ipf&IPv4DontFragment != 0 {
		s = append(s, "DF") // don't fragment else fragment
	}

	if ipf&IPv4MoreFragments != 0 {
		s = append(s, "MF") // More fragments else Last fragment
	}

	return strings.Join(s, "|")
}

// IPv4 Specification lives here https://datatracker.ietf.org/doc/html/rfc791
// Check for further details
type IPv4 struct {
	lp.BaseLayer
	Version    uint8      // Version (4bits)
	IHL        uint8      // Internet Header Length (4bits)
	TOS        uint8      // Type Of Service (8bits)
	Length     uint16     // Total length of the datagram (16bits)
	Id         uint16     // Identification (16bits)
	Flags      IPv4Flag   // Various control flags (3bits)
	FragOffset uint16     // Fragment offset (13bits)
	TTL        uint8      // Time To Live (8bits)
	Protocol   IPProtocol // Protocol (8bits)
	Checksum   uint16     // Header Checksum (16bits)
	SrcAddr    net.IP     // Source Address (32bits)
	DstAddr    net.IP     // Destination Address (32bits)
	Options    []IPv4Option
	Padding    []byte
}

func (ip *IPv4) LayerType() lp.LayerType {
	return lp.LayerTypeIPV4
}

type IPv4Option struct {
	// 1 bit   copied flag indicates that this option is copied into all fragments on fragmentation
	// 2 bits  option class: 0 control, 1 reserver for future use, 2 debugging and measurement, 3 reserved for future use
	// 5 bits  option number
	OptType   uint8
	OptLength uint8
	OptData   []byte
}

func (ipo *IPv4Option) String() string {
	return fmt.Sprintf("IPv4Option(%v:%v", ipo.OptType, ipo.OptData)
}

func (ip *IPv4) DecodeFromBytes(data []byte, p lp.PacketBuilder) error {
	if len(data) < 20 {
		p.SetTruncated()
		return fmt.Errorf("[IP4] invalid IP header of %d less than the expected 20", len(data))
	}
	ip.Version = data[0] >> 4
	ip.IHL = data[0] & 0x0F
	ip.TOS = data[1]
	ip.Length = binary.BigEndian.Uint16(data[2:4]) // uint16(data[2])<<8 | uint16(data[3])
	ip.Id = binary.BigEndian.Uint16(data[4:6])     // uint16(data[4])<<8 | uint16(data[5])

	fragFlags := binary.BigEndian.Uint16(data[6:8])
	ip.Flags = IPv4Flag(fragFlags >> 13)
	ip.FragOffset = fragFlags & 0x1FFF

	ip.TTL = data[8]
	ip.Protocol = IPProtocol(data[9])

	ip.Checksum = binary.BigEndian.Uint16(data[10:12]) // uint16(data[10])<<8 | uint16(data[11])

	ip.SrcAddr = data[12:16]
	ip.DstAddr = data[16:20]

	opts, padding, err := decodeIPv4Options(data[20 : ip.IHL*4])
	if err != nil {
		p.SetTruncated()
		return err
	}
	ip.Options = opts
	ip.Padding = padding

	ip.Contents = data[:ip.IHL*4]
	ip.Payload = data[ip.IHL*4:]

	p.AddLayer(ip)

	switch ip.Protocol {
	case IPProtocolTCP:
		return p.NextDecoder(&TCP{}, ip.Payload)
	case IPProtocolUDP:
		return p.NextDecoder(&UDP{}, ip.Payload)
	default:
		return nil
	}
}

func decodeIPv4Options(data []byte) (opts []IPv4Option, padding []byte, err error) {
	for len(data) > 0 {
		if opts == nil {
			opts = make([]IPv4Option, 0, 4)
		}
		opt := IPv4Option{OptType: data[0]}
		switch opt.OptType {
		case 0: // End of options
			opt.OptLength = 1
			opts = append(opts, opt)
			return opts, data[1:], nil
		case 1: // NOP
			opt.OptLength = 1
			data = data[1:]
			opts = append(opts, opt)
		default:
			if len(data) < 2 {
				return nil, nil, fmt.Errorf("[IP4] invalid IPv4 option length: %d bytes left", len(data))
			}
			opt.OptLength = data[1]
			if len(data) < int(opt.OptLength) {
				return nil, nil, fmt.Errorf("[IP4] option type %d length %d exceeds remaining header", opt.OptType, opt.OptLength)
			}
			if opt.OptLength < 2 {
				return nil, nil, fmt.Errorf("[IP4] invalid IPv4 option type %d length %d", opt.OptType, opt.OptLength)
			}
			opt.OptData = data[2:opt.OptLength]
			data = data[opt.OptLength:]
			opts = append(opts, opt)
		}
	}
	return opts, nil, nil
}
