package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	if m.width == 0 || m.height == 0 {
		return "Initializing..."
	}

	leftWidth := m.width * 70 / 100
	rightWidth := m.width - leftWidth - 1
	if rightWidth < 24 {
		rightWidth = 24
		leftWidth = m.width - rightWidth - 1
	}

	chatHeight := m.height - 5
	if m.showHelp {
		chatHeight -= 1
	}

	leftPanel := m.renderLeftPanel(leftWidth, chatHeight)
	rightPanel := m.renderRightPanel(rightWidth)

	sep := divider().Render("│")

	fullView := lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, sep, rightPanel)
	bgView := base().Width(m.width).Height(m.height).Render(fullView)

	if m.focus == FocusPermission {
		overlay := m.renderPermissionOverlay()
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center,
			overlay,
			lipgloss.WithWhitespaceChars(" "),
			lipgloss.WithWhitespaceForeground(ocBg),
		)
	}

	return bgView
}

func (m Model) renderLeftPanel(width, chatHeight int) string {
	m.chatViewport.Width = width
	m.chatViewport.Height = chatHeight

	chatBox := panelTopBorder("Chat").
		Width(width).
		Height(chatHeight + 2).
		Render(m.chatViewport.View())

	inputArea := m.renderInputArea(width)

	helpLine := ""
	if m.showHelp {
		helpLine = m.renderHelpLine(width)
	}

	m.statusBar.Width = width
	m.statusBar.Loading = m.loading
	m.statusBar.Status = m.statusText
	m.statusBar.Help = "Ctrl+E Send  Ctrl+N Session  Ctrl+S Switch  Ctrl+I Stop  Ctrl+H Help  Ctrl+C Quit"
	m.statusBar.Style = statusBarBg()
	status := m.statusBar.View()

	return lipgloss.JoinVertical(lipgloss.Left, chatBox, inputArea, helpLine, status)
}

func (m Model) renderRightPanel(width int) string {
	margin := lipgloss.NewStyle().Width(width)

	usageBox := panelBorder().Width(width).Render(m.usageInfo.View())
	usageBox = margin.Render(usageBox)

	taskBox := panelBorder().Width(width).Render(m.todoList.View())
	taskBox = margin.Render(taskBox)

	sessionBox := panelBorder().Width(width).Render(m.sessionList.View())
	sessionBox = margin.Render(sessionBox)

	return lipgloss.JoinVertical(lipgloss.Top, usageBox, taskBox, sessionBox)
}

func (m Model) renderInputArea(width int) string {
	m.textarea.SetWidth(width - 2)

	var sb strings.Builder
	if m.errMsg != "" {
		sb.WriteString(errorText().Width(width).Render("! " + m.errMsg))
		sb.WriteString("\n")
	}
	isFocused := m.focus == FocusInput
	if isFocused {
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

	overlayContent := lipgloss.JoinVertical(lipgloss.Left,
		spinner+content,
		"",
		help,
	)

	return overlayBox().MaxWidth(m.width - 10).Render(overlayContent)
}
