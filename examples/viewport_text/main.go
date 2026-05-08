package main

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// 定义各种样式（无背景色）
var (
	// 标题样式：粗体 + 亮青色
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("51")) // 亮青色

	// 强调样式：粗体 + 黄色
	emphasisStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("226")) // 黄色

	// 错误样式：红色 + 斜体（部分终端支持）
	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")). // 亮红色
			Italic(true)

	// 成功样式：绿色
	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("46")) // 亮绿色

	// 警告样式：橙色
	warningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("208")) // 橙色

	// 代码样式：灰色 + 斜体
	codeStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("245")). // 中灰色
			Italic(true)

	// 链接样式：下划线 + 蓝色
	linkStyle = lipgloss.NewStyle().
			Underline(true).
			Foreground(lipgloss.Color("39")) // 亮蓝色

	// 元数据样式：暗灰色
	metaStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")) // 暗灰色

	// 普通文本样式：白色
	normalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("255")) // 白色
)

type model struct {
	viewport   viewport.Model
	ready      bool
	logEntries []string
}

func newModel() model {
	return model{
		logEntries: generateLogEntries(),
	}
}

// 生成示例日志条目（多种样式）
func generateLogEntries() []string {
	entries := []string{}

	// 标题
	entries = append(entries, titleStyle.Render("=== 系统运行日志 ==="))
	entries = append(entries, "")

	// 成功消息
	entries = append(entries, successStyle.Render("✓ 服务启动成功")+" "+metaStyle.Render("[2024-01-15 10:00:01]"))

	// 普通信息
	entries = append(entries, normalStyle.Render("正在加载配置文件...")+" "+metaStyle.Render("[10:00:02]"))

	// 警告消息
	entries = append(entries, warningStyle.Render("⚠ 配置文件中有未使用的字段")+" "+metaStyle.Render("[10:00:03]"))

	// 带强调的文本
	entries = append(entries, normalStyle.Render(fmt.Sprintf("已加载 %s 个模块", emphasisStyle.Render("42"))))

	// 代码示例
	entries = append(entries, "")
	entries = append(entries, codeStyle.Render("func main() {"))
	entries = append(entries, codeStyle.Render("    fmt.Println(\"Hello, World!\")"))
	entries = append(entries, codeStyle.Render("}"))
	entries = append(entries, "")

	// 错误消息
	entries = append(entries, errorStyle.Render("✗ 数据库连接失败: connection refused")+" "+metaStyle.Render("[10:00:05]"))

	// 带链接的文本
	entries = append(entries, normalStyle.Render("更多信息请访问: ")+linkStyle.Render("https://example.com/docs"))

	// 混合样式
	entries = append(entries, "")
	entries = append(entries, titleStyle.Render("=== 性能指标 ==="))

	metrics := []string{
		fmt.Sprintf("CPU: %s", successStyle.Render("12%")),
		fmt.Sprintf("内存: %s", warningStyle.Render("78%")),
		fmt.Sprintf("磁盘: %s", normalStyle.Render("45%")),
		fmt.Sprintf("网络: %s", errorStyle.Render("超时")),
	}

	for _, m := range metrics {
		entries = append(entries, "  "+m)
	}

	// 成功消息
	entries = append(entries, "")
	entries = append(entries, successStyle.Render("✓ 所有检查完成")+" "+metaStyle.Render("[10:00:10]"))

	return entries
}

func (m model) Init() tea.Cmd {
	return nil
}

// 构建 Viewport 内容（关键：先应用样式，后设置到 viewport）
func (m model) buildContent() string {
	var builder strings.Builder

	for i, entry := range m.logEntries {
		// 如果 entry 是空字符串，直接添加换行
		if entry == "" {
			builder.WriteString("\n")
			continue
		}

		// entry 已经包含了样式（在 generateLogEntries 中已应用）
		// 这里只需要原样写入
		builder.WriteString(entry)

		// 最后一行不加换行
		if i < len(m.logEntries)-1 {
			builder.WriteString("\n")
		}
	}

	return builder.String()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if !m.ready {
			// 初始化 viewport
			m.viewport = viewport.New(viewport.WithWidth(msg.Width), viewport.WithHeight(msg.Height-2)) // 留2行给帮助信息
			m.viewport.YPosition = 0
			m.ready = true
		} else {
			m.viewport.SetWidth(msg.Width)
			m.viewport.SetHeight(msg.Height - 2)
		}

		// ✅ 关键：设置内容到 viewport（内容已经带有样式）
		m.viewport.SetContent(m.buildContent())

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "up", "k":
			m.viewport.ScrollUp(1)
		case "down", "j":
			m.viewport.ScrollDown(1)
		case "pgup", "ctrl+b":
			m.viewport.HalfPageUp()
		case "pgdown", "ctrl+f":
			m.viewport.HalfPageDown()
		case "home", "g":
			m.viewport.GotoTop()
		case "end", "G":
			m.viewport.GotoBottom()
		}
	}

	// 处理 viewport 自己的消息
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() tea.View {
	if !m.ready {
		return tea.NewView("\n  Loading...")
	}

	// 帮助信息样式（无背景色）
	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))

	helpText := helpStyle.Render(
		" ↑/↓ 滚动  •  PgUp/PgDown 翻页  •  g/G 顶部/底部  •  q 退出 ",
	)

	// 组合：viewport + 帮助栏
	return tea.NewView(lipgloss.JoinVertical(
		lipgloss.Top,
		m.viewport.View(),
		helpText,
	))
}

func main() {
	p := tea.NewProgram(newModel())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
	}
}
