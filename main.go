package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	acpsdk "github.com/coder/acp-go-sdk"
	"github.com/spf13/cobra"

	"github.com/gausszhou/text-ui-research/client"
	"github.com/gausszhou/text-ui-research/tui"
)

var debug bool

var rootCmd = &cobra.Command{
	Use:   "",
	Short: "A TUI application for interacting with AI agents via ACP protocol",
	RunE: func(cmd *cobra.Command, args []string) error {
		return run(debug)
	},
}

func init() {
	rootCmd.Flags().BoolVar(&debug, "debug", false, "Enable debug logging to local file")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run(debug bool) error {
	logger := setupLogger(debug)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	inCh := make(chan client.InputCommand, 1)
	outCh := make(chan client.OutputEvent, 1)
	sigCh := make(chan interface{}, 1)

	acp := client.NewClient(inCh, outCh, sigCh)

	cmd := exec.CommandContext(ctx, "opencode", "acp")
	conn, err := client.NewConnection(cmd, acp, logger)
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
	if debug {
		logger.Info("agent initialized", "protocol_version", initResp.ProtocolVersion)
	}

	newSess, err := conn.NewSession(ctx, acpsdk.NewSessionRequest{
		Cwd:        client.MustCwd(),
		McpServers: []acpsdk.McpServer{},
	})
	if err != nil {
		return fmt.Errorf("new session failed: %w", err)
	}
	if debug {
		logger.Info("session created", "session_id", newSess.SessionId)
	}

	go acp.Run(ctx, conn)

	model := tui.NewModel(debug, logger, acp, cmd, string(newSess.SessionId), ctx, cancel, inCh, outCh)
	p := tea.NewProgram(model, tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		return err
	}
	return nil
}

func setupLogger(debug bool) *slog.Logger {
	if !debug {
		return slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	}

	logPath := filepath.Join(".", "logs/logger.log")
	f, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	fmt.Fprintf(os.Stderr, "[debug] Logging to %s\n", logPath)
	return slog.New(slog.NewTextHandler(f, &slog.HandlerOptions{Level: slog.LevelDebug}))
}
