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
	"github.com/gausszhou/text-ui-research/tui/layout"
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
	spinner       component.Loading
}

func NewModel(logger *slog.Logger, cmd *exec.Cmd, sessionID string, ctx context.Context, cancel context.CancelFunc, inputCh chan client.InputCommand, outputCh chan client.OutputEvent) *Model {
	initW := 80
	ta := textarea.New()
	ta.Placeholder = "Type a message... (Enter to Send)"
	ta.SetWidth(layout.GetInputWidth(initW))
	ta.SetHeight(layout.InputHeight)
	ta.Focus()
	ta.CharLimit = 0
	ta.ShowLineNumbers = false
	ta.KeyMap.InsertNewline = key.NewBinding(key.WithKeys("shift+enter", "enter"))
	ta.Prompt = theme.AccentStyle().Render("┃ ")

	vp := viewport.New(viewport.WithWidth(layout.GetChatWidth(initW)), viewport.WithHeight(20))

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
		width:        initW,
		height:       24,
		textarea:     ta,
		chatViewport: vp,
		statusText:   "Ready",
		spinner:      component.NewLoading(theme.LoadingSpinner()),
	}
}

func (m *Model) Init() tea.Cmd {
	return tea.Batch(waitForOutput(m.outputCh), spinnerTick())
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if msg.Width < layout.MinWidth || msg.Height < layout.MinHeight {
			return m, nil
		}
		m.width = msg.Width
		m.height = msg.Height
		m.updateSizes()
		m.chatViewport.SetContent(m.renderMessages())
		return m, nil

	case tea.KeyMsg:
		return m.handleKey(msg)

	case tea.MouseMsg:
		return m.handleMouse(msg)

	case outputEventMsg:
		return m.handleOutput(msg.event)

	case loadingTickMsg:
		m.spinner = m.spinner.Tick()
		return m, spinnerTick()
	}
	return m, nil
}

func (m *Model) View() tea.View {
	chat := m.chatViewport.View()
	input := m.renderInput()

	left := m.statusText
	if m.loading {
		left = m.spinner.View() + " " + left
	} else {
		left = "✓ " + left
	}
	pad := 2 * layout.PaddingHorizontal
	status := theme.StatusBar().
		Width(m.width - pad).
		Render(left)

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		chat,
		"\n"+input,
		status,
	)

	view := tea.NewView(theme.PureBlack().
		Width(m.width).
		Height(m.height).
		Padding(0, layout.PaddingHorizontal).
		Render(content))
	view.AltScreen = true
	view.MouseMode = tea.MouseModeAllMotion
	return view
}

func (m *Model) renderInput() string {
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
		m.chatViewport.SetContent(m.renderMessages())
		return m, tea.Batch(sendInput(m.inputCh, client.InputCommand{Type: client.CmdPrompt, Text: text}), waitForOutput(m.outputCh), spinnerTick())

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
		switch msg.Button {
		case tea.MouseWheelUp:
			m.chatViewport.ScrollUp(3)
		case tea.MouseWheelDown:
			m.chatViewport.ScrollDown(3)
		}
		return m, nil
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
	m.chatViewport.SetWidth(layout.GetChatWidth(m.width))
	m.chatViewport.SetHeight(layout.GetChatHeight(m.height))
	m.textarea.SetWidth(layout.GetInputWidth(m.width))
	m.textarea.SetHeight(layout.InputHeight)
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

type outputEventMsg struct {
	event client.OutputEvent
}

type channelClosedMsg struct{}

type loadingTickMsg struct{}

type inputSentMsg struct{}

func sendInput(ch chan client.InputCommand, cmd client.InputCommand) tea.Cmd {
	return func() tea.Msg {
		ch <- cmd
		return inputSentMsg{}
	}
}

func waitForOutput(ch chan client.OutputEvent) tea.Cmd {
	return func() tea.Msg {
		ev, ok := <-ch
		if !ok {
			return channelClosedMsg{}
		}
		return outputEventMsg{event: ev}
	}
}

func spinnerTick() tea.Cmd {
	return tea.Tick(80*time.Millisecond, func(t time.Time) tea.Msg {
		return loadingTickMsg{}
	})
}

var _ tea.Model = (*Model)(nil)
