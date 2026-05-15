package main

import (
	"log/slog"
	"os"

	"github.com/coder/acp-go-sdk"

	"github.com/gausszhou/bubblecode/agent"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo}))
	a := agent.NewMockAgent(logger)
	conn := acp.NewAgentSideConnection(a, os.Stdout, os.Stdin)
	a.SetAgentConnection(conn)
	<-conn.Done()
}
