package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"

	tea "charm.land/bubbletea/v2"
	acpsdk "github.com/coder/acp-go-sdk"

	"github.com/gausszhou/text-ui-research/client"
	"github.com/gausszhou/text-ui-research/tui"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo}))

	cmd := exec.CommandContext(ctx, "agent.exe")

	events := make(chan client.OutputEvent, 100)

	// 必须先创建 acpClient，NewConnection 需要它作为 acp.Client 实现
	acpClient := client.NewACPClient(events)

	conn, err := client.NewConnection(cmd, acpClient, logger)
	if err != nil {
		return fmt.Errorf("connection failed: %w", err)
	}

	initResp, err := conn.Initialize(ctx, acpsdk.InitializeRequest{
		ProtocolVersion: acpsdk.ProtocolVersionNumber,
		ClientCapabilities: acpsdk.ClientCapabilities{
			Fs:       acpsdk.FileSystemCapabilities{ReadTextFile: true, WriteTextFile: true},
			Terminal: true,
		},
	})
	if err != nil {
		return fmt.Errorf("initialize failed: %w", err)
	}
	logger.Info("agent initialized", "protocol_version", initResp.ProtocolVersion)

	newSess, err := conn.NewSession(ctx, acpsdk.NewSessionRequest{
		Cwd:        client.MustCwd(),
		McpServers: []acpsdk.McpServer{},
	})
	if err != nil {
		return fmt.Errorf("new session failed: %w", err)
	}
	logger.Info("session created", "session_id", newSess.SessionId)

	inputCh := make(chan client.InputCommand, 1)

	cl := client.NewClient(inputCh, conn.ClientConn(), events)
	go cl.Run(ctx, newSess.SessionId)

	model := tui.NewModel(logger, cmd, string(newSess.SessionId), ctx, cancel, inputCh, events)
	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		return err
	}

	return nil
}
