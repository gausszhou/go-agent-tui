# go-agent-tui

TUI for AI agents via the [ACP protocol](https://github.com/coder/acp-go-sdk) (`opencode acp`).

## Build

```bash
go build -o go-agent-tui.exe .        # TUI app
go build -o agent.exe ./agent         # standalone example agent
```

## Run

```bash
go run .                              # requires `opencode` in PATH
go run . --debug                      # log to ~/.local/share/go-agent-tui/ (or %LOCALAPPDATA% on Windows)
go run ./agent                        # example agent, connects via stdio
```

## Architecture

- `main.go` — entrypoint. Spawns `opencode acp` as subprocess, connects via stdin/stdout pipes using `acp-go-sdk`.
- `client/client.go` — implements `acp.Client` interface. Sends/receives commands over channels. Also handles fs ops (ReadTextFile, WriteTextFile) and terminal stubs.
- `tui/` — Bubble Tea v1 app. Left-right layout: chat+input+status (68%) | todo+session sidebar.
- `agent/agent.go` — separate `package main` binary (NOT importable). Demo agent with simulated streaming, tool calls, and permission requests.
- `docs/opencode-design.md` — visual design reference (OpenCode-inspired palette). Mirrored in `tui/styles.go`.

## Key conventions

- **Enter** sends message, **Shift+Enter** inserts newline.
- **Double-Esc** interrupts running prompt.
- **Ctrl+P** opens command overlay, **Ctrl+N** new session, **Ctrl+S** switch session, **Ctrl+C** quit.
- Paste protection: keystrokes <20ms apart are treated as paste and inserted as text, not sent.
- Design: warm dark (`#201d1d`) bg, Berkeley Mono aesthetic, flat surfaces (no shadows), 4px border radius by default (6px for inputs). See `tui/styles.go`.

## Repo quirks

- `agent/` has its own `package main` — do not import from main app.
- No tests, no CI, no lint config, no Makefile.
- `opencode acp` must be discoverable in PATH at runtime.
