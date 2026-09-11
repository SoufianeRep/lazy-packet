package layers

import "fmt"

type EthernetType uint16

const (
	EthernetTypeIPv4 EthernetType = 0x0800
	EthernetTypeIPv6 EthernetType = 0x86dd
	EthernetTypeARP  EthernetType = 0x0806
)

func (et EthernetType) String() string {
	switch et {
	case EthernetTypeIPv4:
		return fmt.Sprintf("0x%04x (IPv4)", uint16(et))
	case EthernetTypeIPv6:
		return fmt.Sprintf("0x%04x (IPv6", uint16(et))
	case EthernetTypeARP:
		return fmt.Sprintf("0x%04x (ARP)", uint16(et))
	default:
		return fmt.Sprintf("0x%04x", uint16(et))
	}
}

type IPProtocol uint8

const (
	IPProtocolTCP IPProtocol = 6
	IPProtocolUDP IPProtocol = 17
)

func (p IPProtocol) String() string {
	switch p {
	case IPProtocolUDP:
		return "UDP"
	case IPProtocolTCP:
		return "TCP"
	default:
		return "Unknown"
	}
}
