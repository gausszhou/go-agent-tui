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
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case spinnerTickMsg:
		m.spinner = m.spinner.Tick()
		if m.loading || m.focus == FocusPermission {
			return m, spinnerTick()
		}
		return m, nil

	case renderTickMsg:
		if m.viewportDirty {
			m.viewportDirty = false
			m.updateChatViewport()
		}
		return m, renderTick()

	case outputEventMsg:
		return m.handleOutputEvent(msg.event)

	case tea.KeyMsg:
		return m.handleKeyMsg(msg)

	case tea.MouseMsg:
		return m.handleMouseMsg(msg)
	}

	return m, nil
}

func (m Model) handleOutputEvent(ev client.OutputEvent) (tea.Model, tea.Cmd) {
	if ev.Update != nil {
		return m.handleSessionUpdate(*ev.Update)
	}

	switch ev.Kind {
	case client.EventPromptDone:
		m.promptRunning = false
		m.loading = false
		m.viewportDirty = false
		m.updateChatViewport()
		if m.interrupted {
			m.interrupted = false
			m.statusText = "Ready"
		} else if ev.Error != nil {
			m.statusText = "Error: " + ev.Error.Error()
			m.errMsg = ev.Error.Error()
		} else {
			m.statusText = "Ready"
			m.errMsg = ""
		}
		return m, nil

	case client.EventSessionCreated:
		if ev.Error != nil {
			m.statusText = "Session error: " + ev.Error.Error()
			m.errMsg = ev.Error.Error()
			return m, nil
		}
		count := len(m.sessions) + 1
		s := Session{
			ID:   ev.SessionID,
			Name: "Session " + formatIntStr2(count),
			CWD:  client.MustCwd(),
		}
		m.sessions = append(m.sessions, s)
		m.sessionList.Sessions = append(m.sessionList.Sessions, component.SessionItem{
			ID:     ev.SessionID,
			Name:   s.Name,
			Active: true,
		})
		for i := range m.sessionList.Sessions {
			m.sessionList.Sessions[i].Active = (m.sessionList.Sessions[i].ID == ev.SessionID)
		}
		m.activeSessionID = ev.SessionID
		m.messages = nil
		m.todoList.Items = nil
		m.loading = false
		m.statusText = "New session created"
		return m, nil

	case client.EventSessionLoaded:
		if ev.Error != nil {
			m.statusText = "Load error: " + ev.Error.Error()
			m.errMsg = ev.Error.Error()
		} else {
			m.statusText = "Session loaded"
		}
		m.loading = false
		return m, m.waitForOutput()

	case client.EventPermission:
		if ev.Permission != nil {
			m.focus = FocusPermission
			m.pendingPerm = ev.Permission
			m.questionBox = makePermissionQuestionBox(ev.Permission.Req, min(m.width-10, 60))
			m.statusText = "Permission requested"
		}
		return m, nil

	case client.EventError:
		m.errMsg = ev.Error.Error()
		m.statusText = "Error: " + ev.Error.Error()
		m.loading = false
		return m, nil
	}
	return m, nil
}

func (m Model) handleSessionUpdate(u acp.SessionUpdate) (tea.Model, tea.Cmd) {
	switch {
	case u.UserMessageChunk != nil:
		m.logger.Debug("session update: user message chunk")
		if u.UserMessageChunk.Content.Text != nil && !m.promptRunning {
			m.addMessage(component.ChatMessage{Role: component.RoleUser, Content: u.UserMessageChunk.Content.Text.Text})
			m.viewportDirty = true
		}
	case u.AgentMessageChunk != nil:
		m.logger.Debug("session update: agent message chunk")
		if u.AgentMessageChunk.Content.Text != nil {
			m.appendAgentText(u.AgentMessageChunk.Content.Text.Text)
			m.viewportDirty = true
		}
	case u.ToolCall != nil:
		m.logger.Debug("session update: tool call", "id", string(u.ToolCall.ToolCallId))
		tc := u.ToolCall
		m.addMessage(component.ChatMessage{
			Role:          component.RoleTool,
			ToolCallTitle: tc.Title,
			ToolCallID:    string(tc.ToolCallId),
			ToolStatus:    string(tc.Status),
		})
		m.viewportDirty = true
	case u.ToolCallUpdate != nil:
		tu := u.ToolCallUpdate
		status := ""
		if tu.Status != nil {
			status = string(*tu.Status)
		}
		m.logger.Debug("session update: tool call update", "id", string(tu.ToolCallId), "status", status)
		m.addMessage(component.ChatMessage{
			Role:    component.RoleSystem,
			Content: "Tool " + string(tu.ToolCallId) + " status: " + status,
		})
		m.viewportDirty = true
	case u.Plan != nil:
		m.logger.Debug("session update: plan", "entries", len(u.Plan.Entries))
		plan := u.Plan
		for i, entry := range plan.Entries {
			m.todoList.AddItem(component.TodoItem{
				ID:     fmt.Sprintf("task-%d", i),
				Title:  entry.Content,
				Status: component.TodoPending,
			})
		}
		m.viewportDirty = true
	default:
		m.logger.Debug("session update: unknown type")
	}
	return m, m.waitForOutput()
}

func (m Model) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.focus {
	case FocusPermission:
		return m.handlePermissionKey(msg)
	case FocusSessionList:
		return m.handleSessionListKey(msg)
	case FocusCommandPanel:
		return m.handleCommandPanelKey(msg)
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
			options := m.pendingPerm.Req.Options
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
				resp = acp.RequestPermissionResponse{Outcome: acp.RequestPermissionOutcome{}}
			}
			m.pendingPerm.Response <- resp
			m.pendingPerm = nil
		}
		m.focus = FocusInput
		m.statusText = "Ready"
		return m, m.waitForOutput()
	case "esc":
		if m.pendingPerm != nil {
			m.pendingPerm.Response <- acp.RequestPermissionResponse{
				Outcome: acp.RequestPermissionOutcome{},
			}
			m.pendingPerm = nil
		}
		m.focus = FocusInput
		m.statusText = "Ready"
		return m, m.waitForOutput()
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
			m.statusText = "Loading " + sess.Name + "..."
			m.focus = FocusInput
			m.loading = true
			m.promptRunning = false
			m.messages = nil
			m.todoList.Items = nil
			m.sendInput(client.InputCommand{
				Type:      client.CmdLoadSession,
				SessionID: m.activeSessionID,
			})
			return m, tea.Batch(m.waitForOutput(), spinnerTick())
		}
		m.focus = FocusInput
		return m, nil
	case "esc", "ctrl+s":
		m.focus = FocusInput
		return m, nil
	case "ctrl+c":
		m.cleanup()
		return m, tea.Quit
	}
	return m, nil
}

func (m Model) handleCommandPanelKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.commandPanelIdx > 0 {
			m.commandPanelIdx--
		}
		return m, nil
	case "down", "j":
		if m.commandPanelIdx < 2 {
			m.commandPanelIdx++
		}
		return m, nil
	case "enter":
		switch m.commandPanelIdx {
		case 0:
			m.focus = FocusInput
			m.loading = true
			m.sendInput(client.InputCommand{Type: client.CmdNewSession})
			return m, tea.Batch(m.waitForOutput(), spinnerTick())
		case 1:
			m.focus = FocusSessionList
			m.sessionList.SelectedIdx = 0
			return m, nil
		case 2:
			m.cleanup()
			return m, tea.Quit
		}
		m.focus = FocusInput
		return m, nil
	case "esc", "ctrl+p":
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
		if m.promptRunning {
			return m, nil
		}
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
		m.addMessage(component.ChatMessage{Role: component.RoleUser, Content: text})
		m.updateChatViewport()
		m.promptRunning = true
		m.loading = true
		m.errMsg = ""
		m.statusText = "Processing..."
		m.sendInput(client.InputCommand{
			Type:      client.CmdPrompt,
			SessionID: m.activeSessionID,
			Text:      text,
		})
		return m, tea.Batch(m.waitForOutput(), spinnerTick())

	case "ctrl+n":
		m.focus = FocusInput
		m.loading = true
		m.sendInput(client.InputCommand{Type: client.CmdNewSession})
		return m, tea.Batch(m.waitForOutput(), spinnerTick())

	case "ctrl+p":
		m.focus = FocusCommandPanel
		m.commandPanelIdx = 0
		return m, nil

	case "ctrl+s":
		m.focus = FocusSessionList
		m.sessionList.SelectedIdx = 0
		return m, nil

	case "esc":
		if m.promptRunning {
			now := time.Now()
			if !m.lastEscTime.IsZero() && now.Sub(m.lastEscTime) < 500*time.Millisecond {
				m.interruptPrompt()
				m.lastEscTime = time.Time{}
				return m, nil
			}
			m.lastEscTime = now
			m.statusText = "Press Esc again to interrupt"
			return m, nil
		}
		m.lastEscTime = time.Time{}
		if m.textarea.Value() != "" {
			m.textarea.Reset()
			return m, nil
		}
		return m, nil

	case "ctrl+c":
		m.cleanup()
		return m, tea.Quit

	case "up", "k":
		if m.viewportFocused {
			m.chatViewport.LineUp(1)
			return m, nil
		}
		m.lastKeyTime = time.Now()
		var cmd tea.Cmd
		m.textarea, cmd = m.textarea.Update(msg)
		return m, cmd

	case "down", "j":
		if m.viewportFocused {
			m.chatViewport.LineDown(1)
			return m, nil
		}
		m.lastKeyTime = time.Now()
		var cmd tea.Cmd
		m.textarea, cmd = m.textarea.Update(msg)
		return m, cmd

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

	leftW := m.width * 68 / 100
	barX := leftW - 1

	switch msg.Button {
	case tea.MouseButtonWheelUp, tea.MouseButtonWheelDown:
		var cmd tea.Cmd
		m.chatViewport, cmd = m.chatViewport.Update(msg)
		cmds = append(cmds, cmd)

	case tea.MouseButtonLeft:
		switch msg.Action {
		case tea.MouseActionPress:
			if msg.X == barX && msg.Y >= 0 && msg.Y < m.chatViewport.Height {
				contentLines := visibleLineCount(m.chatViewport.View())
				m.scrollDragging = true
				m.scrollDragStartY = msg.Y
				m.scrollDragStartYOff = m.chatViewport.YOffset
				if contentLines > m.chatViewport.Height && m.chatViewport.Height > 0 {
					thumbH := max(1, m.chatViewport.Height*m.chatViewport.Height/contentLines)
					maxOff := contentLines - m.chatViewport.Height
					thumbY := m.chatViewport.YOffset * (m.chatViewport.Height - thumbH) / maxOff
					if msg.Y < thumbY || msg.Y >= thumbY+thumbH {
						tY := msg.Y * maxOff / (m.chatViewport.Height - thumbH)
						if tY < 0 {
							tY = 0
						}
						if tY > maxOff {
							tY = maxOff
						}
						m.chatViewport.YOffset = tY
					}
				}
			} else {
				m.viewportFocused = msg.X > 0 && msg.X < leftW && msg.Y >= 0 && msg.Y < m.chatViewport.Height
			}
			if m.focus == FocusInput {
				var cmd tea.Cmd
				m.textarea, cmd = m.textarea.Update(msg)
				cmds = append(cmds, cmd)
			}

		case tea.MouseActionMotion:
			if m.scrollDragging && msg.X >= barX-2 && msg.X <= barX+2 {
				contentLines := visibleLineCount(m.chatViewport.View())
				if contentLines > m.chatViewport.Height && m.chatViewport.Height > 0 {
					thumbH := max(1, m.chatViewport.Height*m.chatViewport.Height/contentLines)
					maxOff := contentLines - m.chatViewport.Height
					if m.chatViewport.Height-thumbH > 0 {
						dY := msg.Y - m.scrollDragStartY
						nOff := m.scrollDragStartYOff + dY*maxOff/(m.chatViewport.Height-thumbH)
						if nOff < 0 {
							nOff = 0
						}
						if nOff > maxOff {
							nOff = maxOff
						}
						m.chatViewport.YOffset = nOff
					}
				}
			}

		case tea.MouseActionRelease:
			m.scrollDragging = false
		}

	default:
		if m.focus == FocusInput {
			var cmd tea.Cmd
			m.textarea, cmd = m.textarea.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) waitForOutput() tea.Cmd {
	return func() tea.Msg {
		ev, ok := <-m.outputCh
		if !ok {
			return nil
		}
		return outputEventMsg{event: ev}
	}
}

func formatIntStr2(n int) string {
	return fmt.Sprintf("%d", n)
}
