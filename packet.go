package lazypacket

type Packet interface {
	// String returns a human-readable string representation of the packet.
	String() string

	Layers() []Layer

	Layer(LayerType) Layer
}
