package tui

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/gausszhou/bubblecode/tui/layout"
	"github.com/gausszhou/bubblecode/tui/theme"
)

func (m *Model) View() tea.View {
	chat := m.chatViewport.View()
	input := m.textarea.View()
	status := m.renderStatus()

	content := lipgloss.JoinVertical(lipgloss.Left, chat, "\n"+input, status)

	view := tea.NewView(lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Padding(0, layout.PaddingHorizontal).
		Render(content))
	view.AltScreen = true
	view.MouseMode = tea.MouseModeAllMotion

	if c := m.textarea.Cursor(); c != nil {
		c.Y += lipgloss.Height(chat) + 1
		view.Cursor = c
	}

	return view
}

func (m *Model) renderStatus() string {
	left := m.statusText
	if m.loading {
		left = m.spinner.View() + " " + left
	} else {
		left = "✓ " + left
	}
	return theme.StatusBar().Width(m.width - 2*layout.PaddingHorizontal).Render(left)
}
