package component

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
)

type SessionItem struct {
	ID     string
	Name   string
	Active bool
}

type SessionList struct {
	Sessions    []SessionItem
	Title       string
	SelectedIdx int
	TitleStyle  lipgloss.Style
	ActiveStyle lipgloss.Style
	NormalStyle lipgloss.Style
	SelectStyle lipgloss.Style
}

func NewSessionList(title string) SessionList {
	return SessionList{
		Title:       title,
		TitleStyle:  lipgloss.NewStyle().Foreground(lipgloss.Color("#007aff")).Bold(true),
		ActiveStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("#30d158")),
		NormalStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("#9a9898")),
		SelectStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("#fdfcfc")).Background(lipgloss.Color("#007aff")).Padding(0, 1),
	}
}

func (sl SessionList) View() string {
	var sb strings.Builder
	sb.WriteString(sl.TitleStyle.Render(sl.Title))
	sb.WriteString("\n")

	if len(sl.Sessions) == 0 {
		sb.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#6e6e73")).PaddingLeft(2).Render("No sessions"))
		return sb.String()
	}

	for i, sess := range sl.Sessions {
		marker := ""
		if sess.Active {
			marker = "●"
		}

		if i == sl.SelectedIdx {
			label := fmt.Sprintf("▶ %s %s", marker, truncate(sess.Name, 28))
			sb.WriteString(sl.SelectStyle.Render(label))
		} else {
			prefix := "   "
			label := fmt.Sprintf("%s %s", marker, truncate(sess.Name, 28))
			if sess.Active {
				sb.WriteString(sl.ActiveStyle.Render(prefix + label))
			} else {
				sb.WriteString(sl.NormalStyle.Render(prefix + label))
			}
		}
		sb.WriteString("\n")
	}

	return strings.TrimRight(sb.String(), "\n")
}

func (sl *SessionList) Up() {
	if sl.SelectedIdx > 0 {
		sl.SelectedIdx--
	}
}

func (sl *SessionList) Down() {
	if sl.SelectedIdx < len(sl.Sessions)-1 {
		sl.SelectedIdx++
	}
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-1] + "…"
}
