package main

import (
	"fmt"
	"lazypacket/internal/ui"
	"os"

	tea "charm.land/bubbletea/v2"
)

func main() {
	model := ui.InitModel()
	p := tea.NewProgram(model)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
