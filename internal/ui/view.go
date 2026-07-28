package ui

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

const (
	headerHeight = 4
	FooterHeight = 1
	borderSize   = 2
)

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

	height := m.mainPanelHeight()

	packetsPanel := m.renderMainPanel(height / 2)
	detailsPanel := m.renderDetailsPanel(height / 2)

	layout := lipgloss.JoinVertical(lipgloss.Top, header, packetsPanel, detailsPanel, footer)

	view := tea.NewView(layout)
	view.AltScreen = true
	return view
}
