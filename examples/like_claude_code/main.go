package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/term"
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

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("42"))

	dimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))

	// Alan Code 风格输入框样式
	inputPromptStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("39")).
				Bold(true)

	inputUserStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("226")).
			Bold(true)

	assistantStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("255"))

	// 布局样式
	headerStyle = lipgloss.NewStyle().
			Border(lipgloss.ThickBorder()).
			BorderForeground(lipgloss.Color("39")).
			Padding(0, 2)

	helpBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.ThickBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(0, 2)

	statusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Italic(true)

	// 消息气泡样式
	userBubbleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("226")).
			Bold(true)
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

func getTerminalWidth() int {
	width, _, err := term.GetSize(os.Stdout.Fd())
	if err != nil {
		return 80
	}
	return width
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

// ========== UI 组件 ==========

// Alan Code 风格输入提示符
func getUserPrompt() string {
	return inputPromptStyle.Render("> ")
}

func getUserMessageHeader() string {
	return ">"
}

// 显示思考动画
func showThinkingAnimation() {
	spinner := NewSpinner()
	spinner.message = "Thinking..."
	spinner.Start("Thinking...")
	time.Sleep(2500 * time.Millisecond)
	spinner.Stop()

}

func printHeader() {
	fmt.Println()
	header := headerStyle.Render(
		lipgloss.JoinVertical(
			lipgloss.Center,
			primaryStyle.Bold(true).Render("Alan Code"),
		),
	)
	fmt.Println(header)
	fmt.Println(statusStyle.Render("  Type your question. Use /help for commands."))
	fmt.Println()
}

func printHelp() {
	commands := []string{
		"  /help          " + dimStyle.Render("Show this help message"),
		"  /clear         " + dimStyle.Render("Clear conversation"),
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

func deletePreviousLine() {
	fmt.Print("\033[1A") // 上移一行
	fmt.Print("\033[K")  // 删除当前行
}

func printDivider(width int) {
	fmt.Println(dimStyle.Render(strings.Repeat("─", width)))
}

// ========== 主程序 ==========
func main() {
	rand.Seed(time.Now().UnixNano())

	printHeader()

	ai := NewAIResponseGenerator()
	reader := bufio.NewReader(os.Stdin)
	conversation := &Conversation{
		Messages:    []Message{},
		CurrentFile: "",
		ProjectPath: ".",
	}

	for {
		width := getTerminalWidth()
		printDivider(width)
		fmt.Print(getUserPrompt())
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
		}

		// 显示用户消息

		deletePreviousLine()
		fmt.Println(dimStyle.Background(lipgloss.Color("#282828")).Width(width - 2).Render(getUserMessageHeader() + " " + input))

		// 记录用户消息
		conversation.Messages = append(conversation.Messages, Message{
			Role:      "user",
			Content:   input,
			Timestamp: time.Now(),
		})

		// 检测问题复杂度并显示相应的思考动画

		fmt.Println()
		showThinkingAnimation()

		// 生成并流式输出响应
		response := ai.GenerateResponse(input)
		ai.StreamOutput(response)

		// 记录助手消息
		conversation.Messages = append(conversation.Messages, Message{
			Role:      "assistant",
			Content:   response,
			Timestamp: time.Now(),
		})

		fmt.Println()
	}
}
