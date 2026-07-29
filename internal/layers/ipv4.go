package layers

import (
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
	Version    uint8    // Version (4bits)
	IHL        uint8    // Internet Header Length (4bits)
	TOS        uint8    // Type Of Service (8bits)
	Length     uint16   // Total length of the datagram (16bits)
	Id         uint16   // Identification (16bits)
	Flags      IPv4Flag // Various control flags (3bits)
	FragOffset uint16   // Fragment offset (13bits)
	TTL        uint8    // Time To Live (8bits)
	Protocol   uint8    // Protocol (8bits)
	Checksum   uint16   // Header Checksum (16bits)
	SrcAddr    net.IP   // Source Address (32bits)
	DstAddr    net.IP   // Destination Address (32bits)
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
	return nil
}
