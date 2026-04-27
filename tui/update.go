package tui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/coder/acp-go-sdk"

	"github.com/gausszhou/go-agent-tui/client"
	"github.com/gausszhou/go-agent-tui/tui/component"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.textarea.SetWidth(m.width * 2 / 3)
		m.chatViewport.Width = m.width * 2 / 3
		m.chatViewport.Height = m.height - 6
		return m, nil

	case spinnerTickMsg:
		m.spinner = m.spinner.Tick()
		if m.loading || m.focus == FocusPermission {
			return m, spinnerTick()
		}
		return m, nil

	case initConnMsg:
		if msg.err != nil {
			m.statusText = "Connection failed: " + msg.err.Error()
			m.loading = false
			if m.debug && m.logger != nil {
				m.logger.Error("connection failed", "error", msg.err)
			}
			return m, nil
		}
		m.loading = true
		m.statusText = "Initializing..."
		return m, tea.Batch(
			m.doInitialize(msg.conn, msg.acpClient),
			m.waitForEvents(),
			spinnerTick(),
		)

	case initDoneMsg:
		if msg.err != nil {
			m.statusText = "Init failed: " + msg.err.Error()
			m.loading = false
			m.errMsg = msg.err.Error()
			if m.debug && m.logger != nil {
				m.logger.Error("init failed", "error", msg.err)
			}
			return m, nil
		}
		m.activeSessionID = msg.sessionID
		m.sessions = append(m.sessions, Session{
			ID:   msg.sessionID,
			Name: "Session 1",
			CWD:  client.MustCwd(),
		})
		m.sessionList.Sessions = append(m.sessionList.Sessions, component.SessionItem{
			ID:     msg.sessionID,
			Name:   "Session 1",
			Active: true,
		})
		m.loading = false
		m.statusText = "Ready"
		if m.debug && m.logger != nil {
			m.logger.Info("ready", "session", msg.sessionID)
		}
		return m, nil

	case acpEventMsg:
		switch msg.event.Kind {
		case client.AcpAgentChunk:
			m.appendAgentText(msg.event.Text)
			m.updateChatViewport()
		case client.AcpToolCall:
			tc := msg.event.ToolCall
			if tc != nil {
				title := tc.Title
				m.addMessage(component.ChatMessage{
					Role:          component.RoleTool,
					ToolCallTitle: title,
					ToolCallID:    string(tc.ToolCallId),
					ToolStatus:    string(tc.Status),
				})
				m.updateChatViewport()
			}
		case client.AcpToolUpdate:
			tu := msg.event.ToolUpdate
			if tu != nil {
				status := ""
				if tu.Status != nil {
					status = string(*tu.Status)
				}
				m.addMessage(component.ChatMessage{
					Role:          component.RoleSystem,
					Content:       "Tool " + string(tu.ToolCallId) + " status: " + status,
				})
				m.updateChatViewport()
			}
		case client.AcpPlan:
			plan := msg.event.Plan
			if plan != nil {
				for i, entry := range plan.Entries {
					m.todoList.AddItem(component.TodoItem{
						ID:     fmt.Sprintf("task-%d", i),
						Title:  entry.Content,
						Status: component.TodoPending,
					})
				}
			}
		case client.AcpError:
			m.errMsg = msg.event.Error.Error()
			m.statusText = "Error: " + msg.event.Error.Error()
			m.loading = false
		}
		return m, m.waitForEvents()

	case permEventMsg:
		m.focus = FocusPermission
		m.pendingPerm = &msg.event
		m.questionBox = makePermissionQuestionBox(msg.event.Request, min(m.width-10, 60))
		m.statusText = "Permission requested"
		return m, nil

	case promptDoneMsg:
		m.promptRunning = false
		m.promptCancel = nil
		m.loading = false
		if msg.err != nil {
			m.statusText = "Prompt error: " + msg.err.Error()
			m.errMsg = msg.err.Error()
			if m.debug && m.logger != nil {
				m.logger.Error("prompt done error", "error", msg.err)
			}
		} else {
			m.statusText = "Ready"
		}
		return m, nil

	case sessionCreatedMsg:
		m.loading = false
		if msg.err != nil {
			m.statusText = "Session error: " + msg.err.Error()
			m.errMsg = msg.err.Error()
			return m, nil
		}
		count := len(m.sessions) + 1
		s := Session{
			ID:   msg.sessionID,
			Name: "Session " + formatIntStr2(count),
			CWD:  client.MustCwd(),
		}
		m.sessions = append(m.sessions, s)
		m.sessionList.Sessions = append(m.sessionList.Sessions, component.SessionItem{
			ID:     msg.sessionID,
			Name:   s.Name,
			Active: true,
		})
		for i := range m.sessionList.Sessions {
			m.sessionList.Sessions[i].Active = (m.sessionList.Sessions[i].ID == msg.sessionID)
		}
		m.activeSessionID = msg.sessionID
		m.messages = nil
		m.todoList.Items = nil
		m.statusText = "New session created"
		return m, nil

	case tea.KeyMsg:
		return m.handleKeyMsg(msg)

	case tea.MouseMsg:
		return m.handleMouseMsg(msg)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.focus {
	case FocusPermission:
		return m.handlePermissionKey(msg)

	case FocusSessionList:
		return m.handleSessionListKey(msg)

	case FocusInput:
		return m.handleInputKey(msg)
	}
	return m, nil
}

func (m Model) handlePermissionKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		m.questionBox.Up()
		return m, nil

	case "down", "j":
		m.questionBox.Down()
		return m, nil

	case "enter", " ":
		if m.pendingPerm != nil {
			selectedIdx := m.questionBox.Selected()
			options := m.pendingPerm.Request.Options
			var resp acp.RequestPermissionResponse
			if selectedIdx >= 0 && selectedIdx < len(options) {
				resp = acp.RequestPermissionResponse{
					Outcome: acp.RequestPermissionOutcome{
						Selected: &acp.RequestPermissionOutcomeSelected{
							OptionId: options[selectedIdx].OptionId,
						},
					},
				}
			} else {
				resp = acp.RequestPermissionResponse{
					Outcome: acp.RequestPermissionOutcome{},
				}
			}
			m.pendingPerm.Response <- resp
			m.pendingPerm = nil
		}
		m.focus = FocusInput
		m.statusText = "Ready"
		return m, m.waitForEvents()

	case "esc":
		if m.pendingPerm != nil {
			m.pendingPerm.Response <- acp.RequestPermissionResponse{
				Outcome: acp.RequestPermissionOutcome{},
			}
			m.pendingPerm = nil
		}
		m.focus = FocusInput
		m.statusText = "Ready"
		return m, m.waitForEvents()

	case "ctrl+c":
		if m.pendingPerm != nil {
			m.pendingPerm.Response <- acp.RequestPermissionResponse{
				Outcome: acp.RequestPermissionOutcome{},
			}
			m.pendingPerm = nil
		}
		m.cleanup()
		return m, tea.Quit
	}
	return m, nil
}

func (m Model) handleSessionListKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		m.sessionList.Up()
		return m, nil

	case "down", "j":
		m.sessionList.Down()
		return m, nil

	case "enter":
		idx := m.sessionList.SelectedIdx
		if idx >= 0 && idx < len(m.sessionList.Sessions) {
			sess := m.sessionList.Sessions[idx]
			m.activeSessionID = sess.ID
			for i := range m.sessionList.Sessions {
				m.sessionList.Sessions[i].Active = (m.sessionList.Sessions[i].ID == sess.ID)
			}
			for _, s := range m.sessions {
				if s.ID == sess.ID {
					m.activeSessionID = s.ID
					break
				}
			}
			m.messages = nil
			m.todoList.Items = nil
			m.statusText = "Switched to " + sess.Name
		}
		m.focus = FocusInput
		m.showSessionList = false
		return m, nil

	case "esc", "ctrl+s":
		m.showSessionList = false
		m.focus = FocusInput
		return m, nil

	case "ctrl+c":
		m.cleanup()
		return m, tea.Quit
	}
	return m, nil
}

func (m Model) handleInputKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch key := msg.String(); key {
	case "enter", "ctrl+e":
		now := time.Now()
		pasting := !m.lastKeyTime.IsZero() && now.Sub(m.lastKeyTime) < 20*time.Millisecond
		m.lastKeyTime = now

		if pasting {
			var cmd tea.Cmd
			m.textarea, cmd = m.textarea.Update(msg)
			return m, cmd
		}

		text := strings.TrimSpace(m.textarea.Value())
		if text == "" {
			return m, nil
		}
		m.textarea.Reset()
		return m, tea.Batch(m.submitPrompt(text), spinnerTick())

	case "ctrl+n":
		m.focus = FocusInput
		return m, m.createSession()

	case "ctrl+s":
		if len(m.sessions) > 1 {
			m.showSessionList = true
			m.focus = FocusSessionList
		}
		return m, nil

	case "ctrl+i":
		m.interruptPrompt()
		return m, nil

	case "ctrl+h":
		m.showHelp = !m.showHelp
		return m, nil

	case "ctrl+c":
		m.cleanup()
		return m, tea.Quit

	case "ctrl+up", "ctrl+k":
		m.chatViewport.LineUp(1)
		return m, nil

	case "ctrl+down", "ctrl+j":
		m.chatViewport.LineDown(1)
		return m, nil

	case "pgup":
		m.chatViewport.PageUp()
		return m, nil

	case "pgdown":
		m.chatViewport.PageDown()
		return m, nil

	default:
		m.lastKeyTime = time.Now()
		var cmd tea.Cmd
		m.textarea, cmd = m.textarea.Update(msg)
		return m, cmd
	}
}

func (m Model) handleMouseMsg(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg.Button {
	case tea.MouseButtonWheelUp, tea.MouseButtonWheelDown:
		var cmd tea.Cmd
		m.chatViewport, cmd = m.chatViewport.Update(msg)
		cmds = append(cmds, cmd)

	default:
		if m.focus == FocusInput {
			var cmd tea.Cmd
			m.textarea, cmd = m.textarea.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) updateChatViewport() {
	m.chatViewport.SetContent(m.renderMessages())
	m.chatViewport.GotoBottom()
}

func formatIntStr2(n int) string {
	return fmt.Sprintf("%d", n)
}
