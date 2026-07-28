package layers

import (
	lp "lazypacket"
)

type IPv4 struct {
	lp.BaseLayer
}

func (ip *IPv4) LayerType() lp.LayerType {
	return lp.LayerTypeIPV4
}

func (ip *IPv4) DecodeFromBytes(data []byte, p lp.PacketBuilder) error {
	return nil
}
