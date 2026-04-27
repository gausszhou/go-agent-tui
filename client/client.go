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

type PermissionEvent struct {
	Request  acp.RequestPermissionRequest
	Response chan<- acp.RequestPermissionResponse
}

type AcpEventKind int

const (
	AcpAgentChunk AcpEventKind = iota
	AcpToolCall
	AcpToolUpdate
	AcpPlan
	AcpPromptDone
	AcpError
)

type AcpEvent struct {
	Kind       AcpEventKind
	Text       string
	ToolCall   *acp.SessionUpdateToolCall
	ToolUpdate *acp.SessionToolCallUpdate
	Plan       *acp.SessionUpdatePlan
	Error      error
}

type ACPClient struct {
	eventCh  chan<- AcpEvent
	permCh   chan<- PermissionEvent
	debug    bool
	Logger   *slog.Logger
}

var _ acp.Client = (*ACPClient)(nil)

func NewClient(eventCh chan<- AcpEvent, permCh chan<- PermissionEvent) *ACPClient {
	return &ACPClient{
		eventCh: eventCh,
		permCh:  permCh,
	}
}

func (c *ACPClient) SetDebug(debug bool) {
	c.debug = debug
}

func (c *ACPClient) RequestPermission(ctx context.Context, params acp.RequestPermissionRequest) (acp.RequestPermissionResponse, error) {
	respCh := make(chan acp.RequestPermissionResponse, 1)
	c.permCh <- PermissionEvent{Request: params, Response: respCh}
	select {
	case resp := <-respCh:
		return resp, nil
	case <-ctx.Done():
		return acp.RequestPermissionResponse{}, ctx.Err()
	}
}

func (c *ACPClient) SessionUpdate(ctx context.Context, params acp.SessionNotification) error {
	u := params.Update
	switch {
	case u.AgentMessageChunk != nil:
		content := u.AgentMessageChunk.Content
		if content.Text != nil {
			c.eventCh <- AcpEvent{Kind: AcpAgentChunk, Text: content.Text.Text}
		}
	case u.ToolCall != nil:
		c.eventCh <- AcpEvent{Kind: AcpToolCall, ToolCall: u.ToolCall}
	case u.ToolCallUpdate != nil:
		c.eventCh <- AcpEvent{Kind: AcpToolUpdate, ToolUpdate: u.ToolCallUpdate}
	case u.Plan != nil:
		c.eventCh <- AcpEvent{Kind: AcpPlan, Plan: u.Plan}
	default:
		if c.debug && c.Logger != nil {
			c.Logger.Debug("unhandled update kind", "update", u)
		}
	}
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
	if c.debug && c.Logger != nil {
		c.Logger.Debug("WriteTextFile", "path", params.Path, "bytes", len(params.Content))
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
	if c.debug && c.Logger != nil {
		c.Logger.Debug("ReadTextFile", "path", params.Path, "bytes", len(content))
	}
	return acp.ReadTextFileResponse{Content: content}, nil
}

func (c *ACPClient) CreateTerminal(ctx context.Context, params acp.CreateTerminalRequest) (acp.CreateTerminalResponse, error) {
	if c.debug && c.Logger != nil {
		c.Logger.Debug("CreateTerminal", "params", params)
	}
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
