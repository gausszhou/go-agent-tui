package tui

import (
	"context"
	"log/slog"
	"os/exec"
	"time"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/textarea"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/gausszhou/bubblecode/client"
	"github.com/gausszhou/bubblecode/tui/component"
	"github.com/gausszhou/bubblecode/tui/layout"
	"github.com/gausszhou/bubblecode/tui/theme"
)

const (
	roleUser    = "user"
	roleAgent   = "agent"
	roleThought = "thought"
	roleTool  = "tool"
	rolePlan  = "plan"
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

	promptRunning bool
	loading       bool
	statusText    string
	spinner       component.Loading

	showCommands bool
}

func NewModel(logger *slog.Logger, cmd *exec.Cmd, _ string, ctx context.Context, cancel context.CancelFunc, inputCh chan client.InputCommand, outputCh chan client.OutputEvent) *Model {
	ta := newTextarea()
	vp := viewport.New(viewport.WithWidth(layout.GetChatWidth(layout.InitWidth)), viewport.WithHeight(layout.InitHeight))

	return &Model{
		logger:       logger,
		inputCh:      inputCh,
		outputCh:     outputCh,
		cmd:          cmd,
		ctx:          ctx,
		cancel:       cancel,
		width:        layout.InitWidth,
		height:       layout.InitHeight,
		textarea:     ta,
		chatViewport: vp,
		statusText:   "Ready",
		spinner:      component.NewLoading(theme.LoadingSpinner()),
	}
}

func (m *Model) Init() tea.Cmd {
	return tea.Batch(waitForOutput(m.outputCh), textarea.Blink, spinnerTick())
}

func (m *Model) refreshChat() {
	m.chatViewport.SetContent(m.renderMessages())
	m.chatViewport.GotoBottom()
}

func (m *Model) updateSizes() {
	m.chatViewport.SetWidth(layout.GetChatWidth(m.width))
	m.chatViewport.SetHeight(layout.GetChatHeight(m.height))
	m.textarea.SetWidth(layout.GetInputWidth(m.width))
	m.textarea.SetHeight(layout.InputHeight)
}

func (m *Model) cleanup() {
	m.cancel()
	if m.cmd != nil && m.cmd.Process != nil {
		_ = m.cmd.Process.Kill()
	}
}

func newTextarea() textarea.Model {
	ta := textarea.New()
	ta.Placeholder = "Type a message... (Enter to Send)"
	ta.SetWidth(layout.GetInputWidth(layout.InitWidth))
	ta.SetHeight(layout.InputHeight)
	ta.Focus()
	ta.CharLimit = 0
	ta.ShowLineNumbers = false
	ta.Prompt = theme.AccentStyle().Render("┃ ")
	ta.KeyMap.InsertNewline = key.NewBinding(key.WithKeys("shift+enter", "enter"))

	s := ta.Styles()
	s.Focused.CursorLine = lipgloss.NewStyle()
	s.Blurred.CursorLine = lipgloss.NewStyle()
	ta.SetStyles(s)

	return ta
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
