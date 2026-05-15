package bubblecode

import (
	"log/slog"
	"os"

	"github.com/coder/acp-go-sdk"
	"github.com/spf13/cobra"

	"github.com/gausszhou/text-ui-research/agent"
)

func acpCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "acp",
		Short: "Start as an ACP server (stdio mode)",
		Long:  `Start the mock agent as an ACP server communicating over stdin/stdout.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runACPServer(cmd)
		},
	}

	return cmd
}

func runACPServer(cmd *cobra.Command) error {
	logger := slog.New(slog.NewTextHandler(cmd.ErrOrStderr(), &slog.HandlerOptions{Level: slog.LevelInfo}))

	agentInstance := agent.NewMockAgent(logger)
	conn := acp.NewAgentSideConnection(agentInstance, os.Stdout, os.Stdin)
	agentInstance.SetAgentConnection(conn)

	logger.Info("ACP server started, waiting for connections...")

	<-conn.Done()
	logger.Info("ACP server stopped")

	return nil
}
