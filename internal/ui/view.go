package ui

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

const (
	headerHeight = 4
	FooterHeight = 1
)

func (m Model) renderHeader() string {
	return lipgloss.
		NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		Width(m.width).
		Height(headerHeight).
		Render("Header Box")
}

func (m Model) renderMainPanel(height int) string {
	return lipgloss.
		NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		Width(m.width / 2).
		Height(height).
		Render("Hello this is a box")
}

func (m Model) renderDetailsPanel(height int) string {
	return lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).Height(height).Width(m.width / 2).Render("Details Box")
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

	usedHeight := lipgloss.Height(header) + lipgloss.Height(footer)

	packetsPanel := m.renderMainPanel(m.height - usedHeight)
	detailsPanel := m.renderDetailsPanel(m.height - usedHeight)
	mainPanel := lipgloss.JoinHorizontal(lipgloss.Top, packetsPanel, detailsPanel)

	layout := lipgloss.JoinVertical(lipgloss.Top, header, mainPanel, footer)

	view := tea.NewView(layout)
	view.AltScreen = true
	return view
}
