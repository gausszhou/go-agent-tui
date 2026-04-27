package client

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/coder/acp-go-sdk"
)

type CommandType int

const (
	CmdPrompt   CommandType = iota
	CmdInterrupt
	CmdNewSession
	CmdLoadSession
)

type InputCommand struct {
	Type      CommandType
	SessionID string
	Text      string
	CWD       string
}

type EventKind int

const (
	EventPromptDone     EventKind = iota
	EventSessionCreated
	EventSessionLoaded
	EventPermission
	EventError
)

type PermissionRequest struct {
	Req      acp.RequestPermissionRequest
	Response chan<- acp.RequestPermissionResponse
}

type OutputEvent struct {
	Kind       EventKind
	Update     *acp.SessionUpdate
	Permission *PermissionRequest
	SessionID  string
	Error      error
}

type ACPClient struct {
	Input  <-chan InputCommand
	Output chan<- OutputEvent
	Signal chan<- interface{}
}

var _ acp.Client = (*ACPClient)(nil)

func NewClient(input <-chan InputCommand, output chan<- OutputEvent, signal chan<- interface{}) *ACPClient {
	return &ACPClient{
		Input:  input,
		Output: output,
		Signal: signal,
	}
}

func (c *ACPClient) Run(ctx context.Context, conn *acp.ClientSideConnection) {
	var promptCancel context.CancelFunc

	for {
		select {
		case cmd, ok := <-c.Input:
			if !ok {
				return
			}
			switch cmd.Type {
			case CmdPrompt:
				pCtx, cancel := context.WithCancel(ctx)
				promptCancel = cancel
				go doPrompt(pCtx, conn, c.Output, cmd)
			case CmdInterrupt:
				if promptCancel != nil {
					promptCancel()
					promptCancel = nil
				}
			case CmdNewSession:
				go doNewSession(ctx, conn, c.Output, cmd)
			case CmdLoadSession:
				go doLoadSession(ctx, conn, c.Output, cmd)
			}
		case <-ctx.Done():
			return
		}
	}
}

func (c *ACPClient) RequestPermission(ctx context.Context, params acp.RequestPermissionRequest) (acp.RequestPermissionResponse, error) {
	respCh := make(chan acp.RequestPermissionResponse, 1)
	c.Output <- OutputEvent{
		Kind: EventPermission,
		Permission: &PermissionRequest{
			Req:      params,
			Response: respCh,
		},
	}
	select {
	case resp := <-respCh:
		return resp, nil
	case <-ctx.Done():
		return acp.RequestPermissionResponse{}, ctx.Err()
	}
}

func (c *ACPClient) SessionUpdate(ctx context.Context, params acp.SessionNotification) error {
	c.Output <- OutputEvent{Update: &params.Update}
	return nil
}

func (c *ACPClient) WriteTextFile(ctx context.Context, params acp.WriteTextFileRequest) (acp.WriteTextFileResponse, error) {
	if !filepath.IsAbs(params.Path) {
		return acp.WriteTextFileResponse{}, fmt.Errorf("path must be absolute: %s", params.Path)
	}
	dir := filepath.Dir(params.Path)
	if dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return acp.WriteTextFileResponse{}, fmt.Errorf("mkdir %s: %w", dir, err)
		}
	}
	if err := os.WriteFile(params.Path, []byte(params.Content), 0o644); err != nil {
		return acp.WriteTextFileResponse{}, fmt.Errorf("write %s: %w", params.Path, err)
	}
	return acp.WriteTextFileResponse{}, nil
}

func (c *ACPClient) ReadTextFile(ctx context.Context, params acp.ReadTextFileRequest) (acp.ReadTextFileResponse, error) {
	if !filepath.IsAbs(params.Path) {
		return acp.ReadTextFileResponse{}, fmt.Errorf("path must be absolute: %s", params.Path)
	}
	b, err := os.ReadFile(params.Path)
	if err != nil {
		return acp.ReadTextFileResponse{}, fmt.Errorf("read %s: %w", params.Path, err)
	}
	content := string(b)
	if params.Line != nil || params.Limit != nil {
		lines := strings.Split(content, "\n")
		start := 0
		if params.Line != nil && *params.Line > 0 {
			start = min(max(*params.Line-1, 0), len(lines))
		}
		end := len(lines)
		if params.Limit != nil && *params.Limit > 0 {
			if start+*params.Limit < end {
				end = start + *params.Limit
			}
		}
		content = strings.Join(lines[start:end], "\n")
	}
	return acp.ReadTextFileResponse{Content: content}, nil
}

func (c *ACPClient) CreateTerminal(ctx context.Context, params acp.CreateTerminalRequest) (acp.CreateTerminalResponse, error) {
	return acp.CreateTerminalResponse{TerminalId: "term-1"}, nil
}

func (c *ACPClient) TerminalOutput(ctx context.Context, params acp.TerminalOutputRequest) (acp.TerminalOutputResponse, error) {
	return acp.TerminalOutputResponse{}, nil
}

func (c *ACPClient) ReleaseTerminal(ctx context.Context, params acp.ReleaseTerminalRequest) (acp.ReleaseTerminalResponse, error) {
	return acp.ReleaseTerminalResponse{}, nil
}

func (c *ACPClient) WaitForTerminalExit(ctx context.Context, params acp.WaitForTerminalExitRequest) (acp.WaitForTerminalExitResponse, error) {
	return acp.WaitForTerminalExitResponse{}, nil
}

func (c *ACPClient) KillTerminal(ctx context.Context, params acp.KillTerminalRequest) (acp.KillTerminalResponse, error) {
	return acp.KillTerminalResponse{}, nil
}

func NewConnection(cmd *exec.Cmd, client acp.Client, logger *slog.Logger) (*acp.ClientSideConnection, error) {
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("stdin pipe: %w", err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("stdout pipe: %w", err)
	}
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("cmd start: %w", err)
	}
	conn := acp.NewClientSideConnection(client, stdin, stdout)
	conn.SetLogger(logger)
	return conn, nil
}

func MustCwd() string {
	wd, err := os.Getwd()
	if err != nil {
		return "."
	}
	return wd
}

func doPrompt(ctx context.Context, conn *acp.ClientSideConnection, output chan<- OutputEvent, cmd InputCommand) {
	_, err := conn.Prompt(ctx, acp.PromptRequest{
		SessionId: acp.SessionId(cmd.SessionID),
		Prompt:    []acp.ContentBlock{acp.TextBlock(cmd.Text)},
	})
	output <- OutputEvent{Kind: EventPromptDone, Error: err}
}

func doNewSession(ctx context.Context, conn *acp.ClientSideConnection, output chan<- OutputEvent, cmd InputCommand) {
	cwd := cmd.CWD
	if cwd == "" {
		cwd = MustCwd()
	}
	newSess, err := conn.NewSession(ctx, acp.NewSessionRequest{
		Cwd:        cwd,
		McpServers: []acp.McpServer{},
	})
	if err != nil {
		output <- OutputEvent{Kind: EventSessionCreated, Error: err}
		return
	}
	output <- OutputEvent{Kind: EventSessionCreated, SessionID: string(newSess.SessionId)}
}

func doLoadSession(ctx context.Context, conn *acp.ClientSideConnection, output chan<- OutputEvent, cmd InputCommand) {
	cwd := cmd.CWD
	if cwd == "" {
		cwd = MustCwd()
	}
	_, err := conn.LoadSession(ctx, acp.LoadSessionRequest{
		SessionId: acp.SessionId(cmd.SessionID),
		Cwd:       cwd,
		McpServers: []acp.McpServer{},
	})
	output <- OutputEvent{Kind: EventSessionLoaded, Error: err}
}
