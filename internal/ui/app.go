package ui

import (
	lp "lazypacket"
	"strings"
	"time"

	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
)

type FrameMsg struct {
	Packet    *lp.Packet
	TimeStamp time.Time
}
type Entry struct {
	Packet    *lp.Packet
	TimeStamp time.Time
}
type Model struct {
	width  int
	height int

	packetList viewport.Model
	detailView viewport.Model
	packets    []Entry
	selected   int
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
		lines[i] = formatPacket(i+1, e, e.TimeStamp.Sub(packets[0].TimeStamp))
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
		m.packets = append(m.packets, Entry{Packet: msg.Packet, TimeStamp: msg.TimeStamp})

		wasAtBottom := m.packetList.AtBottom()
		m.packetList.SetContent(packetLines(m.packets))
		if wasAtBottom {
			m.packetList.GotoBottom()
		}

		return m, nil
	}

	return m, nil
}
