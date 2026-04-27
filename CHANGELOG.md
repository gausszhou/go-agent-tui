# Changelog

All notable changes to go-agent-tui will be documented in this file.

## [Unreleased]

### Added
- ACP protocol integration via `acp-go-sdk` v0.12.0 with `opencode acp` agent
- TUI application with left-right layout using `bubbletea` v1 and `lipgloss` v1
- Component library: Loading, Button, QuestionBox, MessageBox, TodoList, CommandPanel, StatusBar, UsagePanel, SessionList
- Streaming agent output with tool call and plan tracking
- Permission inquiry overlay with keyboard navigation (QuestionBox)
- Session management: create (Ctrl+N), switch (Ctrl+S), and list sessions
- `--debug` flag for local file logging
- OpenCode-inspired warm dark color palette (`#201d1d` background)
- Mouse wheel scrolling support for chat viewport
- `rmhubbert/bubbletea-overlay` for permission dialog compositing
- Prompt-style input with `❯` symbol (blue when focused, gray when unfocused)
- Ctrl+P commands panel overlay (New Session, Switch Session)
- Time-based paste detection to prevent multi-line paste from sending multiple messages
- Double-Esc to interrupt prompt during conversation
- Session loading via `LoadSession` to restore chat history on session switch
- `UserMessageChunk` handling for displaying loaded user messages

### Changed
- Input area: border replaced with `❯` prompt symbol
- Input behavior: Enter sends, Shift+Enter inserts newline
- UsageInfo panel: removed Tokens and Cost, keep only Model info
- Status bar: fixed left alignment with chat and input area
- Session switching: now calls `LoadSession` and restores history

### Removed
- Ctrl+I (Interrupt) key binding — replaced by double-Esc
- Ctrl+H (Help toggle) key binding

### Fixed
- First user message now appears immediately in chat history after sending
- Input area left alignment with other UI elements
- Multi-line paste no longer sends multiple messages
