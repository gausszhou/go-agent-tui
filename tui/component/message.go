package component

import (
	"charm.land/lipgloss/v2"

	"github.com/gausszhou/bubblecode/tui/theme"
)

type Message struct {
	Role    string
	Content string
}

func (m Message) Render(w int) string {
	if w <= 0 {
		w = 80
	}

	fgColor := theme.ThemeText
	switch m.Role {
	case "user":
		fgColor = theme.ThemeUser
	case "agent":
		fgColor = theme.ThemeAgent
	case "thought":
		fgColor = theme.ThemeThought
	case "tool":
		fgColor = theme.ThemeTool
	case "result":
		fgColor = theme.ThemeSuccess
	}

	return lipgloss.NewStyle().Width(w).Foreground(fgColor).Render(m.Content)
}
