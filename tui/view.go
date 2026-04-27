package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	overlay "github.com/rmhubbert/bubbletea-overlay"
)

func (m Model) View() string {
	if m.width == 0 || m.height == 0 {
		return "Initializing..."
	}

	bg := m.renderMainView()

	switch m.focus {
	case FocusPermission:
		fg := m.renderPermissionOverlay()
		return overlay.Composite(fg, bg, overlay.Center, overlay.Center, 0, 0)
	case FocusCommandPanel:
		fg := m.renderCommandPanelOverlay()
		return overlay.Composite(fg, bg, overlay.Center, overlay.Center, 0, 0)
	}

	return bg
}

func (m Model) renderMainView() string {
	leftW := m.width * 68 / 100
	rightW := m.width - leftW - 1
	if rightW < 22 {
		rightW = 22
		leftW = m.width - rightW - 1
	}
	chatH := m.height - 4
	if m.showHelp {
		chatH--
	}

	left := m.renderLeft(leftW, chatH)
	right := m.renderRight(rightW)

	return lipgloss.JoinHorizontal(lipgloss.Top, left, divider().Render("│"), right)
}

func (m Model) renderLeft(width, chatH int) string {
	m.chatViewport.Width = width
	m.chatViewport.Height = chatH
	m.chatViewport.SetContent(m.renderMessages())
	m.chatViewport.GotoBottom()

	chat := m.chatViewport.View()

	input := m.renderInput(width)

	helpLine := ""
	if m.showHelp {
		helpLine = m.renderHelpLine(width)
	}

	m.statusBar.Width = width
	m.statusBar.Loading = m.loading
	m.statusBar.Status = m.statusText
	m.statusBar.Help = "Enter Send  Esc Esc Interrupt  Ctrl+P Commands  Ctrl+C Quit"
	m.statusBar.Style = statusBarBg()
	status := m.statusBar.View()

	return lipgloss.JoinVertical(lipgloss.Left, chat, input, helpLine, status)
}

func (m Model) renderRight(width int) string {
	var parts []string

	model := base().
		BorderTop(true).BorderStyle(lipgloss.NormalBorder()).
		BorderTopForeground(border()).
		Width(width).Padding(0, 1).
		Render(m.usageInfo.View())
	parts = append(parts, model)

	tasks := base().
		BorderTop(true).BorderStyle(lipgloss.NormalBorder()).
		BorderTopForeground(border()).
		Width(width).Padding(0, 1).
		Render(m.todoList.View())
	parts = append(parts, tasks)

	sess := base().
		BorderTop(true).BorderStyle(lipgloss.NormalBorder()).
		BorderTopForeground(border()).
		Width(width).Padding(0, 1).
		Render(m.sessionList.View())
	parts = append(parts, sess)

	return lipgloss.JoinVertical(lipgloss.Top, parts...)
}

func (m Model) renderInput(width int) string {
	m.textarea.SetWidth(width - 4)

	focused := m.focus == FocusInput
	promptColor := muted()
	if focused {
		promptColor = accent()
	}

	prompt := lipgloss.NewStyle().Foreground(promptColor).Render("❯ ")
	content := m.textarea.View()

	var sb strings.Builder
	if m.errMsg != "" {
		sb.WriteString(errorText().Width(width).Render("! " + m.errMsg))
		sb.WriteString("\n")
	}

	lines := strings.Split(content, "\n")
	for i, line := range lines {
		if i > 0 {
			sb.WriteString("\n")
		}
		if i == 0 {
			sb.WriteString(prompt + line)
		} else {
			sb.WriteString(lipgloss.NewStyle().PaddingLeft(3).Render(line))
		}
	}

	return sb.String()
}

func (m Model) renderHelpLine(width int) string {
	m.commandPanel.Style = m.commandPanel.Style.Width(width)
	return m.commandPanel.View()
}

func (m Model) renderPermissionOverlay() string {
	w := min(m.width-10, 64)
	m.questionBox.Width = w
	m.questionBox.Style = overlayBox()
	m.questionBox.TitleStyle = lipgloss.NewStyle().Foreground(warning()).Bold(true)

	content := m.questionBox.View()

	spinner := ""
	if m.loading {
		spinner = m.spinner.View() + " "
	}

	help := helpLabel().Render("↑↓ navigate  Enter select  Esc deny  Ctrl+C quit")

	ov := lipgloss.JoinVertical(lipgloss.Left,
		spinner+content,
		"",
		help,
	)

	return overlayBox().MaxWidth(m.width - 10).Render(ov)
}

func (m Model) renderCommandPanelOverlay() string {
	var sb strings.Builder

	title := lipgloss.NewStyle().Foreground(accent()).Bold(true).Render("Commands")
	sb.WriteString(title)
	sb.WriteString("\n\n")

	options := []struct {
		label string
		key   string
	}{
		{"New Session", "N"},
		{"Switch Session", "S"},
	}

	for i, opt := range options {
		prefix := "  "
		if i == m.commandPanelIdx {
			prefix = "▶ "
		}
		if i == m.commandPanelIdx {
			sb.WriteString(lipgloss.NewStyle().Foreground(text()).Background(accent()).Padding(0, 1).Render(prefix + opt.label))
		} else {
			sb.WriteString(lipgloss.NewStyle().Foreground(muted()).Padding(0, 1).Render(prefix + opt.label))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("\n")
	help := helpLabel().Render("↑↓ navigate  Enter select  Esc cancel")
	sb.WriteString(help)

	content := sb.String()
	return overlayBox().Width(30).Render(content)
}
