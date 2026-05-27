package tui

import (
	"fmt"
	"strings"

	flex "github.com/gausszhou/bubbleflex"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/gausszhou/bubblecode/tui/component"
	"github.com/gausszhou/bubblecode/tui/layout"
	"github.com/gausszhou/bubblecode/tui/overlay"
	"github.com/gausszhou/bubblecode/tui/theme"
)

func (m *Model) View() tea.View {
	chat := m.chatViewport.View()
	input := m.textarea.View()
	status := m.renderStatus()

	content := lipgloss.JoinVertical(lipgloss.Left, chat, "\n"+input, status)

	if m.showCommands {
		overlayContent := m.renderCommandOverlay()
		content = overlay.CompositeMasked(overlayContent, content, overlay.Center, overlay.Center, 0, 0, true)
	}

	view := tea.NewView(lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Padding(0, layout.PaddingHorizontal).
		Render(content))
	view.AltScreen = true
	view.MouseMode = tea.MouseModeAllMotion

	if c := m.textarea.Cursor(); c != nil && !m.showCommands {
		c.Y += lipgloss.Height(chat) + 1
		view.Cursor = c
	}

	return view
}

func (m *Model) renderCommandOverlay() string {
	panel := component.DefaultCommands()
	bg := theme.ThemeBgOverlay
	var sb strings.Builder
	sb.WriteString(theme.AccentStyle().Background(bg).Render("Commands"))
	sb.WriteString("\n\n")
	for _, cmd := range panel.Commands {
		sb.WriteString("  ")
		sb.WriteString(theme.CommandKeyStyle.Background(bg).Render(cmd.Key))
		sb.WriteString("  ")
		sb.WriteString(theme.CommandDescStyle.Background(bg).Render(cmd.Desc))
		sb.WriteString("\n")
	}
	sb.WriteString("\n")
	sb.WriteString(theme.HelpLabel().Background(bg).Render("Esc to close"))
	return theme.OverlayBox().Render(sb.String())
}

func (m *Model) renderMessages() string {
	w := m.chatViewport.Width()
	var sb strings.Builder
	for i, msg := range m.messages {
		if i > 0 {
			sb.WriteString("\n")
		}
		sb.WriteString(msg.Render(w))
	}
	return sb.String()
}

func comma(n int) string {
	s := fmt.Sprintf("%d", n)
	if len(s) <= 3 {
		return s
	}
	var buf []byte
	for i, c := range s {
		buf = append(buf, byte(c))
		if c == '-' {
			continue
		}
		fromRight := len(s) - i
		if fromRight > 3 && fromRight%3 == 1 {
			buf = append(buf, ',')
		}
	}
	return string(buf)
}

func (m *Model) renderStatus() string {
	left := m.statusText
	if m.loading {
		left = m.spinner.View() + " " + left
	} else {
		left = "✓ " + left
	}
	right := fmt.Sprintf("%s chars", comma(m.chars))
	line := flex.New(flex.Row).
		JustifyContent(flex.SpaceBetween).
		Width(m.width - 2*layout.PaddingHorizontal).
		Join(left, right)
	return theme.StatusBar().Width(m.width - 2*layout.PaddingHorizontal).Render(line)
}
