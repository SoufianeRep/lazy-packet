package main

import (
	"fmt"
	"lazypacket/internal/capture"
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
			frame, err := h.ReadFrame()
			if err != nil {
				fmt.Fprintf(os.Stderr, "read failed: %v\n", err)
				break
			}

			p.Send(ui.FrameMsg{Frame: frame})
		}
	}()

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
