package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	overlay "github.com/gausszhou/go-agent-tui/tui/overlay"
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
	case FocusSessionList:
		fg := m.renderSessionOverlay()
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
	chatH := m.height - 7

	left := m.renderLeft(leftW, chatH)
	right := m.renderRight(rightW)

	return lipgloss.JoinHorizontal(lipgloss.Top, left, right)
}

func (m Model) renderLeft(width, chatH int) string {
	barW := 1
	vpW := width - barW

	m.chatViewport.Width = vpW
	m.chatViewport.Height = chatH
	vp := m.chatViewport.View()

	contentLines := m.chatViewport.TotalLineCount()
	sb := renderScrollbar(chatH, contentLines, m.chatViewport.YOffset)
	chat := lipgloss.JoinHorizontal(lipgloss.Top, vp, sb)

	input := m.renderInput(width)

	help := "Enter Send  Ctrl+P Commands  Ctrl+C Quit"
	if m.promptRunning {
		help = "Esc Esc Interrupt  Ctrl+P Commands  Ctrl+C Quit"
	}

	m.statusBar.Width = width
	m.statusBar.Loading = m.loading
	m.statusBar.Status = m.statusText
	m.statusBar.Help = help
	m.statusBar.Style = statusBarBg()
	status := m.statusBar.View()

	return lipgloss.JoinVertical(lipgloss.Left, chat, input, lipgloss.NewStyle().Height(1).Render(""), status)
}

func (m Model) renderRight(width int) string {
	return lipgloss.NewStyle().Width(width).Height(m.height).
		Background(lipgloss.Color("#141414")).
		Render(lipgloss.NewStyle().Width(width).Padding(0, 1).
			Background(lipgloss.Color("#141414")).
			Foreground(text()).
			Render(m.todoList.View()))
}

func (m Model) renderInput(width int) string {
	promptWidth := 2
	m.textarea.SetWidth(width - promptWidth)

	promptColor := muted()
	if m.focus == FocusInput {
		promptColor = accent()
	}
	m.textarea.Prompt = lipgloss.NewStyle().Width(promptWidth).Foreground(promptColor).Render("┃ ")

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
			sb.WriteString(line)
		} else {
			sb.WriteString(line)
		}
	}

	return lipgloss.NewStyle().Width(width).Background(bg()).Render(sb.String())
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

	itemW := 45

	for i, opt := range options {
		if i == m.commandPanelIdx {
			sb.WriteString(lipgloss.NewStyle().Width(itemW).Foreground(text()).Background(accent()).Padding(0, 1).Render("▶ " + opt.label))
		} else {
			sb.WriteString(lipgloss.NewStyle().Width(itemW).Foreground(muted()).Background(lipgloss.Color("#1e1e1e")).Padding(0, 1).Render("  " + opt.label))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("\n")
	sb.WriteString(helpLabel().Background(lipgloss.Color("#1e1e1e")).Render("↑↓ navigate  Enter select  Esc cancel"))

	content := sb.String()
	return overlayBox().Width(50).Render(content)
}

func (m Model) renderSessionOverlay() string {
	var sb strings.Builder

	title := lipgloss.NewStyle().Foreground(accent()).Bold(true).Render("Sessions")
	sb.WriteString(title)
	sb.WriteString("\n\n")

	itemW := 45

	for i, sess := range m.sessionList.Sessions {
		marker := " "
		color := muted()
		if sess.Active {
			marker = "●"
			color = success()
		}
		if i == m.sessionList.SelectedIdx {
			sb.WriteString(lipgloss.NewStyle().Width(itemW).Foreground(text()).Background(accent()).Padding(0, 1).Render("▶ " + marker + " " + sess.Name))
		} else {
			sb.WriteString(lipgloss.NewStyle().Width(itemW).Foreground(color).Background(lipgloss.Color("#1e1e1e")).Padding(0, 1).Render("  " + marker + " " + sess.Name))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("\n")
	sb.WriteString(helpLabel().Background(lipgloss.Color("#1e1e1e")).Render("↑↓ navigate  Enter select  Esc cancel"))

	return overlayBox().Width(50).Render(sb.String())
}

func renderScrollbar(height, totalLines, yOffset int) string {
	if height <= 0 {
		return ""
	}
	trackStyle := lipgloss.NewStyle().Foreground(dim())
	thumbStyle := lipgloss.NewStyle().Foreground(border())

	var sb strings.Builder
	if totalLines <= height {
		for i := 0; i < height; i++ {
			sb.WriteString(trackStyle.Render("│"))
			sb.WriteByte('\n')
		}
		return sb.String()
	}

	thumbH := max(1, height*height/totalLines)
	maxOffset := totalLines - height
	if maxOffset <= 0 {
		for i := 0; i < height; i++ {
			sb.WriteString(trackStyle.Render("│"))
			sb.WriteByte('\n')
		}
		return sb.String()
	}
	thumbY := yOffset * (height - thumbH) / maxOffset

	for i := 0; i < height; i++ {
		if i >= thumbY && i < thumbY+thumbH {
			sb.WriteString(thumbStyle.Render("█"))
		} else {
			sb.WriteString(trackStyle.Render("│"))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}
