# go-agent-tui

A terminal user interface for interacting with AI agents via the [ACP protocol](https://github.com/coder/acp-go-sdk).

Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) v1 and [Lip Gloss](https://github.com/charmbracelet/lipgloss) v1.

![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go)

## Requirements

- Go 1.22+
- `opencode` CLI in PATH (provides the ACP agent subprocess)

## Build

```bash
go build -o go-agent-tui.exe .        # TUI app
go build -o agent.exe ./agent         # standalone example agent
```

## Usage

```bash
go run .                              # requires `opencode` in PATH
go run . --debug                      # log to ./go-agent-tui.log
go run ./agent                        # run the example agent (stdio)
```

## Key Bindings

| Key | Action |
|---|---|
| `Enter` | Send message |
| `Shift+Enter` | Insert newline |
| `Esc` `Esc` | Interrupt running prompt |
| `Ctrl+P` | Open commands panel |
| `Ctrl+N` | New session |
| `Ctrl+S` | Open session switcher |
| `↑`/`↓` / `k`/`j` | Scroll chat viewport (when focused) |
| `PgUp`/`PgDn` | Page scroll chat |
| `Ctrl+C` | Quit |

In overlays (Commands/Sessions):

| Key | Action |
|---|---|
| `↑`/`↓` / `k`/`j` | Navigate items |
| `Enter` | Confirm selection |
| `Esc` | Cancel / close overlay |

## Architecture

```
main.go → spawns opencode acp subprocess
          ↓ stdin/stdout pipes
    client/ → ACP client (sends/receives commands)
          ↓ channels
       tui/ → Bubble Tea TUI (chat + sidebar)
```

- **Left panel** (68%): chat viewport + input area + status bar
- **Right panel** (32%): todo/task list
- **Overlays**: commands panel (`Ctrl+P`), session list (`Ctrl+S`), permission requests

## Components

| Package | Description |
|---|---|
| `tui/` | Main TUI model, view, update, styles |
| `client/` | ACP protocol client, connection management |
| `tui/component/` | Reusable UI components (TodoList, StatusBar, SessionList, etc.) |
| `agent/` | Standalone example agent binary (separate `package main`) |

## Design

Warm dark palette inspired by OpenCode:

- Background: `#201d1d`
- Surface: `#302c2c`
- Accent: `#007aff` (blue)
- Text: `#fdfcfc`
- Muted: `#9a9898`

Flat surfaces, 4px border radius, Berkeley Mono aesthetic.
