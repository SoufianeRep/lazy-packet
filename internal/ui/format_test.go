package ui

import (
	"fmt"
	"testing"
	"time"

	lp "lazypacket"
	"lazypacket/internal/layers"
)

// A hand-built 34-byte frame: 14-byte Ethernet header (EtherType 0x0800)
// followed by a 20-byte IPv4 header (IHL=5, no options), protocol TCP.
// Payload is empty — formatPacket never looks past the IPv4 header.
var testFrame = []byte{
	// Ethernet: dst mac, src mac, EtherType
	0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF,
	0x11, 0x22, 0x33, 0x44, 0x55, 0x66,
	0x08, 0x00,
	// IPv4 header
	0x45,       // version 4, IHL 5
	0x00,       // TOS
	0x00, 0x22, // total length = 34
	0x00, 0x00, // id
	0x40, 0x00, // flags (DF) / frag offset
	0x40,       // TTL 64
	0x06,       // protocol: TCP
	0x00, 0x00, // checksum
	0xC0, 0xA8, 0x01, 0x2A, // src 192.168.1.42
	0x5D, 0xB8, 0xD8, 0x22, // dst 93.184.216.34
}

func newTestEntries(n int) []Entry {
	pkt := lp.NewPacket(testFrame, &layers.Ethernet{})

	base := time.Now()
	entries := make([]Entry, n)
	for i := range entries {
		entries[i] = Entry{
			Packet:    pkt,
			TimeStamp: base.Add(time.Duration(i) * time.Millisecond),
		}
	}
	return entries
}

func BenchmarkPacketLines(b *testing.B) {
	for _, n := range []int{100, 1_000, 10_000, 100_000} {
		entries := newTestEntries(n)

		if _, ok := entries[0].Packet.NetworkLayer().(*layers.IPv4); !ok {
			b.Fatalf("fixture frame did not decode an IPv4 layer — check testFrame bytes")
		}

		b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				packetLines(entries)
			}
		})
	}
}

// BenchmarkUpdateFrameMsg mirrors what Update actually does on every
// FrameMsg: format the whole packet list, then hand it to the viewport.
// This is the number that matters for real-world feel, since it includes
// viewport's own internal rescan on top of packetLines' formatting cost.
func BenchmarkUpdateFrameMsg(b *testing.B) {
	for _, n := range []int{100, 1_000, 10_000, 100_000} {
		entries := newTestEntries(n)
		vp := newViewport()

		b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				wasAtBottom := vp.AtBottom()
				vp.SetContent(packetLines(entries))
				if wasAtBottom {
					vp.GotoBottom()
				}
			}
		})
	}
}
