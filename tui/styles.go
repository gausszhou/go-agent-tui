package tui

import "github.com/charmbracelet/lipgloss"

var (
	ocBg      = lipgloss.Color("#201d1d")
	ocSurface = lipgloss.Color("#302c2c")
	ocBorder  = lipgloss.Color("#646262")
	ocText    = lipgloss.Color("#fdfcfc")
	ocMuted   = lipgloss.Color("#9a9898")
	ocDim     = lipgloss.Color("#6e6e73")
	ocAccent  = lipgloss.Color("#007aff")
	ocSuccess = lipgloss.Color("#30d158")
	ocWarning = lipgloss.Color("#ff9f0a")
	ocDanger  = lipgloss.Color("#ff3b30")
)

func bg() lipgloss.Color      { return ocBg }
func surface() lipgloss.Color { return ocSurface }
func border() lipgloss.Color  { return ocBorder }
func text() lipgloss.Color    { return ocText }
func muted() lipgloss.Color   { return ocMuted }
func dim() lipgloss.Color     { return ocDim }
func accent() lipgloss.Color  { return ocAccent }
func success() lipgloss.Color { return ocSuccess }
func warning() lipgloss.Color { return ocWarning }
func danger() lipgloss.Color  { return ocDanger }

func agentLabel() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(accent()).Bold(true)
}

func userLabel() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(success()).Bold(true)
}

func thoughtLabel() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(dim()).Italic(true)
}

func toolLabel() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(warning()).Bold(true)
}

func systemLabel() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(muted())
}

func overlayBox() lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(warning()).
		Background(lipgloss.Color("#1e1e1e")).
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
		Background(lipgloss.Color("#0c0c0c")).
		Foreground(muted())
}
