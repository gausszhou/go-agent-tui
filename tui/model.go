package tui

import (
	"context"
	"log/slog"
	"os/exec"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/coder/acp-go-sdk"

	"github.com/gausszhou/go-agent-tui/client"
	"github.com/gausszhou/go-agent-tui/tui/component"
)

type FocusArea int

const (
	FocusInput FocusArea = iota
	FocusPermission
	FocusSessionList
	FocusCommandPanel
)

type Session struct {
	ID   string
	Name string
	CWD  string
}

type spinnerTickMsg struct {
	time time.Time
}

type renderTickMsg struct{}

type outputEventMsg struct {
	event client.OutputEvent
}

type Model struct {
	width  int
	height int
	debug  bool
	logger *slog.Logger

	acp      *client.ACPClient
	inputCh  chan client.InputCommand
	outputCh chan client.OutputEvent
	cmd      *exec.Cmd
	ctx      context.Context
	cancel   context.CancelFunc

	focus FocusArea

	pendingPerm *client.PermissionRequest
	questionBox component.QuestionBox

	messages []component.ChatMessage

	textarea textarea.Model

	spinner component.Loading
	loading bool

	sessions        []Session
	activeSessionID string
	sessionList     component.SessionList

	usageInfo component.UsageInfo
	todoList  component.TodoList
	statusBar component.StatusBar

	commandPanelIdx int

	chatViewport viewport.Model

	promptRunning bool
	interrupted   bool

	errMsg     string
	statusText string

	lastKeyTime         time.Time
	lastEscTime         time.Time
	viewportFocused     bool
	scrollDragging      bool
	scrollDragStartY    int
	scrollDragStartYOff int
	viewportDirty       bool
}

func NewModel(debug bool, logger *slog.Logger, acp *client.ACPClient, cmd *exec.Cmd, sessionID string, ctx context.Context, cancel context.CancelFunc, inputCh chan client.InputCommand, outputCh chan client.OutputEvent) Model {
	ta := textarea.New()
	ta.Placeholder = "Type a message... (Enter to send, Shift+Enter for newline)"
	ta.SetWidth(80)
	ta.SetHeight(5)
	ta.Focus()
	ta.CharLimit = 0
	ta.ShowLineNumbers = false
	ta.KeyMap.InsertNewline = key.NewBinding(key.WithKeys("shift+enter", "enter"))

	vp := viewport.New(80, 20)

	return Model{
		debug:  debug,
		logger: logger,

		acp:      acp,
		inputCh:  inputCh,
		outputCh: outputCh,
		cmd:      cmd,
		ctx:      ctx,
		cancel:   cancel,

		activeSessionID: sessionID,

		focus: FocusInput,

		textarea: ta,
		spinner:  component.NewLoading(loadingSpinner()),
		loading:  false,

		sessions: []Session{
			{ID: sessionID, Name: "Session 1", CWD: client.MustCwd()},
		},
		sessionList: component.SessionList{
			Sessions: []component.SessionItem{
				{ID: sessionID, Name: "Session 1", Active: true},
			},
		},
		usageInfo:    component.NewUsageInfo(),
		todoList:     component.NewTodoList("Tasks"),
		statusBar:    component.NewStatusBar(),
		chatViewport: vp,

		statusText: "Ready",
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(spinnerTick(), renderTick())
}

func (m *Model) sendInput(cmd client.InputCommand) {
	select {
	case m.inputCh <- cmd:
	default:
	}
}

func (m *Model) appendAgentText(text string) {
	if len(m.messages) > 0 {
		last := &m.messages[len(m.messages)-1]
		if last.Role == component.RoleAgent && last.ToolCallID == "" {
			last.Content += text
			return
		}
	}
	m.messages = append(m.messages, component.ChatMessage{Role: component.RoleAgent, Content: text})
}

func (m *Model) addMessage(msg component.ChatMessage) {
	m.messages = append(m.messages, msg)
}

func (m *Model) renderMessages() string {
	var sb strings.Builder
	for _, msg := range m.messages {
		sb.WriteString(msg.Render(m.chatViewport.Width, userLabel(), agentLabel(), toolLabel(), systemLabel()))
		sb.WriteString("\n")
	}
	return sb.String()
}

func (m *Model) updateChatViewport() {
	m.chatViewport.SetContent(m.renderMessages())
	m.chatViewport.GotoBottom()
}

func (m *Model) interruptPrompt() {
	m.sendInput(client.InputCommand{Type: client.CmdInterrupt})
	m.interrupted = true
	m.statusText = "Interrupted"
}

func (m *Model) cleanup() {
	m.interruptPrompt()
	if m.cmd != nil && m.cmd.Process != nil {
		_ = m.cmd.Process.Kill()
	}
	m.cancel()
}

func spinnerTick() tea.Cmd {
	return tea.Tick(time.Millisecond*80, func(t time.Time) tea.Msg {
		return spinnerTickMsg{time: t}
	})
}

func renderTick() tea.Cmd {
	return tea.Tick(33*time.Millisecond, func(t time.Time) tea.Msg {
		return renderTickMsg{}
	})
}

func makePermissionQuestionBox(req acp.RequestPermissionRequest, width int) component.QuestionBox {
	title := "Permission Required"
	if req.ToolCall.Title != nil {
		title = *req.ToolCall.Title
	}

	message := "Agent is requesting permission to proceed."
	if req.ToolCall.ToolCallId != "" {
		message = "Tool call requires permission: " + string(req.ToolCall.ToolCallId)
	}

	options := make([]string, len(req.Options))
	for i, opt := range req.Options {
		options[i] = opt.Name + " (" + string(opt.Kind) + ")"
	}

	return component.NewQuestionBox(title, message, options, width)
}
