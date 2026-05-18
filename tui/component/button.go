package component

import "charm.land/lipgloss/v2"

type Button struct {
	Label   string
	Focused bool
	Normal  lipgloss.Style
	Focus   lipgloss.Style
}

func NewButton(label string) Button {
	return Button{
		Label:  label,
		Normal: lipgloss.NewStyle().Foreground(lipgloss.Color("#9a9898")).Padding(0, 1),
		Focus:  lipgloss.NewStyle().Foreground(lipgloss.Color("#fdfcfc")).Background(lipgloss.Color("#007aff")).Padding(0, 1),
	}
}

func (b Button) View() string {
	if b.Focused {
		return b.Focus.Render(b.Label)
	}
	return b.Normal.Render(b.Label)
}
