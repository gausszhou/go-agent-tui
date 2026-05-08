package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/term"
)

var (
	blueColor = lipgloss.NewStyle().Foreground(lipgloss.Color("39"))
	dimStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	userStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("226")).Bold(true)
	aiStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("255"))
)

type message struct {
	role    string
	content string
}

type model struct {
	input    textinput.Model
	messages []message
	width    int
	height   int
	quitting bool
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m *model) addMessage(role, content string) {
	m.messages = append(m.messages, message{role: role, content: content})
}

func (m *model) generateAIResponse(userInput string) string {
	input := strings.ToLower(userInput)

	switch {
	case strings.Contains(input, "hello") || strings.Contains(input, "hi"):
		return "你好！我是 Alan Code，有什么可以帮你的吗？"
	case strings.Contains(input, "read") || strings.Contains(input, "读取"):
		return "📖 正在读取文件...\n\n```go\npackage main\n\nfunc main() {\n    fmt.Println(\"Hello World\")\n}\n```"
	case strings.Contains(input, "write") || strings.Contains(input, "写入"):
		return "✍️ 文件已创建成功！"
	case strings.Contains(input, "search") || strings.Contains(input, "搜索"):
		return "🔍 搜索完成！找到 3 处匹配结果。"
	case strings.Contains(input, "test"):
		return "🧪 测试运行完成！\n\n✅ 3/4 测试通过"
	default:
		return fmt.Sprintf("收到: %s\n\n我可以帮你：\n• 读取文件\n• 写入文件\n• 搜索代码\n• 运行测试", userInput)
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.input.Width = msg.Width - 2

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			input := strings.TrimSpace(m.input.Value())
			if input != "" {
				// 添加用户消息
				m.addMessage("user", input)

				// 生成并添加 AI 响应
				aiResponse := m.generateAIResponse(input)
				m.addMessage("assistant", aiResponse)

				// 清空输入框
				m.input.Reset()
			}
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.quitting {
		return "\n再见！👋\n"
	}

	if m.width == 0 {
		return "Loading..."
	}

	divider := dimStyle.Render(strings.Repeat("─", m.width))

	// 构建所有历史消息
	var content strings.Builder

	for _, msg := range m.messages {
		switch msg.role {
		case "user":
			content.WriteString(fmt.Sprintf("%s %s\n\n", userStyle.Render(">"), msg.content))
		case "assistant":
			content.WriteString(fmt.Sprintf("%s %s\n\n", aiStyle.Render("●"), msg.content))
		}
	}

	// 构建输入行
	inputLine := fmt.Sprintf("%s%s", blueColor.Render("> "), m.input.View())

	// 组合：历史消息 + 下横线 + 输入框 + 下横线
	return fmt.Sprintf("%s%s\n%s\n%s",
		content.String(),
		divider,
		inputLine,
		divider,
	)
}

func main() {
	// 获取终端大小
	width, height, err := term.GetSize(os.Stdout.Fd())
	if err != nil {
		width = 80
		height = 24
	}

	p := tea.NewProgram(&model{
		width:    width,
		height:   height,
		messages: []message{},
		input:    textinput.New(),
	})

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
