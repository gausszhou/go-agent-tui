package component

import (
	"github.com/charmbracelet/lipgloss"
)

type StatusBar struct {
	Loading bool
	Help    string
	Status  string
	Width   int
	Style   lipgloss.Style
}

func NewStatusBar() StatusBar {
	return StatusBar{
		Style: lipgloss.NewStyle().
			Background(lipgloss.Color("#302c2c")).
			Foreground(lipgloss.Color("#9a9898")),
	}
}

func (sb StatusBar) View() string {
	left := sb.Status
	if sb.Loading {
		left = "⏳ " + left
	} else {
		left = "✓ " + left
	}

	right := sb.Help

	leftW := lipgloss.Width(left)
	rightW := lipgloss.Width(right)
	gapW := sb.Width - leftW - rightW
	if gapW < 1 {
		gapW = 1
	}

	return sb.Style.Width(sb.Width).Render(
		left + lipgloss.NewStyle().Width(gapW).Render("") + right,
	)
}
