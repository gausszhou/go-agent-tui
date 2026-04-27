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
)

type Session struct {
	ID   string
	Name string
	CWD  string
}

type spinnerTickMsg struct {
	time time.Time
}

type acpEventMsg struct {
	event client.AcpEvent
}

type permEventMsg struct {
	event client.PermissionEvent
}

type initConnMsg struct {
	conn    *acp.ClientSideConnection
	cmd     *exec.Cmd
	acpClient *client.ACPClient
	err     error
}

type initDoneMsg struct {
	sessionID  string
	err        error
}

type promptDoneMsg struct {
	err error
}

type sessionCreatedMsg struct {
	sessionID string
	err       error
}

type Model struct {
	width  int
	height int
	debug  bool
	logger *slog.Logger

	acpEventCh  chan client.AcpEvent
	permEventCh chan client.PermissionEvent

	acpClient *client.ACPClient
	conn      *acp.ClientSideConnection
	cmd       *exec.Cmd
	ctx       context.Context
	cancel    context.CancelFunc

	focus FocusArea

	pendingPerm *client.PermissionEvent
	questionBox component.QuestionBox

	messages []component.ChatMessage

	textarea textarea.Model

	spinner component.Loading
	loading bool

	sessions          []Session
	activeSessionID   string
	sessionList       component.SessionList
	showSessionList   bool

	usageInfo component.UsageInfo
	todoList  component.TodoList
	statusBar component.StatusBar

	showHelp     bool
	commandPanel component.CommandPanel

	chatViewport viewport.Model

	promptCancel context.CancelFunc
	promptRunning bool

	errMsg      string
	statusText  string
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
	ta.KeyMap.InsertNewline = key.NewBinding(key.WithKeys("shift+enter"))

	vp := viewport.New(80, 20)

	return Model{
		debug:  debug,
		logger: logger,

		ctx:    ctx,
		cancel: cancel,
		conn:   nil,
		cmd:    nil,

		acpEventCh:  make(chan client.AcpEvent, 100),
		permEventCh: make(chan client.PermissionEvent, 10),

		focus: FocusInput,

		textarea: ta,
		spinner:  component.NewLoading(loadingSpinner()),
		loading:  true,

		sessions:       []Session{},
		sessionList:    component.NewSessionList("Sessions"),
		usageInfo:      component.NewUsageInfo(),
		todoList:       component.NewTodoList("Tasks"),
		statusBar:      component.NewStatusBar(),
		commandPanel:   component.DefaultCommands(),
		chatViewport:   vp,

		statusText: "Connecting...",
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		func() tea.Msg {
			cmd := exec.CommandContext(m.ctx, "opencode", "acp")
			acpClient := client.NewClient(m.acpEventCh, m.permEventCh)
			acpClient.Logger = m.logger
			acpClient.SetDebug(m.debug)

			if m.debug && m.logger != nil {
				m.logger.Info("starting agent process", "cmd", cmd.String())
			}

			conn, err := client.NewConnection(cmd, acpClient, m.logger)
			if err != nil {
				return initConnMsg{err: err}
			}
			return initConnMsg{conn: conn, cmd: cmd, acpClient: acpClient}
		},
		spinnerTick(),
	)
}

func (m *Model) doInitialize(conn *acp.ClientSideConnection, acpClient *client.ACPClient) tea.Cmd {
	m.conn = conn
	m.acpClient = acpClient

	return func() tea.Msg {
		initResp, err := m.conn.Initialize(m.ctx, acp.InitializeRequest{
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

		newSess, err := m.conn.NewSession(m.ctx, acp.NewSessionRequest{
			Cwd:        client.MustCwd(),
			McpServers: []acp.McpServer{},
		})
		if err != nil {
			return initDoneMsg{err: err}
		}
		if m.debug && m.logger != nil {
			m.logger.Info("session created", "session_id", newSess.SessionId)
		}

		return initDoneMsg{sessionID: string(newSess.SessionId)}
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

func (m *Model) submitPrompt(text string) tea.Cmd {
	m.addMessage(component.ChatMessage{Role: component.RoleUser, Content: text})

	sessionID := m.activeSessionID
	ctx, cancel := context.WithCancel(m.ctx)
	m.promptCancel = cancel
	m.promptRunning = true
	m.loading = true
	m.errMsg = ""
	m.statusText = "Processing..."

	if m.debug && m.logger != nil {
		m.logger.Info("submitting prompt", "session", sessionID, "text", text)
	}

	return tea.Batch(
		func() tea.Msg {
			_, err := m.conn.Prompt(ctx, acp.PromptRequest{
				SessionId: acp.SessionId(sessionID),
				Prompt:    []acp.ContentBlock{acp.TextBlock(text)},
			})
			return promptDoneMsg{err: err}
		},
		m.waitForEvents(),
		spinnerTick(),
	)
}

func (m *Model) createSession() tea.Cmd {
	m.loading = true
	m.statusText = "Creating session..."
	return func() tea.Msg {
		newSess, err := m.conn.NewSession(m.ctx, acp.NewSessionRequest{
			Cwd:        client.MustCwd(),
			McpServers: []acp.McpServer{},
		})
		if err != nil {
			return sessionCreatedMsg{err: err}
		}
		return sessionCreatedMsg{sessionID: string(newSess.SessionId)}
	}
}

func (m *Model) waitForEvents() tea.Cmd {
	return func() tea.Msg {
		select {
		case evt, ok := <-m.acpEventCh:
			if !ok {
				return nil
			}
			return acpEventMsg{event: evt}
		case evt, ok := <-m.permEventCh:
			if !ok {
				return nil
			}
			return permEventMsg{event: evt}
		}
	}
}

func (m *Model) interruptPrompt() {
	if m.promptCancel != nil {
		m.promptCancel()
		m.promptCancel = nil
		m.statusText = "Interrupted"
	}
}

func (m *Model) cleanup() {
	if m.promptCancel != nil {
		m.promptCancel()
		m.promptCancel = nil
	}
	if m.cmd != nil && m.cmd.Process != nil {
		_ = m.cmd.Process.Kill()
	}
	m.cancel()
}

func (m *Model) renderMessages() string {
	var sb strings.Builder
	for _, msg := range m.messages {
		sb.WriteString(msg.Render(m.chatViewport.Width, userLabel(), agentLabel(), toolLabel(), systemLabel()))
		sb.WriteString("\n")
	}
	return sb.String()
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
