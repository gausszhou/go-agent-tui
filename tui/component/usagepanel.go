package component

import "github.com/charmbracelet/lipgloss"

type UsageInfo struct {
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
	title := ui.TitleStyle.Render("Model")
	model := ui.LabelStyle.Render("Model: ") + ui.ValueStyle.Render(ui.ModelName)
	if ui.ModelName == "" {
		model = ui.LabelStyle.Render("—")
	}

	return lipgloss.JoinVertical(lipgloss.Left,
		title,
		model,
	)
}
