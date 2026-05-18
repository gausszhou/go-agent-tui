package component

import (
	"github.com/gausszhou/text-ui-research/tui/theme"
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

	role := m.Role
	if role == "thought" {
		role = "thinking"
	}

	return theme.ChatBg(w).Foreground(fgColor).Render(m.Content)
}
