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

	left := m.renderLeft(leftW, chatH)
	right := m.renderRight(rightW)

	return lipgloss.JoinHorizontal(lipgloss.Top, left, divider().Render("│"), right)
}

func (m Model) renderLeft(width, chatH int) string {
	barW := 1
	vpW := width - barW

	m.chatViewport.Width = vpW
	m.chatViewport.Height = chatH
	content := m.renderMessages()
	m.chatViewport.SetContent(content)
	vp := m.chatViewport.View()

	sb := renderScrollbar(chatH, strings.Count(content, "\n")+1, m.chatViewport.YOffset)
	chat := lipgloss.JoinHorizontal(lipgloss.Top, vp, sb)

	input := m.renderInput(width)

	help := "Enter Send  Ctrl+P Commands  Ctrl+C Quit"
	if m.promptRunning {
		help = "Enter Send  Esc Esc Interrupt  Ctrl+P Commands  Ctrl+C Quit"
	}

	m.statusBar.Width = width
	m.statusBar.Loading = m.loading
	m.statusBar.Status = m.statusText
	m.statusBar.Help = help
	m.statusBar.Style = statusBarBg()
	status := m.statusBar.View()

	return lipgloss.JoinVertical(lipgloss.Left, chat, input, status)
}

func (m Model) renderRight(width int) string {
	var parts []string

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
	promptWidth := 2
	m.textarea.SetWidth(width - promptWidth)

	promptColor := muted()
	if m.focus == FocusInput {
		promptColor = accent()
	}
	prompt := lipgloss.NewStyle().Foreground(promptColor).Bold(true).Render("❯ ")

	var sb strings.Builder
	if m.errMsg != "" {
		sb.WriteString(errorText().Width(width).Render("! " + m.errMsg))
		sb.WriteString("\n")
	}

	lines := strings.Split(m.textarea.View(), "\n")
	for i, line := range lines {
		if i > 0 {
			sb.WriteString("\n")
		}
		if i == 0 {
			sb.WriteString(prompt + line)
		} else {
			sb.WriteString("   " + line)
		}
	}

	return sb.String()
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
	}{
		{"New Session"},
		{"Switch Session"},
		{"Quit"},
	}

	for i, opt := range options {
		if i == m.commandPanelIdx {
			sb.WriteString(lipgloss.NewStyle().Foreground(text()).Background(accent()).Padding(0, 1).Render("▶ " + opt.label))
		} else {
			sb.WriteString(lipgloss.NewStyle().Foreground(muted()).Padding(0, 1).Render("  " + opt.label))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("\n")
	sb.WriteString(helpLabel().Render("↑↓ navigate  Enter select  Esc cancel"))

	content := sb.String()
	return overlayBox().Width(30).Render(content)
}

func renderScrollbar(height, totalLines, yOffset int) string {
	if totalLines <= height || height <= 0 {
		return lipgloss.NewStyle().Width(1).Height(height).Render("")
	}

	thumbH := max(1, height*height/totalLines)
	maxOffset := totalLines - height
	if maxOffset <= 0 {
		return lipgloss.NewStyle().Width(1).Height(height).Render("")
	}
	thumbY := yOffset * (height - thumbH) / maxOffset

	trackStyle := lipgloss.NewStyle().Foreground(dim())
	thumbStyle := lipgloss.NewStyle().Foreground(border())

	var sb strings.Builder
	for i := 0; i < height; i++ {
		if i >= thumbY && i < thumbY+thumbH {
			sb.WriteString(thumbStyle.Render("█"))
		} else {
			sb.WriteString(trackStyle.Render("│"))
		}
		sb.WriteString("\n")
	}
	return sb.String()
}
