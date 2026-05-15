package component

import (
	"strings"

	"charm.land/lipgloss/v2"

	"github.com/gausszhou/bubblecode/tui/theme"
)

// Role-specific background colors (one shade lighter than border)
var (
	bgAccent  = lipgloss.Color("#1a3050")
	bgThought = lipgloss.Color("#282530")
	bgTool    = lipgloss.Color("#2d2828")
	bgSuccess = lipgloss.Color("#1a3620")
	bgDanger  = lipgloss.Color("#361e1e")
	bgWarning = lipgloss.Color("#362c18")
)

type Message struct {
	Role    string
	Content string
	Status  string
}

func (m Message) Render(w int) string {
	if w <= 0 {
		w = 80
	}

	content := m.Content
	if m.Role == "agent" {
		mdWidth := w - 3
		if mdWidth < 10 {
			mdWidth = 10
		}
		content = RenderMarkdown(content, mdWidth, "")
	}

	var (
		borderColor = theme.ThemeBorder
		bgColor     = theme.ThemeSurface
		contentFg   = theme.ThemeText
	)

	switch m.Role {
	case "user":
		borderColor = theme.ThemeAccent
		bgColor = bgAccent
		contentFg = theme.ThemeText

	case "agent":
		borderColor = theme.ThemeAccent
		bgColor = theme.ThemeBg
		contentFg = theme.ThemeText

	case "thought":
		borderColor = theme.ThemeDim
		bgColor = bgThought
		contentFg = theme.ThemeMuted

	case "tool":
		if m.Status == "error" || m.Status == "failed" {
			borderColor = theme.ThemeDanger
			bgColor = bgDanger
			contentFg = theme.ThemeDanger
		} else if m.Status == "completed" {
			borderColor = theme.ThemeSuccess
			bgColor = bgSuccess
			contentFg = theme.ThemeSuccess
		} else {
			borderColor = theme.ThemeTool
			bgColor = bgTool
			contentFg = theme.ThemeText
		}

	case "result":
		isErr := m.Status == "error" || m.Status == "failed"
		if !isErr && strings.HasPrefix(strings.ToLower(m.Content), "error") {
			isErr = true
		}
		if isErr {
			borderColor = theme.ThemeDanger
			bgColor = bgDanger
			contentFg = theme.ThemeDanger
		} else {
			borderColor = theme.ThemeSuccess
			bgColor = bgSuccess
			contentFg = theme.ThemeSuccess
		}

	case "plan":
		borderColor = theme.ThemeWarning
		bgColor = bgWarning
		contentFg = theme.ThemeText

	default:
		borderColor = theme.ThemeBorder
		bgColor = theme.ThemeSurface
		contentFg = theme.ThemeText
	}

	return lipgloss.NewStyle().
		Border(lipgloss.ThickBorder(), false, false, false, true).
		BorderForeground(borderColor).
		Background(bgColor).
		Width(w).
		Padding(1).
		Foreground(contentFg).
		Render(content) + "\n"
}
