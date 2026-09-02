package ui

import (
	"fmt"
	"lazypacket/internal/layers"
	"time"
)

type Summarizer interface {
	Summary() string
}

func formatPacket(i int, packet Entry, elapsed time.Duration) string {
	ip, ok := packet.Packet.NetworkLayer().(*layers.IPv4)
	if !ok {
		return fmt.Sprintf("#%d\t%.6f\t unknown", i, elapsed.Seconds())
	}

	fixed := fmt.Sprintf("#%-6d %-10.5f %-16s %-16s %-6s", i, elapsed.Seconds(), ip.SrcAddr.String(), ip.DstAddr.String(), ip.Protocol.String())

	info := ""
	if s, ok := packet.Packet.TransportLayer().(Summarizer); ok {
		info = s.Summary()
	}

	return fixed + " " + info
}
