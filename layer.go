package lazypacket

type Layer interface {
	// LayerType says what kind of layer this is.
	LayerType() LayerType
	// LayerContents is the raw bytes that belong to *this* layer's header
	// (e.g. just the 14 Ethernet header bytes, not what follows).
	LayerContents() []byte
	// LayerPayload is everything left over after this layer's header —
	// i.e. exactly the bytes handed to the next layer's decoder.
	LayerPayload() []byte
}

type BaseLayer struct {
	Contents []byte
	Payload  []byte
}

// LayerContents returns the bytes of the packet layer.
func (b *BaseLayer) LayerContents() []byte {
	return b.Contents
}

// LayerPayload returns the bytes contained within the packet layer.
func (b *BaseLayer) LayerPayload() []byte {
	return b.Payload
}
