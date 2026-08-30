package ui

import (
	"fmt"
	"lazypacket/internal/layers"
	"time"
)

func formatPacket(i int, packet Entry, elapsed time.Duration) string {
	ip, ok := packet.Packet.NetworkLayer().(*layers.IPv4)
	if !ok {
		return fmt.Sprintf("#%d\t%.6f\t unknown", i, elapsed.Seconds())
	}
	srcAddr := ip.SrcAddr.String()
	dstAddr := ip.DstAddr.String()
	protocol := ip.Protocol.String()

	return fmt.Sprintf("#%-6d %-10.5f %-16s %-16s %-6s", i, elapsed.Seconds(), srcAddr, dstAddr, protocol)
}
