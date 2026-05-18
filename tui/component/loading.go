package component

import "charm.land/lipgloss/v2"

type Loading struct {
	Frame int
	Style lipgloss.Style
}

var spinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

func NewLoading(style lipgloss.Style) Loading {
	return Loading{Style: style}
}

func (l Loading) Tick() Loading {
	l.Frame = (l.Frame + 1) % len(spinnerFrames)
	return l
}

func (l Loading) View() string {
	if l.Frame >= len(spinnerFrames) {
		return ""
	}
	return l.Style.Render(spinnerFrames[l.Frame])
}
