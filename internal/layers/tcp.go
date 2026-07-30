package layers

import lp "lazypacket"

type TCP struct {
	lp.BaseLayer
}

func (t *TCP) LayerType() lp.LayerType {
	return lp.LayerTypeTCP
}

func (t *TCP) DecodeFromBytes(data []byte, p lp.PacketBuilder) error {
	return nil
}
