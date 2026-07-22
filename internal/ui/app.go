package ui

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
)

type Model struct {
	width  int
	height int
}

func InitModel() Model {
	return Model{
		width:  0,
		height: 0,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m Model) View() tea.View {
	placeholder := fmt.Sprintf("%d x %d", m.width, m.height)
	return tea.NewView(placeholder)
}
