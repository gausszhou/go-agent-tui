package component

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type QuestionBox struct {
	Title         string
	Message       string
	Options       []string
	SelectedIdx   int
	Width         int
	Style         lipgloss.Style
	ActiveStyle   lipgloss.Style
	InactiveStyle lipgloss.Style
	TitleStyle    lipgloss.Style
}

func NewQuestionBox(title, message string, options []string, width int) QuestionBox {
	return QuestionBox{
		Title:         title,
		Message:       message,
		Options:       options,
		SelectedIdx:   0,
		Width:         width,
		Style:         lipgloss.NewStyle().Border(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("#ff9f0a")).Padding(1, 2),
		ActiveStyle:   lipgloss.NewStyle().Foreground(lipgloss.Color("#fdfcfc")).Background(lipgloss.Color("#007aff")).Padding(0, 1),
		InactiveStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("#9a9898")).Padding(0, 1),
		TitleStyle:    lipgloss.NewStyle().Foreground(lipgloss.Color("#ff9f0a")).Bold(true),
	}
}

func (q QuestionBox) View() string {
	var sb strings.Builder

	sb.WriteString(q.TitleStyle.Render(q.Title))
	sb.WriteString("\n\n")
	sb.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#fdfcfc")).Render(q.Message))
	sb.WriteString("\n\n")

	for i, opt := range q.Options {
		if i == q.SelectedIdx {
			prefix := "▶ "
			sb.WriteString(q.ActiveStyle.Render(fmt.Sprintf("%s%d. %s", prefix, i+1, opt)))
		} else {
			prefix := "  "
			sb.WriteString(q.InactiveStyle.Render(fmt.Sprintf("%s%d. %s", prefix, i+1, opt)))
		}
		sb.WriteString("\n")
	}

	content := sb.String()
	return q.Style.Width(q.Width).Render(content)
}

func (q *QuestionBox) Up() {
	if q.SelectedIdx > 0 {
		q.SelectedIdx--
	}
}

func (q *QuestionBox) Down() {
	if q.SelectedIdx < len(q.Options)-1 {
		q.SelectedIdx++
	}
}

func (q QuestionBox) Selected() int {
	return q.SelectedIdx
}
