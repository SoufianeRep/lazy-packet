package protocol

import (
	"errors"
	"net"
)

type TCI struct {
	PCP uint8  // Priority Code Point (3 bits)
	DEI bool   // Drop Eligible Indicator (1 bit)
	VID uint16 // VLAN Identifier (12 bits)
}
type VLANTag struct {
	TPID uint16 // Tag Protocol Identifier, e.g. 0x8100 for 802.1Q
	TCI         // Tag Control Information
}
type Ethernet struct {
	BaseLayer

	DstMAC    net.HardwareAddr // destination MAC address
	SrcMAC    net.HardwareAddr // source MAC address
	VLAN      *VLANTag         // optional 802.1Q tag; nil if the frame is untagged
	EtherType uint16           // payload protocol, e.g. 0x0800 for IPv4 (read after any VLAN tag)
}

func DecodeFromBytes(data []byte) (*Ethernet, []byte, error) {
	var err error
	if len(data) < 14 {
		err = errors.New("ethernet frame too small")
	}

	eth := &Ethernet{}
	eth.DstMAC = net.HardwareAddr(data[0:6])
	eth.SrcMAC = net.HardwareAddr(data[6:12])
	eth.EtherType = uint16(data[12])<<8 | uint16(data[13])
	if eth.EtherType < 0x0600 {
		// This is a length field, not an EtherType.  We need to check for a VLAN tag.
	}

	eth.BaseLayer = BaseLayer{Contents: data[:14], Payload: data[14:]}

	return eth, data, err
}
