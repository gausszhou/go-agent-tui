package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	viewport viewport.Model
	ready    bool
	content  string
	height   int
	width    int
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// 窗口尺寸变化时更新 viewport
		m.height = msg.Height
		m.width = msg.Width
		m.viewport = viewport.New(msg.Width, msg.Height-2) // 留出边距
		m.viewport.YPosition = 1
		m.viewport.SetContent(m.content)
		m.ready = true

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "up", "k":
			m.viewport.LineUp(1)
		case "down", "j":
			m.viewport.LineDown(1)
		case "pgup", "b":
			m.viewport.PageUp()
		case "pgdn", "f", " ":
			m.viewport.PageDown()
		case "home":
			m.viewport.GotoTop()
		case "end":
			m.viewport.GotoBottom()
		}

	case tea.MouseMsg:
		// 处理鼠标滚轮事件
		switch msg.Button {
		case tea.MouseButtonWheelUp:
			m.viewport.ScrollUp(1)
		case tea.MouseButtonWheelDown:
			m.viewport.ScrollDown(1)
		}

	default:
		// 其他消息传递给 viewport
		m.viewport, cmd = m.viewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if !m.ready {
		return "Loading..."
	}

	// 使用 lipgloss 添加边框和标题
	style := lipgloss.NewStyle().
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(0, 1).
		Width(m.width - 2)

	return style.Render(m.viewport.View())
}

func main() {
	// 获取当前源文件所在目录
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		fmt.Println("无法获取当前文件路径")
		return
	}
	dir := filepath.Dir(filename)

	// 构造同目录下的目标文件路径
	filePath := filepath.Join(dir, "data.md")
	buffer, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading markdown file: %v\n", err)
		return
	}
	markdown := string(buffer)
	// 使用 glamour 渲染 Markdown 为带 ANSI 颜色的字符串
	rendered, err := glamour.Render(markdown, "dark")
	if err != nil {
		rendered = markdown // 降级为纯文本
	}

	p := tea.NewProgram(
		model{content: rendered},
		tea.WithAltScreen(),       // 备用屏幕，锁定终端滚动
		tea.WithMouseCellMotion(), // 启用鼠标事件
	)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
