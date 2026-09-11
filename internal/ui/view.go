package ui

import (
	"image/color"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

const (
	headerHeight = 4
	footerHeight = 1
	borderSize   = 2
)

var (
	focusedBorderColor = lipgloss.Color("#AD58B4")
	dimBorderColor     = lipgloss.Color("#3C3C3C")
)

type Layout struct {
	ContentWidth     int
	HeaderHeight     int
	PacketsHeight    int
	DetailsHeight    int
	FooterHeight     int
	InterfacesHeight int
}

func (m Model) layout() Layout {
	contentWidth := max(0, m.width-borderSize)

	used := headerHeight + borderSize + footerHeight + borderSize
	mainHeight := max(0, m.height-used)
	packetsHeight := mainHeight / 2
	detailsHeight := mainHeight - packetsHeight

	return Layout{
		ContentWidth:     contentWidth,
		HeaderHeight:     headerHeight,
		PacketsHeight:    packetsHeight,
		DetailsHeight:    detailsHeight,
		FooterHeight:     footerHeight,
		InterfacesHeight: max(0, m.height-footerHeight-borderSize),
	}
}

func (m Model) borderColor(panel FocusedPanel) color.Color {
	if m.focus == panel {
		return focusedBorderColor
	}
	return dimBorderColor
}

func (m Model) renderInterfaceList() string {
	lines := make([]string, len(m.ifaces))
	for i, iface := range m.ifaces {
		cursor := "  "
		if i == m.ifaceCursor {
			cursor = "> "
		}
		lines[i] = cursor + iface.Name
	}

	var body string
	if len(lines) == 0 {
		body = "no interfaces found"
	}

	list := lipgloss.NewStyle().Align(lipgloss.Left).Render(strings.Join(lines, "\n"))
	body = lipgloss.JoinVertical(lipgloss.Center, "Logo", list)

	l := m.layout()

	return lipgloss.
		NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		Width(l.ContentWidth).
		Height(l.InterfacesHeight).
		Align(lipgloss.Center, lipgloss.Center).
		Render(body)
}

func (m Model) renderHeader() string {
	l := m.layout()

	return lipgloss.
		NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		Width(l.ContentWidth).
		Height(l.HeaderHeight).
		Render("Header Box")
}

func (m Model) packetPanelTop() int {
	return m.layout().HeaderHeight + borderSize
}

func (m Model) renderMainPanel() string {
	l := m.layout()

	return lipgloss.
		NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(m.borderColor(PacketsFocus)).
		Width(l.ContentWidth).
		Height(l.PacketsHeight).
		Render(m.packetList.View())
}

func (m Model) renderDetailsPanel() string {
	l := m.layout()

	return lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(m.borderColor(DetailsFocus)).
		Width(l.ContentWidth).
		Height(l.DetailsHeight).
		Render(m.detailView.View())
}

func (m Model) renderCommandFooter() string {
	l := m.layout()

	return lipgloss.
		NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		Width(l.ContentWidth).
		Height(l.FooterHeight).
		Render("Footer Box")
}

func (m Model) View() tea.View {
	var content string
	if m.screenMode == Interfaces {
		content = lipgloss.JoinVertical(lipgloss.Top, m.renderInterfaceList(), m.renderCommandFooter())
	} else {
		content = lipgloss.JoinVertical(lipgloss.Top, m.renderHeader(), m.renderMainPanel(), m.renderDetailsPanel(), m.renderCommandFooter())
	}

	view := tea.NewView(content)
	view.AltScreen = true
	view.MouseMode = tea.MouseModeCellMotion
	return view
}
