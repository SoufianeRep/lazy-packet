package ui

import (
	"fmt"
	lp "lazypacket"
	"lazypacket/internal/capture"
	"strings"

	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
)

type Entry struct {
	Frame  capture.Frame
	Packet lp.Packet
}

type Model struct {
	width  int
	height int

	packetList viewport.Model
	detailView viewport.Model
	packets    []Entry
	selected   int
}

type FrameMsg struct {
	Frame capture.Frame
}

func InitModel() Model {
	return Model{
		width:      0,
		height:     0,
		packetList: newViewport(),
		detailView: newViewport(),
	}
}

func packetLines(packets []Entry) string {
	lines := make([]string, len(packets))
	for i, e := range packets {
		lines[i] = fmt.Sprintf("#%d - %s - %d bytes", i, e.Frame.Timestamp.Format("15:04:05"), len(e.Frame.Data))
	}
	return strings.Join(lines, "\n")
}

func newViewport() viewport.Model {
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

		m.packetList.SetWidth(m.width)
		m.packetList.SetHeight(m.mainPanelHeight() / 2)

		m.detailView.SetWidth(m.width)
		m.detailView.SetHeight(m.mainPanelHeight() / 2)

		return m, nil

	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}

	case FrameMsg:
		m.packets = append(m.packets, Entry{Frame: msg.Frame})

		wasAtBottom := m.packetList.AtBottom()
		m.packetList.SetContent(packetLines(m.packets))
		if wasAtBottom {
			m.packetList.GotoBottom()
		}

		return m, nil
	}

	return m, nil
}
