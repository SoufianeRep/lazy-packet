package layers

type EthernetType uint16

const (
	EthernetTypeIPv4 EthernetType = 0x0800
	EthernetTypeIPv6 EthernetType = 0x86dd
	EthernetTypeARP  EthernetType = 0x0806
)

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
