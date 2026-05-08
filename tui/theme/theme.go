package theme

import (
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

	ThemeUser    = lipgloss.Color("#fdfcfc")
	ThemeAgent   = lipgloss.Color("#9a9898")
	ThemeThought = lipgloss.Color("#6e6e73")
	ThemeTool    = lipgloss.Color("#646262")
	ThemeSystem  = lipgloss.Color("#9a9898")

	ThemeBgDark    = lipgloss.Color("#0c0c0c")
	ThemeBgOverlay = lipgloss.Color("#1e1e1e")
	ThemeInputBg   = lipgloss.Color("#1a1a1a")
	ThemeChatBg    = lipgloss.Color("#000000")
)

func BaseStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Background(ThemeBg).
		Foreground(ThemeText)
}

func TextStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(ThemeText)
}

func MutedStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(ThemeMuted)
}

func DimStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(ThemeDim)
}

func AccentStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(ThemeAccent)
}

func SurfaceStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(ThemeSurface)
}

func BorderStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(ThemeBorder)
}

func PureBlack() lipgloss.Style {
	return lipgloss.NewStyle().Background(lipgloss.Color("#000000")).Foreground(ThemeText)
}

func StatusBar() lipgloss.Style {
	return lipgloss.NewStyle().Background(ThemeBgDark).Foreground(ThemeMuted)
}

func NoBg() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(ThemeText)
}

func ChatBg(w int) lipgloss.Style {
	return lipgloss.NewStyle().Background(ThemeChatBg).Width(w)
}

var (
	StyleUser    = lipgloss.NewStyle().Foreground(ThemeUser).Bold(true)
	StyleAgent   = lipgloss.NewStyle().Foreground(ThemeAgent).Bold(true)
	StyleThought = lipgloss.NewStyle().Foreground(ThemeThought).Italic(true)
	StyleTool    = lipgloss.NewStyle().Foreground(ThemeTool).Bold(true)
	StyleSystem  = lipgloss.NewStyle().Foreground(ThemeSystem)
	StyleContent = lipgloss.NewStyle().Foreground(ThemeMuted)
)
