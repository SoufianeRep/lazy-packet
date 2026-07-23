package protocol

type BaseLayer struct {
	Contents []byte // The bytes that make up this layer
	Payload  []byte // The bytes contained by (but not part of) this layer
}

// LayerContents returns the bytes of the packet layer.
func (b *BaseLayer) LayerContents() []byte {
	return b.Contents
}

// LayerPayload returns the bytes contained within the packet layer.
func (b *BaseLayer) LayerPayload() []byte {
	return b.Payload
}
