package bubblecode

import (
	"context"
	"fmt"
	"log/slog"
	"os/exec"

	tea "charm.land/bubbletea/v2"
	acpsdk "github.com/coder/acp-go-sdk"
	"github.com/spf13/cobra"

	"github.com/gausszhou/text-ui-research/client"
	"github.com/gausszhou/text-ui-research/tui"
)

func runTUI(cmd *cobra.Command) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := slog.New(slog.NewTextHandler(cmd.ErrOrStderr(), &slog.HandlerOptions{Level: slog.LevelInfo}))

	agentCmd := exec.CommandContext(ctx, "agent.exe")

	events := make(chan client.OutputEvent, 100)

	acpClient := client.NewACPClient(events)

	conn, err := client.NewConnection(agentCmd, acpClient, logger)
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

	model := tui.NewModel(logger, agentCmd, string(newSess.SessionId), ctx, cancel, inputCh, events)
	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		return err
	}

	return nil
}
