package lazypacket

type Decoder interface {
	Decode(data []byte, p PacketBuilder) error
}

type PacketBuilder interface {
	AddLayer(l Layer)

	SetTruncated()

	SetLinkLayer(l Layer)

	SetNetworkLayer(l Layer)

	setTransportLayer(l Layer)

	NextDecoder(next Decoder, data []byte) error
}
