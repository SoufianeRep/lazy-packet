package ui

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

const (
	headerHeight = 4
	FooterHeight = 1
	borderSize   = 2
)

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

	return lipgloss.
		NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		Width(m.width).
		Height(m.height-FooterHeight-borderSize).
		Align(lipgloss.Center, lipgloss.Center).
		Render(body)
}

func (m Model) renderHeader() string {
	return lipgloss.
		NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		Width(m.width).
		Height(headerHeight).
		Render("Header Box")
}

func (m Model) mainPanelHeight() int {
	header := m.renderHeader()
	footer := m.renderCommandFooter()
	return m.height - lipgloss.Height(header) - lipgloss.Height(footer) - borderSize
}

func (m Model) renderMainPanel(height int) string {
	return lipgloss.
		NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		Width(m.width).
		Height(height).
		Render(m.packetList.View())
}

func (m Model) renderDetailsPanel(height int) string {
	return lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).Height(height).Width(m.width).Render("Details Box")
}

func (m Model) renderCommandFooter() string {
	return lipgloss.
		NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		Width(m.width).
		Height(FooterHeight).
		Render("Footer Box")
}

func (m Model) View() tea.View {
	header := m.renderHeader()
	footer := m.renderCommandFooter()

	var layout string
	if m.screenMode == Interfaces {
		layout = lipgloss.JoinVertical(lipgloss.Top, m.renderInterfaceList(), footer)
	} else {
		height := m.mainPanelHeight()

		packetsPanel := m.renderMainPanel(height / 2)
		detailsPanel := m.renderDetailsPanel(height / 2)

		layout = lipgloss.JoinVertical(lipgloss.Top, header, packetsPanel, detailsPanel, footer)
	}

	view := tea.NewView(layout)
	view.AltScreen = true
	return view
}
