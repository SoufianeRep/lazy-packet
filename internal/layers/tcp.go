package layers

import (
	"encoding/binary"
	"fmt"
	lp "lazypacket"
)

type TCP struct {
	lp.BaseLayer
	SrcPort uint16
	DstPort uint16
	SeqNum  uint32
}

func (t *TCP) LayerType() lp.LayerType {
	return lp.LayerTypeTCP
}

func (t *TCP) DecodeFromBytes(data []byte, p lp.PacketBuilder) error {
	if len(data) < 20 {
		p.SetTruncated()
		return fmt.Errorf("[TCP] tcp header too small, only %d bytes minimum 20 bytes", len(data))
	}

	t.SrcPort = binary.BigEndian.Uint16(data[0:2])
	t.DstPort = binary.BigEndian.Uint16(data[2:4])
	t.SeqNum = binary.BigEndian.Uint32(data[4:8])

	return nil
}
