package tui

import (
	"image/color"

	"charm.land/lipgloss/v2"
	"github.com/gausszhou/text-ui-research/tui/theme"
)

func bg() color.Color      { return theme.ThemeBg }
func surface() color.Color { return theme.ThemeSurface }
func border() color.Color  { return theme.ThemeBorder }
func text() color.Color    { return theme.ThemeText }
func muted() color.Color   { return theme.ThemeMuted }
func dim() color.Color     { return theme.ThemeDim }
func accent() color.Color  { return theme.ThemeAccent }
func success() color.Color { return theme.ThemeSuccess }
func warning() color.Color { return theme.ThemeWarning }
func danger() color.Color  { return theme.ThemeDanger }

func overlayBox() lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(warning()).
		Background(theme.ThemeBgOverlay).
		Padding(1, 2)
}

func loadingSpinner() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(accent())
}

func helpLabel() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(dim()).Padding(0, 1)
}

func statusBarBg() lipgloss.Style {
	return lipgloss.NewStyle().
		Background(theme.ThemeBgDark).
		Foreground(muted())
}
