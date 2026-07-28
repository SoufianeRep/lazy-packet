package lazypacket

type Decoder interface {
	DecodeFromBytes(data []byte, p PacketBuilder) error
}

type PacketBuilder interface {
	AddLayer(l Layer)

	SetTruncated()

	SetLinkLayer(l Layer)

	SetNetworkLayer(l Layer)

	SetTransportLayer(l Layer)

	NextDecoder(next Decoder, data []byte) error
}
