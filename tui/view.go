package tui

import (
	"strings"

	overlay "github.com/rmhubbert/bubbletea-overlay"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	if m.width == 0 || m.height == 0 {
		return "Initializing..."
	}

	bg := m.renderMainView()

	if m.focus == FocusPermission {
		fg := m.renderPermissionOverlay()
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
	m.statusBar.Help = "Enter Send  Ctrl+S Switch  Ctrl+N New  Ctrl+I Stop  Ctrl+C Quit"
	m.statusBar.Style = statusBarBg()
	status := m.statusBar.View()

	return lipgloss.JoinVertical(lipgloss.Left, chat, input, helpLine, status)
}

func (m Model) renderRight(width int) string {
	var parts []string

	usage := base().
		BorderTop(true).BorderStyle(lipgloss.NormalBorder()).
		BorderTopForeground(border()).
		Width(width).Padding(0, 1).
		Render(m.usageInfo.View())
	parts = append(parts, usage)

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
	m.textarea.SetWidth(width - 2)

	var sb strings.Builder
	if m.errMsg != "" {
		sb.WriteString(errorText().Width(width).Render("! " + m.errMsg))
		sb.WriteString("\n")
	}

	if m.focus == FocusInput {
		sb.WriteString(inputBoxFocused().Width(width).Render(m.textarea.View()))
	} else {
		sb.WriteString(inputBox().Width(width).Render(m.textarea.View()))
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
