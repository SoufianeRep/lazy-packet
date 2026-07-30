package layers

import lp "lazypacket"

type UDP struct {
	lp.BaseLayer
}

func (u *UDP) LayerType() lp.LayerType {
	return lp.LayerTypeUDP
}

func (u *UDP) DecodeFromBytes(data []byte, p lp.PacketBuilder) error {
	return nil
}
