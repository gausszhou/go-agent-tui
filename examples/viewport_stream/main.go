package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// 定义各种样式（无背景色）
var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("51"))

	emphasisStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("226"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196"))

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("46"))

	warningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("208"))

	codeStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("245"))

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("39"))

	normalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("255"))
)

type MessageType int

const (
	MsgInfo MessageType = iota
	MsgSuccess
	MsgWarning
	MsgError
	MsgCode
	MsgEmphasis
)

type StreamMessage struct {
	Type    MessageType
	Content string
}

type model struct {
	viewport         viewport.Model
	ready            bool
	allContent       []string
	streamActive     bool
	messages         []StreamMessage
	currentMsg       *StreamMessage
	charIndex        int
	needsRefresh     bool
	lastContentHash  string
	refreshCountdown int
	// 新增：标记是否需要重新初始化
	needReset bool
}

// 生成随机消息队列
func generateRandomMessages() []StreamMessage {
	messages := []StreamMessage{}

	messages = append(messages, StreamMessage{
		Type:    MsgInfo,
		Content: "=== 系统实时日志流 ===",
	})

	// 生成 50-80 条消息
	numMessages := rand.Intn(30) + 50

	for i := 0; i < numMessages; i++ {
		msgType := MessageType(rand.Intn(6))

		var content string
		switch msgType {
		case MsgInfo:
			contents := []string{
				"正在处理请求...",
				"加载配置文件",
				"初始化模块",
				"建立连接",
				"解析数据包",
				"验证证书",
				"缓存命中",
				"查询数据库",
				"渲染模板",
				"压缩资源",
				"分配内存",
				"启动服务",
			}
			content = contents[rand.Intn(len(contents))]
		case MsgSuccess:
			contents := []string{
				"操作成功",
				"数据已保存",
				"连接建立成功",
				"文件已上传",
				"任务完成",
				"验证通过",
				"同步完成",
				"部署成功",
			}
			content = contents[rand.Intn(len(contents))]
		case MsgWarning:
			contents := []string{
				"磁盘空间不足",
				"网络延迟较高",
				"重试请求",
				"缓存过期",
				"配置已弃用",
				"连接池满",
				"降级处理",
				"响应缓慢",
			}
			content = contents[rand.Intn(len(contents))]
		case MsgError:
			contents := []string{
				"连接超时",
				"权限不足",
				"文件不存在",
				"解析失败",
				"服务不可用",
				"校验失败",
				"资源冲突",
				"内存溢出",
			}
			content = contents[rand.Intn(len(contents))]
		case MsgCode:
			contents := []string{
				"func handleRequest() {",
				"    return result, nil",
				"}",
				"if err != nil {",
				"    log.Fatal(err)",
				"}",
				"for item := range items {",
				"    process(item)",
				"}",
				"type Response struct {",
				"    Code int",
				"    Data string",
				"}",
			}
			content = contents[rand.Intn(len(contents))]
		case MsgEmphasis:
			contents := []string{
				"重要：请检查配置",
				"注意：版本更新可用",
				"关键指标异常",
				"用户会话已过期",
				"需要人工确认",
				"安全警告",
			}
			content = contents[rand.Intn(len(contents))]
		}

		messages = append(messages, StreamMessage{
			Type:    msgType,
			Content: content,
		})
	}

	messages = append(messages, StreamMessage{
		Type:    MsgSuccess,
		Content: "=== 日志流输出完成，共 " + fmt.Sprintf("%d", len(messages)) + " 条消息 ===",
	})

	return messages
}

// 根据消息类型应用样式
func styleMessage(msg StreamMessage, content string) string {
	switch msg.Type {
	case MsgInfo:
		return infoStyle.Render(content)
	case MsgSuccess:
		return successStyle.Render("✓ " + content)
	case MsgWarning:
		return warningStyle.Render("⚠ " + content)
	case MsgError:
		return errorStyle.Render("✗ " + content)
	case MsgCode:
		return codeStyle.Render(content)
	case MsgEmphasis:
		return emphasisStyle.Render("★ " + content)
	default:
		return normalStyle.Render(content)
	}
}

type streamNextCharMsg struct{}

func startStreaming() tea.Cmd {
	return func() tea.Msg {
		return streamNextCharMsg{}
	}
}

// 重置 model 到初始状态
func (m *model) reset() {
	m.allContent = make([]string, 0)
	m.streamActive = true
	m.messages = generateRandomMessages()
	m.currentMsg = nil
	m.charIndex = 0
	m.needsRefresh = true
	m.refreshCountdown = 0
	m.lastContentHash = ""
	m.needReset = false
}

func newModel() model {
	m := model{
		allContent:       make([]string, 0),
		streamActive:     true,
		messages:         generateRandomMessages(),
		currentMsg:       nil,
		charIndex:        0,
		needsRefresh:     true,
		refreshCountdown: 0,
		needReset:        false,
	}
	return m
}

func (m model) Init() tea.Cmd {
	// 启动流式输出
	return startStreaming()
}

// 计算内容的哈希值，用于判断是否需要刷新
func (m model) contentHash() string {
	return strings.Join(m.allContent, "\n")
}

// 智能刷新：只在必要时更新 viewport
func (m *model) smartRefresh() {
	if !m.ready {
		return
	}

	currentHash := m.contentHash()
	if currentHash != m.lastContentHash {
		m.viewport.SetContent(currentHash)
		m.lastContentHash = currentHash
		m.viewport.GotoBottom()
	}
	m.needsRefresh = false
}

func (m *model) processStream() tea.Cmd {
	// 如果需要重置，先完成重置
	if m.needReset {
		m.reset()
		// 重置后需要刷新视图
		m.needsRefresh = true
		m.smartRefresh()
	}

	// 如果没有当前消息，从队列中取下一个
	if m.currentMsg == nil {
		if len(m.messages) == 0 {
			m.streamActive = false
			m.needsRefresh = true
			return nil
		}

		msg := m.messages[0]
		m.messages = m.messages[1:]
		m.currentMsg = &msg
		m.charIndex = 0

		// 每条消息之间的延迟
		return tea.Tick(15*time.Millisecond, func(t time.Time) tea.Msg {
			return streamNextCharMsg{}
		})
	}

	// 当前消息还有字符未输出
	if m.charIndex < len(m.currentMsg.Content) {
		// 每次输出 2-4 个字符
		outputCount := rand.Intn(3) + 2
		m.charIndex += outputCount
		if m.charIndex > len(m.currentMsg.Content) {
			m.charIndex = len(m.currentMsg.Content)
		}

		partial := m.currentMsg.Content[:m.charIndex]

		if m.charIndex == outputCount { // 第一次输出这条消息
			styledLine := styleMessage(*m.currentMsg, partial)
			m.allContent = append(m.allContent, styledLine)
			m.needsRefresh = true
		} else { // 更新最后一行
			if len(m.allContent) > 0 {
				styledLine := styleMessage(*m.currentMsg, partial)
				m.allContent[len(m.allContent)-1] = styledLine
				m.needsRefresh = true
			}
		}

		// 批量刷新计数器：每2次字符输出才刷新一次，减少闪烁
		m.refreshCountdown++
		if m.refreshCountdown >= 2 {
			m.smartRefresh()
			m.refreshCountdown = 0
		}

		// 字符间隔
		nextDelay := time.Duration(rand.Intn(10)+5) * time.Millisecond
		return tea.Tick(nextDelay, func(t time.Time) tea.Msg {
			return streamNextCharMsg{}
		})
	}

	// 当前消息完成，确保最后刷新
	m.smartRefresh()
	m.currentMsg = nil
	m.charIndex = 0
	m.needsRefresh = true

	// 消息间延迟
	nextDelay := time.Duration(rand.Intn(20)+10) * time.Millisecond
	return tea.Tick(nextDelay, func(t time.Time) tea.Msg {
		return streamNextCharMsg{}
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if !m.ready {
			m.viewport = viewport.New(msg.Width, msg.Height-3)
			m.viewport.YPosition = 0
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - 3
		}
		m.needsRefresh = true
		m.smartRefresh()

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "up", "k":
			m.viewport.LineUp(1)
		case "down", "j":
			m.viewport.LineDown(1)
		case "pgup", "ctrl+b":
			m.viewport.HalfViewUp()
		case "pgdown", "ctrl+f":
			m.viewport.HalfViewDown()
		case "home", "g":
			m.viewport.GotoTop()
		case "end", "G":
			m.viewport.GotoBottom()
		case "r", "R":
			// ✅ 修复：标记需要重置，而不是创建新 model
			m.needReset = true
			m.streamActive = true
			// 清空当前显示
			m.allContent = []string{}
			m.needsRefresh = true
			m.smartRefresh()
			// 立即开始新的流式输出
			return m, startStreaming()
		case "s", "S":
			// 跳过当前消息
			if m.currentMsg != nil {
				fullContent := m.currentMsg.Content
				styledLine := styleMessage(*m.currentMsg, fullContent)
				if len(m.allContent) > 0 {
					m.allContent[len(m.allContent)-1] = styledLine
				} else {
					m.allContent = append(m.allContent, styledLine)
				}
				m.currentMsg = nil
				m.charIndex = 0
				m.needsRefresh = true
				m.smartRefresh()
				return m, tea.Tick(10*time.Millisecond, func(t time.Time) tea.Msg {
					return streamNextCharMsg{}
				})
			}
		}

	case streamNextCharMsg:
		cmd = m.processStream()
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	// 定期刷新（用于保证最终一致性）
	if m.needsRefresh {
		m.smartRefresh()
	}

	// 处理 viewport 自己的更新
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if !m.ready {
		return "\n  Initializing stream..."
	}

	// 如果正在重置或没有内容且流式输出未激活，显示初始状态
	if m.needReset || (len(m.allContent) == 0 && m.streamActive) {
		return "\n  Starting stream...\n"
	}

	statusStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))

	var statusText string
	if m.streamActive {
		remaining := len(m.messages)
		if m.currentMsg != nil {
			remaining++
		}
		statusText = statusStyle.Render(
			fmt.Sprintf(" ⚡ 流式输出中... (剩余: %d 条)  •  s 跳过当前  •  r 重新开始  •  q 退出 ",
				remaining))
	} else {
		statusText = statusStyle.Render(
			" ✓ 输出完成  •  r 重新开始  •  q 退出 ")
	}

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))

	helpText := helpStyle.Render(
		" ↑/↓ 滚动  •  PgUp/PgDown 翻页  •  g/G 顶部/底部  •  s 跳过当前  •  r 重新开始  •  q 退出 ",
	)

	return lipgloss.JoinVertical(
		lipgloss.Top,
		m.viewport.View(),
		statusText,
		helpText,
	)
}

func main() {
	rand.Seed(time.Now().UnixNano())

	p := tea.NewProgram(
		newModel(),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
	}
}
