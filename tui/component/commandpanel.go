package component

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type Command struct {
	Key  string
	Desc string
}

type CommandPanel struct {
	Commands  []Command
	Style     lipgloss.Style
	KeyStyle  lipgloss.Style
	DescStyle lipgloss.Style
}

func NewCommandPanel(commands []Command) CommandPanel {
	return CommandPanel{
		Commands:  commands,
		Style:     lipgloss.NewStyle().Padding(0, 1),
		KeyStyle:  lipgloss.NewStyle().Foreground(lipgloss.Color("#007aff")).Bold(true),
		DescStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("#6e6e73")),
	}
}

func (cp CommandPanel) View() string {
	var sb strings.Builder
	for _, cmd := range cp.Commands {
		sb.WriteString(cp.KeyStyle.Render(cmd.Key))
		sb.WriteString(" ")
		sb.WriteString(cp.DescStyle.Render(cmd.Desc))
		sb.WriteString("  ")
	}
	return cp.Style.Render(strings.TrimRight(sb.String(), " "))
}

func DefaultCommands() CommandPanel {
	return NewCommandPanel([]Command{
		{Key: "Enter", Desc: "Send"},
		{Key: "Shift+Enter", Desc: "Newline"},
		{Key: "Ctrl+N", Desc: "New Session"},
		{Key: "Ctrl+S", Desc: "Switch Session"},
		{Key: "Ctrl+I", Desc: "Interrupt"},
		{Key: "Ctrl+H", Desc: "Help"},
		{Key: "Ctrl+C", Desc: "Quit"},
	})
}
