package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// ========== 全局样式定义 ==========
var (
	// 主题色
	primaryStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("39")).
			Bold(true)

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("46")).
			Bold(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Bold(true)

	warningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("208")).
			Bold(true)

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("42"))

	dimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))

	// Alan Code 风格输入框样式
	inputPromptStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("39")).
				Bold(true)

	inputBracketStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("240"))

	inputUserStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("226")).
			Bold(true)

	inputAtStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))

	assistantStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("255"))

	codeStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("214"))

	dividerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))

	// 布局样式
	headerStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("39")).
			Padding(0, 1)

	helpBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(0, 2)

	statusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Italic(true)

	// 消息气泡样式
	userBubbleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("226")).
			Bold(true)

	assistantBubbleStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("255"))
)

// ========== 上下文管理 ==========
type Conversation struct {
	Messages    []Message
	CurrentFile string
	ProjectPath string
}

type Message struct {
	Role      string
	Content   string
	Timestamp time.Time
}

// ========== 模拟 AI 响应生成器 ==========
type AIResponseGenerator struct {
	typingSpeed int
}

func NewAIResponseGenerator() *AIResponseGenerator {
	return &AIResponseGenerator{
		typingSpeed: 10,
	}
}

// 根据用户输入生成模拟响应
func (g *AIResponseGenerator) GenerateResponse(userInput string) string {
	input := strings.ToLower(userInput)

	switch {
	case strings.Contains(input, "hello") || strings.Contains(input, "hi"):
		return "Hello! I'm Alan Code, your AI coding assistant. I can help you with:\n\n" +
			"• Writing and reviewing code\n" +
			"• Debugging issues\n" +
			"• Explaining complex concepts\n" +
			"• Running tests and commands\n\n" +
			"How can I help you today?"

	case strings.Contains(input, "explain"):
		return "Let me explain that for you:\n\n```go\n// This is a Go function example\nfunc Explain() string {\n    return \"I'll break this down step by step\"\n}\n```\n\nThe key points to understand are:\n1. **Structure**: The code follows a clear pattern\n2. **Best Practices**: Error handling is included\n3. **Performance**: It's optimized for speed\n\nWould you like me to elaborate on any part?"

	case strings.Contains(input, "debug") || strings.Contains(input, "error"):
		return "🔍 **Debugging Analysis**\n\nI've analyzed the error. Here's what I found:\n\n**Root Cause:**\nVariable `result` is being used before initialization\n\n**Solution:**\n```go\nvar result string\nresult = computeValue()\n```\n\n**Prevention:**\n• Always initialize variables before use\n• Use `:=` for declaration + initialization\n• Enable `go vet` in your CI pipeline\n\nShould I apply this fix?"

	case strings.Contains(input, "perf") || strings.Contains(input, "optimize"):
		return "⚡ **Performance Optimization Recommendations**\n\nBased on my analysis, here are the top improvements:\n\n**1. Reduce allocations** (预计提升 30%)\n   • Use object pooling\n   • Pre-allocate slices with capacity\n\n**2. Optimize database queries** (预计提升 50%)\n   • Add proper indexes\n   • Use connection pooling\n\n**3. Implement caching** (预计提升 80%)\n   • Redis for frequent queries\n   • In-memory cache for config\n\nWant me to implement these optimizations?"

	case strings.Contains(input, "test"):
		return "🧪 **Test Strategy**\n\nI recommend the following test approach:\n\n**Unit Tests:**\n```go\nfunc TestFunction(t *testing.T) {\n    tests := []struct{\n        name string\n        input int\n        want int\n    }{\n        {\"positive\", 5, 5},\n        {\"negative\", -1, 0},\n    }\n    // ... test logic\n}\n```\n\n**Integration Tests:** Test API endpoints\n**E2E Tests:** Use Cypress/Playwright\n\nTest coverage currently: ~65%\nTarget coverage: 85%"

	case strings.Contains(input, "review"):
		return "📋 **Code Review**\n\nI've reviewed the changes. Here's my feedback:\n\n✅ **What's working well:**\n• Clear naming conventions\n• Good separation of concerns\n\n⚠️ **Issues to address:**\n• Missing error handling in line 42\n• Potential nil pointer in `processData()`\n• Unused import: `fmt`\n\n💡 **Suggestions:**\n• Add documentation for exported functions\n• Consider using context for timeouts\n\nWould you like me to fix these issues?"

	default:
		return "I understand you're asking about **" + userInput + "**.\n\nHere's what I can do:\n\n🔄 **Searching codebase...**\n📖 **Reading documentation...**\n💡 **Generating solution...**\n\nBased on my analysis, I recommend:\n\n```go\n// Example implementation\nfunc ProcessRequest(req *Request) (*Response, error) {\n    // TODO: Implement based on requirements\n    return &Response{}, nil\n}\n```\n\nCould you provide more context about your specific use case? This will help me give you a more accurate solution."
	}
}

// 流式输出
func (g *AIResponseGenerator) StreamOutput(text string) {
	lines := strings.Split(text, "\n")
	inCodeBlock := false

	for _, line := range lines {
		if strings.Contains(line, "```") {
			inCodeBlock = !inCodeBlock
			fmt.Println(line)
		} else if inCodeBlock {
			fmt.Println(line)
			time.Sleep(5 * time.Millisecond)
		} else if line == "" {
			fmt.Println()
		} else {
			for _, char := range line {
				fmt.Print(string(char))
				time.Sleep(time.Duration(g.typingSpeed) * time.Millisecond)
			}
			fmt.Println()
		}
		time.Sleep(20 * time.Millisecond)
	}
}

// 显示 Spinner（带进度百分比）
func showSpinnerWithProgress(message string, duration time.Duration) {
	spinnerChars := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	done := make(chan bool)
	startTime := time.Now()

	go func() {
		i := 0
		for {
			select {
			case <-done:
				fmt.Print("\r\033[K")
				fmt.Println(successStyle.Render("✓") + " " + message)
				return
			default:
				elapsed := time.Since(startTime)
				percent := int(float64(elapsed) / float64(duration) * 100)
				if percent > 100 {
					percent = 100
				}
				spinner := spinnerChars[i%len(spinnerChars)]
				progressBar := strings.Repeat("█", percent/5) + strings.Repeat("░", 20-percent/5)
				fmt.Printf("\r%s %s %s [%s] %d%%",
					infoStyle.Render(spinner),
					message,
					dimStyle.Render(progressBar),
					dimStyle.Render(fmt.Sprintf("%d%%", percent)))
				time.Sleep(80 * time.Millisecond)
				i++
			}
		}
	}()

	time.Sleep(duration)
	close(done)
	time.Sleep(100 * time.Millisecond)
}

// 简单 Spinner
func showSpinner(message string, duration time.Duration) {
	spinnerChars := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	done := make(chan bool)

	go func() {
		i := 0
		for {
			select {
			case <-done:
				fmt.Print("\r\033[K")
				fmt.Println(successStyle.Render("✓") + " " + message)
				return
			default:
				fmt.Printf("\r%s %s", infoStyle.Render(spinnerChars[i%len(spinnerChars)]), message)
				time.Sleep(80 * time.Millisecond)
				i++
			}
		}
	}()

	time.Sleep(duration)
	close(done)
	time.Sleep(100 * time.Millisecond)
}

// ========== UI 组件 ==========

// Alan Code 风格输入提示符
func getUserPrompt() string {
	return inputPromptStyle.Render("> ") + inputUserStyle.Render("User") + inputPromptStyle.Render(" ")
}

// 显示思考动画（多阶段）
func showThinkingAnimation(complexity int) {
	// 根据问题复杂度决定思考时间
	// complexity: 1-简单, 2-中等, 3-复杂
	var stages []string
	var totalDuration time.Duration

	switch complexity {
	case 1:
		stages = []string{
			"🤔 理解问题...",
			"💭 准备回答...",
		}
		totalDuration = 800 * time.Millisecond
	case 2:
		stages = []string{
			"🤔 分析问题...",
			"🔍 搜索相关知识...",
			"💭 组织回答...",
		}
		totalDuration = 1500 * time.Millisecond
	default:
		stages = []string{
			"🤔 深入分析...",
			"🔍 检索代码库...",
			"📖 阅读文档...",
			"⚡ 生成解决方案...",
			"💭 优化回答...",
		}
		totalDuration = 2500 * time.Millisecond
	}

	stageDuration := totalDuration / time.Duration(len(stages))

	for _, stage := range stages {
		fmt.Printf("\r\033[K%s %s", infoStyle.Render("●"), stage)
		time.Sleep(stageDuration)
	}
	fmt.Print("\r\033[K")
	fmt.Println(infoStyle.Render("✓") + " 准备就绪")
}

// 检测问题复杂度
func detectComplexity(input string) int {
	input = strings.ToLower(input)

	// 复杂问题关键词
	complexKeywords := []string{"explain", "how", "why", "debug", "performance", "optimize", "architecture", "design"}
	// 中等问题关键词
	mediumKeywords := []string{"review", "test", "what", "when", "where", "help"}

	for _, kw := range complexKeywords {
		if strings.Contains(input, kw) {
			return 3 // 复杂
		}
	}

	for _, kw := range mediumKeywords {
		if strings.Contains(input, kw) {
			return 2 // 中等
		}
	}

	return 1 // 简单
}

func printHeader() {
	fmt.Println()
	header := headerStyle.Render(
		lipgloss.JoinVertical(
			lipgloss.Center,
			primaryStyle.Bold(true).Render("Alan Code"),
			dimStyle.Render("AI Coding Assistant"),
		),
	)
	fmt.Println(header)
	fmt.Println()
	fmt.Println(statusStyle.Render("  Type your question. Use /help for commands."))
	fmt.Println()
}

func printHelp() {
	commands := []string{
		"  /help          " + dimStyle.Render("Show this help message"),
		"  /clear         " + dimStyle.Render("Clear conversation"),
		"  /review        " + dimStyle.Render("Review code changes"),
		"  /test          " + dimStyle.Render("Run tests"),
		"  /explain       " + dimStyle.Render("Explain code"),
		"  /exit          " + dimStyle.Render("Exit Alan Code"),
		"",
		"  " + dimStyle.Render("Tip: Just type your question naturally!"),
	}

	helpContent := lipgloss.JoinVertical(
		lipgloss.Left,
		primaryStyle.Render("📖 Available Commands"),
		"",
		strings.Join(commands, "\n"),
	)

	helpBox := helpBoxStyle.Render(helpContent)
	fmt.Println(helpBox)
	fmt.Println()
}

func printDivider() {
	fmt.Println(dimStyle.Render(strings.Repeat("─", 60)))
}

func getUserMessageHeader() string {
	return userBubbleStyle.Render("You:")
}

func getAssistantHeader() string {
	return assistantStyle.Render("💬 ") + successStyle.Render("Alan Code") + assistantStyle.Render(":")
}

// ========== 主程序 ==========
func main() {
	rand.Seed(time.Now().UnixNano())

	// 清屏
	fmt.Print("\033[2J\033[H")

	// 打印头部
	printHeader()

	ai := NewAIResponseGenerator()
	reader := bufio.NewReader(os.Stdin)
	conversation := &Conversation{
		Messages:    []Message{},
		CurrentFile: "",
		ProjectPath: ".",
	}

	for {
		// 显示 Alan Code 风格提示符
		fmt.Print(getUserPrompt())

		// 读取用户输入
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(errorStyle.Render("Error reading input"))
			continue
		}
		input = strings.TrimSpace(input)

		if input == "" {
			continue
		}

		// 处理命令
		switch input {
		case "/exit", "/quit", "exit", "quit":
			fmt.Println()
			fmt.Println(successStyle.Render("  👋 Goodbye! Happy coding!"))
			fmt.Println()
			return

		case "/help", "help", "?":
			printHelp()
			continue

		case "/clear", "clear":
			fmt.Print("\033[2J\033[H")
			printHeader()
			conversation.Messages = []Message{}
			fmt.Println(successStyle.Render("  ✨ Conversation cleared"))
			fmt.Println()
			continue

		case "/review":
			fmt.Println()
			showSpinnerWithProgress("Reviewing code changes", 1500*time.Millisecond)
			fmt.Println()
			response := `**Code Review Summary**

Files changed: 3
Lines added: +47
Lines removed: -12

**Issues found:**
• Line 23: Missing error handling
• Line 56: Potential nil pointer

**Suggestions:**
• Add unit tests for new functions
• Consider extracting duplicate logic

Would you like me to fix these issues?`
			fmt.Println(getAssistantHeader())
			fmt.Println()
			ai.StreamOutput(response)
			printDivider()
			continue

		case "/test":
			fmt.Println()
			showSpinnerWithProgress("Running tests", 2000*time.Millisecond)
			fmt.Println()
			response := `🧪 **Test Results**

✅ TestUserLogin - PASS (0.23s)
✅ TestDataFetch - PASS (0.45s)
⚠️ TestAuthMiddleware - SKIPPED
❌ TestDatabaseConnection - FAIL (0.12s)

**Coverage:** 73.5%

**Failed test details:**
TestDatabaseConnection timed out

**Recommendation:** Check if PostgreSQL is running`
			fmt.Println(getAssistantHeader())
			fmt.Println()
			ai.StreamOutput(response)
			printDivider()
			continue

		case "/explain":
			fmt.Println()
			showSpinnerWithProgress("Analyzing code", 1800*time.Millisecond)
			fmt.Println()
			response := `📚 **Explanation**

To get the best explanation:
1. Select the code you want explained
2. Ask specific questions about the code
3. Use "explain how X works" format

**Example:** "explain how the authentication flow works"`
			fmt.Println(getAssistantHeader())
			fmt.Println()
			ai.StreamOutput(response)
			printDivider()
			continue
		}

		// 显示用户消息
		fmt.Println()
		fmt.Println(getUserMessageHeader())
		fmt.Println(dimStyle.Render("  " + input))

		// 记录用户消息
		conversation.Messages = append(conversation.Messages, Message{
			Role:      "user",
			Content:   input,
			Timestamp: time.Now(),
		})

		// 检测问题复杂度并显示相应的思考动画
		complexity := detectComplexity(input)

		// 显示思考动画
		fmt.Println()
		showThinkingAnimation(complexity)

		// 显示助手头部
		fmt.Println()
		fmt.Println(getAssistantHeader())
		fmt.Println()

		// 生成并流式输出响应
		response := ai.GenerateResponse(input)
		ai.StreamOutput(response)

		// 记录助手消息
		conversation.Messages = append(conversation.Messages, Message{
			Role:      "assistant",
			Content:   response,
			Timestamp: time.Now(),
		})

		// 输出分隔线
		fmt.Println()
		printDivider()
	}
}
