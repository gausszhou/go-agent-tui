package tui

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os/exec"
	"strings"
	"time"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/textarea"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/coder/acp-go-sdk"

	"github.com/gausszhou/text-ui-research/client"
	"github.com/gausszhou/text-ui-research/tui/component"
	"github.com/gausszhou/text-ui-research/tui/theme"
)

type Model struct {
	logger   *slog.Logger
	inputCh  chan client.InputCommand
	outputCh chan client.OutputEvent
	cmd      *exec.Cmd
	ctx      context.Context
	cancel   context.CancelFunc

	width  int
	height int

	textarea     textarea.Model
	chatViewport viewport.Model

	messages []component.Message

	loading       bool
	promptRunning bool
	statusText    string
	lastKeyTime   time.Time
}

func NewModel(logger *slog.Logger, cmd *exec.Cmd, sessionID string, ctx context.Context, cancel context.CancelFunc, inputCh chan client.InputCommand, outputCh chan client.OutputEvent) *Model {
	ta := textarea.New()
	ta.Placeholder = "Type a message... (Enter to Send)"
	ta.SetWidth(80)
	ta.SetHeight(5)
	ta.Focus()
	ta.CharLimit = 0
	ta.ShowLineNumbers = false
	ta.KeyMap.InsertNewline = key.NewBinding(key.WithKeys("shift+enter", "enter"))

	vp := viewport.New(viewport.WithWidth(80), viewport.WithHeight(20))

	styles := textarea.DefaultDarkStyles()
	styles.Focused.Base = styles.Focused.Base.Background(theme.ThemeInputBg)
	styles.Blurred.Base = styles.Blurred.Base.Background(theme.ThemeInputBg)
	ta.SetStyles(styles)

	return &Model{
		logger:       logger,
		inputCh:      inputCh,
		outputCh:     outputCh,
		cmd:          cmd,
		ctx:          ctx,
		cancel:       cancel,
		textarea:     ta,
		chatViewport: vp,
		statusText:   "Ready",
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.updateSizes()
		return m, nil

	case tea.KeyMsg:
		return m.handleKey(msg)

	case tea.MouseMsg:
		return m.handleMouse(msg)

	case outputEventMsg:
		return m.handleOutput(msg.event)
	}
	return m, nil
}

func (m *Model) View() tea.View {
	if m.width == 0 || m.height == 0 {
		return tea.NewView("Initializing...")
	}

	chat := m.chatViewport.View()
	input := m.renderInput()
	status := theme.StatusBar().
		Width(m.width - 4).
		Render(m.statusText)

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		chat,
		"\n"+input,
		status,
	)

	view := tea.NewView(theme.PureBlack().
		Width(m.width).
		Height(m.height).
		Padding(0, 2).
		Render(content))
	view.AltScreen = true
	view.MouseMode = tea.MouseModeAllMotion
	return view
}

func (m *Model) renderInput() string {
	m.textarea.Prompt = theme.AccentStyle().Render("┃ ")
	return m.textarea.View()
}

func (m *Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		if m.promptRunning {
			return m, nil
		}
		text := m.textarea.Value()
		if text == "" {
			return m, nil
		}
		m.textarea.Reset()
		m.messages = append(m.messages, component.Message{Role: "user", Content: text})
		m.promptRunning = true
		m.loading = true
		m.statusText = "Processing..."
		m.inputCh <- client.InputCommand{Type: client.CmdPrompt, Text: text}
		m.chatViewport.SetContent(m.renderMessages())
		return m, waitForOutput(m.outputCh)

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
	switch msg := msg.(type) {
	case tea.MouseWheelMsg:
		var cmd tea.Cmd
		m.chatViewport, cmd = m.chatViewport.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m *Model) handleOutput(ev client.OutputEvent) (tea.Model, tea.Cmd) {
	if ev.Update != nil {
		update := ev.Update.Update
		if update.AgentMessageChunk != nil && update.AgentMessageChunk.Content.Text != nil {
			text := update.AgentMessageChunk.Content.Text.Text
			if len(m.messages) > 0 && m.messages[len(m.messages)-1].Role == "agent" {
				m.messages[len(m.messages)-1].Content += text
			} else {
				m.messages = append(m.messages, component.Message{Role: "agent", Content: text})
			}
		}
		if update.AgentThoughtChunk != nil && update.AgentThoughtChunk.Content.Text != nil {
			text := update.AgentThoughtChunk.Content.Text.Text
			if len(m.messages) > 0 && m.messages[len(m.messages)-1].Role == "thought" {
				m.messages[len(m.messages)-1].Content += text
			} else {
				m.messages = append(m.messages, component.Message{Role: "thought", Content: text})
			}
		}
		if update.ToolCall != nil {
			tc := update.ToolCall
			label := tc.Title
			input, _ := json.Marshal(tc.RawInput)
			m.messages = append(m.messages, component.Message{Role: "tool", Content: label + "\n" + string(input)})
		}
		if update.ToolCallUpdate != nil {
			tu := update.ToolCallUpdate
			status := "completed"
			if tu.Status != nil {
				status = string(*tu.Status)
			}
			output := ""
			if tu.RawOutput != nil {
				output = fmt.Sprintf("%v", tu.RawOutput)
			}
			if output != "" {
				m.messages = append(m.messages, component.Message{Role: "result", Content: status + ": " + output})
			}
		}
		if update.Plan != nil {
			var lines []string
			for _, e := range update.Plan.Entries {
				mark := " "
				if e.Status == acp.PlanEntryStatusCompleted {
					mark = "✓"
				} else if e.Status == acp.PlanEntryStatusInProgress {
					mark = "→"
				}
				lines = append(lines, fmt.Sprintf("[%s] %s", mark, e.Content))
			}
			m.messages = append(m.messages, component.Message{Role: "plan", Content: strings.Join(lines, "\n")})
		}
		m.chatViewport.SetContent(m.renderMessages())
		m.chatViewport.GotoBottom()
		return m, waitForOutput(m.outputCh)
	}

	switch ev.Kind {
	case "done":
		m.promptRunning = false
		m.loading = false
		m.statusText = "Ready"
		m.chatViewport.SetContent(m.renderMessages())
		m.chatViewport.GotoBottom()
	case "error":
		m.promptRunning = false
		m.loading = false
		m.statusText = "Error: " + ev.Error.Error()
	}
	return m, nil
}

func (m *Model) updateSizes() {
	chatW := m.width - 4
	chatH := m.height - 8
	m.chatViewport.SetWidth(chatW)
	m.chatViewport.SetHeight(chatH)
	m.chatViewport.Style = lipgloss.NewStyle()

	inputW := m.width - 4
	m.textarea.SetWidth(inputW)
	m.textarea.SetHeight(5)
}

func (m *Model) renderMessages() string {
	w := m.chatViewport.Width()
	var content string
	for _, msg := range m.messages {
		content += msg.Render(w)
	}
	return content
}

func (m *Model) cleanup() {
	if m.cmd != nil && m.cmd.Process != nil {
		_ = m.cmd.Process.Kill()
	}
	m.cancel()
}

func waitForOutput(ch chan client.OutputEvent) tea.Cmd {
	return func() tea.Msg {
		ev, ok := <-ch
		if !ok {
			return nil
		}
		return outputEventMsg{event: ev}
	}
}

type outputEventMsg struct {
	event client.OutputEvent
}

var _ tea.Model = (*Model)(nil)
