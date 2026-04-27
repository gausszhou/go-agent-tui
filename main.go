package main

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/gausszhou/go-agent-tui/tui"
)

var debug bool

var rootCmd = &cobra.Command{
	Use:   "go-agent-tui",
	Short: "A TUI application for interacting with AI agents via ACP protocol",
	RunE: func(cmd *cobra.Command, args []string) error {
		logger := setupLogger(debug)

		model := tui.NewModel(debug, logger)
		program := tea.NewProgram(
			model,
			tea.WithAltScreen(),
			tea.WithMouseCellMotion(),
		)

		if _, err := program.Run(); err != nil {
			if debug && logger != nil {
				logger.Error("program error", "error", err)
			}
			return err
		}

		return nil
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

func setupLogger(debug bool) *slog.Logger {
	if !debug {
		return slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	}

	logDir := logDir()
	if err := os.MkdirAll(logDir, 0o755); err != nil {
		return slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	logPath := filepath.Join(logDir, "go-agent-tui.log")
	f, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	fmt.Fprintf(os.Stderr, "[debug] Logging to %s\n", logPath)
	return slog.New(slog.NewTextHandler(f, &slog.HandlerOptions{Level: slog.LevelDebug}))
}

func logDir() string {
	switch runtime.GOOS {
	case "windows":
		dir := os.Getenv("LOCALAPPDATA")
		if dir == "" {
			dir = os.Getenv("APPDATA")
		}
		if dir == "" {
			dir, _ = os.UserHomeDir()
		}
		return filepath.Join(dir, "go-agent-tui")
	default:
		dir, _ := os.UserHomeDir()
		return filepath.Join(dir, ".local", "share", "go-agent-tui")
	}
}
