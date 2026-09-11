package ui

import (
	"fmt"
	lp "lazypacket"
	"lazypacket/internal/layers"
	"strings"
	"time"

	"charm.land/lipgloss/v2"
)

type Summarizer interface {
	Summary() string
}

type Detailer interface {
	Detail() string
}

var selectedRowStyle = lipgloss.NewStyle().Reverse(true)

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

func formatDetail(e Entry) string {
	var blocks []string
	for _, l := range []lp.Layer{e.Packet.LinkLayer(), e.Packet.NetworkLayer(), e.Packet.TransportLayer()} {
		if l == nil {
			continue
		}

		if d, ok := l.(Detailer); ok {
			blocks = append(blocks, d.Detail())
		}
	}

	return strings.Join(blocks, "\n\n")
}
