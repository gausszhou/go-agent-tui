package component

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

type UsageInfo struct {
	Tokens     int
	Cost       float64
	ModelName  string
	TitleStyle lipgloss.Style
	LabelStyle lipgloss.Style
	ValueStyle lipgloss.Style
}

func NewUsageInfo() UsageInfo {
	return UsageInfo{
		TitleStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("#007aff")).Bold(true),
		LabelStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("#6e6e73")),
		ValueStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("#fdfcfc")),
	}
}

func (ui UsageInfo) View() string {
	title := ui.TitleStyle.Render("Usage")
	model := ui.LabelStyle.Render("Model: ") + ui.ValueStyle.Render(ui.ModelName)
	tokens := ui.LabelStyle.Render("Tokens: ") + ui.ValueStyle.Render(formatTokens(ui.Tokens))
	cost := ui.LabelStyle.Render("Cost: ") + ui.ValueStyle.Render(formatCost(ui.Cost))

	return lipgloss.JoinVertical(lipgloss.Left,
		title,
		model,
		tokens,
		cost,
	)
}

func formatTokens(n int) string {
	if n == 0 {
		return "—"
	}
	if n >= 1000 {
		return fmt.Sprintf("%.1fk", float64(n)/1000.0)
	}
	return fmt.Sprintf("%d", n)
}

func formatCost(c float64) string {
	if c == 0 {
		return "$0.00"
	}
	return fmt.Sprintf("$%.2f", c)
}
