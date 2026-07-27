package lazypacket

type LayerType int64

const (
	LayerTypeUnknow LayerType = iota
	LayerTypeEthernet
	LayerTypeIPV4
	LayerTypeIPv6
	LayerTypeTCP
	LayerTypeUDP
)
