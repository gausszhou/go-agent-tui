package tui

import (
	"context"
	"log/slog"
	"os/exec"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/coder/acp-go-sdk"

	"github.com/gausszhou/go-agent-tui/client"
	"github.com/gausszhou/go-agent-tui/tui/component"
)

type FocusArea int

const (
	FocusInput     FocusArea = iota
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

type outputEventMsg struct {
	event client.OutputEvent
}

type initDoneMsg struct {
	sessionID string
	err       error
	acpClient *client.ACPClient
	cmd       *exec.Cmd
}

type Model struct {
	width  int
	height int
	debug  bool
	logger *slog.Logger

	inputCh  chan client.InputCommand
	outputCh chan client.OutputEvent

	acpClient *client.ACPClient
	cmd       *exec.Cmd
	ctx       context.Context
	cancel    context.CancelFunc

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
	showSessionList bool

	usageInfo component.UsageInfo
	todoList  component.TodoList
	statusBar component.StatusBar

	commandPanelIdx int

	chatViewport viewport.Model

	promptRunning bool

	errMsg     string
	statusText string

	lastKeyTime         time.Time
	lastEscTime         time.Time
	viewportFocused     bool
	scrollDragging      bool
	scrollDragStartY    int
	scrollDragStartYOff int
}

func NewModel(debug bool, logger *slog.Logger) Model {
	ctx, cancel := context.WithCancel(context.Background())

	ta := textarea.New()
	ta.Placeholder = "Type a message... (Enter to send, Shift+Enter for newline)"
	ta.SetWidth(80)
	ta.SetHeight(3)
	ta.Focus()
	ta.CharLimit = 0
	ta.ShowLineNumbers = false
	ta.KeyMap.InsertNewline = key.NewBinding(key.WithKeys("shift+enter", "enter"))

	vp := viewport.New(80, 20)

	inCh := make(chan client.InputCommand, 100)
	outCh := make(chan client.OutputEvent, 100)

	return Model{
		debug:  debug,
		logger: logger,

		ctx:    ctx,
		cancel: cancel,

		inputCh:  inCh,
		outputCh: outCh,

		focus: FocusInput,

		textarea: ta,
		spinner:  component.NewLoading(loadingSpinner()),
		loading:  true,

		sessions:     []Session{},
		sessionList:  component.NewSessionList("Sessions"),
		usageInfo:    component.NewUsageInfo(),
		todoList:     component.NewTodoList("Tasks"),
		statusBar:    component.NewStatusBar(),
		chatViewport: vp,

		statusText: "Connecting...",
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		func() tea.Msg {
			in := m.inputCh
			out := m.outputCh
			acpClient := client.NewClient(in, out, nil)

			if m.debug && m.logger != nil {
				m.logger.Info("starting agent process")
			}

			cmd := exec.CommandContext(m.ctx, "opencode", "acp")
			conn, err := client.NewConnection(cmd, acpClient, m.logger)
			if err != nil {
				return initDoneMsg{err: err}
			}

			initResp, err := conn.Initialize(m.ctx, acp.InitializeRequest{
				ProtocolVersion: acp.ProtocolVersionNumber,
				ClientCapabilities: acp.ClientCapabilities{
					Fs:       acp.FileSystemCapabilities{ReadTextFile: true, WriteTextFile: true},
					Terminal: true,
				},
			})
			if err != nil {
				return initDoneMsg{err: err}
			}
			if m.debug && m.logger != nil {
				m.logger.Info("agent initialized", "protocol_version", initResp.ProtocolVersion)
			}

			newSess, err := conn.NewSession(m.ctx, acp.NewSessionRequest{
				Cwd:        client.MustCwd(),
				McpServers: []acp.McpServer{},
			})
			if err != nil {
				return initDoneMsg{err: err}
			}
			if m.debug && m.logger != nil {
				m.logger.Info("session created", "session_id", newSess.SessionId)
			}

			go acpClient.Run(m.ctx, conn)

			return initDoneMsg{
				sessionID: string(newSess.SessionId),
				acpClient: acpClient,
				cmd:       cmd,
			}
		},
		spinnerTick(),
	)
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
