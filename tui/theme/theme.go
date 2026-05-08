package theme

import (
	"image/color"

	"charm.land/lipgloss/v2"
)

var (
	ThemeBg      = lipgloss.Color("#201d1d")
	ThemeSurface = lipgloss.Color("#302c2c")
	ThemeBorder  = lipgloss.Color("#646262")
	ThemeText    = lipgloss.Color("#fdfcfc")
	ThemeMuted   = lipgloss.Color("#9a9898")
	ThemeDim     = lipgloss.Color("#6e6e73")
	ThemeAccent  = lipgloss.Color("#007aff")
	ThemeSuccess = lipgloss.Color("#30d158")
	ThemeWarning = lipgloss.Color("#ff9f0a")
	ThemeDanger  = lipgloss.Color("#ff3b30")

	ThemeUser    = lipgloss.Color("#4cd964")
	ThemeAgent   = lipgloss.Color("#ff9533")
	ThemeThought = lipgloss.Color("#6e6e73")
	ThemeTool    = lipgloss.Color("#ffb300")
	ThemeSystem  = lipgloss.Color("#9a9898")

	ThemeBgDark    = lipgloss.Color("#0c0c0c")
	ThemeBgOverlay = lipgloss.Color("#1e1e1e")
)

func themeBg() color.Color      { return ThemeBg }
func themeSurface() color.Color { return ThemeSurface }
func themeBorder() color.Color  { return ThemeBorder }
func themeText() color.Color    { return ThemeText }
func themeMuted() color.Color   { return ThemeMuted }
func themeDim() color.Color     { return ThemeDim }
func themeAccent() color.Color  { return ThemeAccent }
func themeSuccess() color.Color { return ThemeSuccess }
func themeWarning() color.Color { return ThemeWarning }
func themeDanger() color.Color  { return ThemeDanger }

var (
	StyleUser    = lipgloss.NewStyle().Foreground(ThemeUser).Bold(true)
	StyleAgent   = lipgloss.NewStyle().Foreground(ThemeAgent).Bold(true)
	StyleThought = lipgloss.NewStyle().Foreground(ThemeThought).Italic(true)
	StyleTool    = lipgloss.NewStyle().Foreground(ThemeTool).Bold(true)
	StyleSystem  = lipgloss.NewStyle().Foreground(ThemeSystem)
	StyleContent = lipgloss.NewStyle().Foreground(ThemeMuted)
)
