package layers

import (
	"encoding/binary"
	"errors"
	"net"

	lp "lazypacket"
)

const hdrLen = 14

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
	lp.BaseLayer

	DstMAC    net.HardwareAddr // destination MAC address
	SrcMAC    net.HardwareAddr // source MAC address
	VLAN      *VLANTag         // optional 802.1Q tag; nil if the frame is untagged
	EtherType EthernetType     // payload protocol, e.g. 0x0800 for IPv4 (read after any VLAN tag)
}

func (eth *Ethernet) LayerType() lp.LayerType {
	return lp.LayerTypeEthernet
}

func (eth *Ethernet) DecodeFromBytes(data []byte, p lp.PacketBuilder) error {
	if len(data) < hdrLen {
		return errors.New("ethernet frame too small")
	}

	eth.DstMAC = net.HardwareAddr(data[0:6])
	eth.SrcMAC = net.HardwareAddr(data[6:12])

	//uint16(data[12])<<8 | uint16(data[13])
	eth.EtherType = EthernetType(binary.BigEndian.Uint16(data[12:14]))
	if eth.EtherType < 0x0600 {
		// This is a length field, not an EtherType.  We need to check for a VLAN tag.
		// TODO: Add support for this later
	}
	eth.Contents = data[:hdrLen]
	eth.Payload = data[hdrLen:]

	p.AddLayer(eth)
	p.SetLinkLayer(eth)

	switch eth.EtherType {
	case EthernetTypeIPv4:
		return p.NextDecoder(&IPv4{}, eth.Payload)
	case EthernetTypeIPv6:
		// TODO: p.NextDecoder() with the correct decoder
	case EthernetTypeARP:
		// TODO: p.NextDecoder() with the correct decoder
	}
	return nil
}
