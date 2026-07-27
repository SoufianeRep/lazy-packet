package layers

type EthernetType uint16

const (
	EthernetTypeIPv4 EthernetType = 0x0800
	EthernetTypeIPv6 EthernetType = 0x86dd
)

type IPProtocol uint8

const (
	IPProtocolTCP IPProtocol = 6
	IPProtocolUDP IPProtocol = 17
)
