package component

import (
	"charm.land/lipgloss/v2"
	"github.com/gausszhou/bubblecode/tui/theme"
)

type Button struct {
	Label   string
	Focused bool
	Normal  lipgloss.Style
	Focus   lipgloss.Style
}

func NewButton(label string) Button {
	return Button{
		Label:  label,
		Normal: theme.ButtonNormalStyle,
		Focus:  theme.ButtonFocusStyle,
	}
}

func (b Button) View() string {
	if b.Focused {
		return b.Focus.Render(b.Label)
	}
	return b.Normal.Render(b.Label)
}
