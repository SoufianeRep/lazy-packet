package lazypacket

type Packet struct {
	data []byte

	layers []Layer

	link      Layer
	network   Layer
	transport Layer

	truncated bool
	err       error
}

func (p *Packet) AddLayer(l Layer) {
	p.layers = append(p.layers, l)
}

func (p *Packet) SetTruncated() {
	p.truncated = true
}

func (p *Packet) SetLinkLayer(l Layer) {
	if p.link == nil {
		p.link = l
	}
}

func (p *Packet) SetNetworkLayer(l Layer) {
	if p.network == nil {
		p.network = l
	}
}

func (p *Packet) SetTransportLayer(l Layer) {
	if p.transport == nil {
		p.transport = l
	}
}

func (p *Packet) NextDecoder(next Decoder, data []byte) error {
	return next.DecodeFromBytes(data, p)
}
