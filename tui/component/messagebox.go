package component

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
)

type MessageRole int

const (
	RoleUser MessageRole = iota
	RoleAgent
	RoleThought
	RoleSystem
	RoleTool
)

type ChatMessage struct {
	Role          MessageRole
	Content       string
	ToolCallTitle string
	ToolCallID    string
	ToolStatus    string
}

func (m ChatMessage) Render(width int, userStyle, agentStyle, thoughtStyle, toolStyle, systemStyle lipgloss.Style) string {
	prefix := ""
	var style lipgloss.Style

	switch m.Role {
	case RoleUser:
		prefix = "You"
		style = userStyle
	case RoleAgent:
		prefix = "Agent"
		style = agentStyle
	case RoleThought:
		prefix = "Thought"
		style = thoughtStyle
	case RoleTool:
		prefix = fmt.Sprintf("Tool: %s (%s)", m.ToolCallTitle, m.ToolStatus)
		style = toolStyle
	case RoleSystem:
		prefix = "System"
		style = systemStyle
	}

	contentWidth := width - 4
	if contentWidth < 20 {
		contentWidth = 20
	}
	wrapped := wordWrap(m.Content, contentWidth)

	var sb strings.Builder
	sb.WriteString(style.Render(prefix))
	sb.WriteString("\n")
	sb.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#9a9898")).PaddingLeft(2).Width(contentWidth).Render(wrapped))

	return sb.String()
}

func wordWrap(text string, width int) string {
	if width <= 0 {
		return text
	}
	var result strings.Builder
	for _, line := range strings.Split(text, "\n") {
		if result.Len() > 0 {
			result.WriteByte('\n')
		}
		if len(line) == 0 {
			continue
		}
		remaining := line
		for len(remaining) > width {
			idx := strings.LastIndexFunc(remaining[:width+1], func(r rune) bool { return r == ' ' })
			if idx <= 0 {
				idx = width
			}
			result.WriteString(strings.TrimRight(remaining[:idx], " "))
			result.WriteByte('\n')
			remaining = remaining[idx:]
			if len(remaining) > 0 && remaining[0] == ' ' {
				remaining = remaining[1:]
			}
		}
		if len(remaining) > 0 {
			result.WriteString(strings.TrimRight(remaining, " "))
		}
	}
	return result.String()
}
