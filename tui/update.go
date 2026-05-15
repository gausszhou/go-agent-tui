package tui

import (
	"charm.land/bubbles/v2/cursor"
	tea "charm.land/bubbletea/v2"

	"github.com/gausszhou/bubblecode/client"
	"github.com/gausszhou/bubblecode/tui/component"
	"github.com/gausszhou/bubblecode/tui/layout"
)

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m.handleResize(msg)

	case tea.KeyPressMsg:
		if msg.Key().Mod == tea.ModCtrl && msg.Key().Code == 'c' {
			m.cleanup()
			return m, tea.Quit
		}
		return m.handleKey(msg)

	case tea.KeyReleaseMsg:
		return m, nil

	case tea.KeyMsg:
		return m.handleKey(msg)

	case tea.MouseMsg:
		return m.handleMouse(msg)

	case outputEventMsg:
		return m.handleOutputEvent(msg.event)

	case loadingTickMsg:
		m.spinner = m.spinner.Tick()
		return m, spinnerTick()

	case cursor.BlinkMsg:
		return m.handleBlink(msg)
	}
	return m, nil
}

func (m *Model) handleResize(msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	if msg.Width < layout.MinWidth || msg.Height < layout.MinHeight {
		return m, nil
	}
	m.width = msg.Width
	m.height = msg.Height
	m.updateSizes()
	m.chatViewport.SetContent(m.renderMessages())
	return m, nil
}

func (m *Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		return m.sendPrompt()

	case "ctrl+c":
		m.cleanup()
		return m, tea.Quit

	case "up", "k":
		m.chatViewport.ScrollUp(1)
		return m, nil

	case "down", "j":
		m.chatViewport.ScrollDown(1)
		return m, nil
	}

	var cmd tea.Cmd
	m.textarea, cmd = m.textarea.Update(msg)
	return m, cmd
}

func (m *Model) handleMouse(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	if wheel, ok := msg.(tea.MouseWheelMsg); ok {
		switch wheel.Button {
		case tea.MouseWheelUp:
			m.chatViewport.ScrollUp(3)
		case tea.MouseWheelDown:
			m.chatViewport.ScrollDown(3)
		}
	}
	return m, nil
}

func (m *Model) handleOutputEvent(ev client.OutputEvent) (tea.Model, tea.Cmd) {
	if ev.Update != nil {
		m.processUpdate(ev.Update.Update)
		m.refreshChat()
		return m, waitForOutput(m.outputCh)
	}

	switch ev.Kind {
	case "done":
		m.promptRunning = false
		m.loading = false
		m.statusText = "Ready"
		m.refreshChat()
	case "error":
		m.promptRunning = false
		m.loading = false
		m.statusText = "Error: " + ev.Error.Error()
	}
	return m, nil
}

func (m *Model) handleBlink(msg cursor.BlinkMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.textarea, cmd = m.textarea.Update(msg)
	return m, cmd
}

func (m *Model) sendPrompt() (tea.Model, tea.Cmd) {
	if m.promptRunning {
		return m, nil
	}
	text := m.textarea.Value()
	if text == "" {
		return m, nil
	}

	m.textarea.Reset()
	m.messages = append(m.messages, component.Message{Role: roleUser, Content: text})
	m.promptRunning = true
	m.loading = true
	m.statusText = "Processing..."
	m.chatViewport.SetContent(m.renderMessages())

	return m, tea.Batch(
		sendInput(m.inputCh, client.InputCommand{Type: client.CmdPrompt, Text: text}),
		waitForOutput(m.outputCh),
		spinnerTick(),
	)
}
