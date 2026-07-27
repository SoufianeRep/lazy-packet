package ui

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
)

type Model struct {
	width  int
	height int

	frames  [][]byte
	packets viewport.Model
}

type FrameMsg struct {
	Data []byte
}

func InitModel() Model {
	return Model{
		width:   0,
		height:  0,
		packets: InitFramesModel(),
	}
}

func frameLines(frames [][]byte) string {
	lines := make([]string, len(frames))
	for i, f := range frames {
		lines[i] = fmt.Sprintf("#%d - %d bytes", i, len(f))
	}
	return strings.Join(lines, "\n")
}

func InitFramesModel() viewport.Model {
	vp := viewport.New()
	vp.SetWidth(0)
	vp.SetHeight(0)

	return vp
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		m.packets.SetWidth(m.width/2 - borderSize)
		m.packets.SetHeight(m.mainPanelHeight() - borderSize)

		return m, nil

	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}

	case FrameMsg:
		m.frames = append(m.frames, msg.Data)

		wasAtBottom := m.packets.AtBottom()
		m.packets.SetContent(frameLines(m.frames))
		if wasAtBottom {
			m.packets.GotoBottom()
		}

		return m, nil
	}

	return m, nil
}
