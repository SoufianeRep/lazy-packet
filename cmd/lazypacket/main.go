package main

import (
	"fmt"
	lp "lazypacket"
	"lazypacket/internal/capture"
	"lazypacket/internal/layers"
	"lazypacket/internal/ui"
	"os"

	tea "charm.land/bubbletea/v2"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: lazypacket <interface>")
		os.Exit(1)
	}

	h, err := capture.Open(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "open failed: %v\n", err)
		os.Exit(1)
	}
	defer h.Close()

	p := tea.NewProgram(ui.InitModel())
	go func() {
		for {
			data, ts, err := h.ReadFrame()
			if err != nil {
				fmt.Fprintf(os.Stderr, "read failed: %v\n", err)
				break
			}
			packet := lp.NewPacket(data, &layers.Ethernet{})

			p.Send(ui.FrameMsg{Packet: packet, TimeStamp: ts})
		}
	}()

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
