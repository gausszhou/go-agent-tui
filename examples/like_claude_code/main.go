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

	toolStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("208")).
			Bold(true)

	toolResultStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("228"))

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

// ========== Tool Call 相关结构 ==========

// ToolDefinition 定义了一个可用的工具
type ToolDefinition struct {
	Name        string
	Description string
	Parameters  map[string]ToolParameter
}

// ToolParameter 定义工具参数
type ToolParameter struct {
	Type        string
	Description string
	Required    bool
}

// ToolCall 表示一次工具调用
type ToolCall struct {
	ID        string
	Name      string
	Arguments map[string]interface{}
}

// ToolResult 表示工具执行结果
type ToolResult struct {
	ToolCallID string
	Content    string
	Error      error
}

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
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`
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

// 检测用户意图并生成相应的工具调用
func (g *AIResponseGenerator) DetectToolCalls(userInput string) []ToolCall {
	input := strings.ToLower(userInput)
	var toolCalls []ToolCall
	toolCalls = append(toolCalls, ToolCall{
		ID:   generateToolCallID(),
		Name: "read_file",
		Arguments: map[string]interface{}{
			"path": extractFilePath(userInput),
		},
	})
	toolCalls = append(toolCalls, ToolCall{
		ID:   generateToolCallID(),
		Name: "write_file",
		Arguments: map[string]interface{}{
			"path":    "output.go",
			"content": generateSampleCode(input),
		},
	})

	toolCalls = append(toolCalls, ToolCall{
		ID:   generateToolCallID(),
		Name: "execute_command",
		Arguments: map[string]interface{}{
			"command": extractCommand(userInput),
		},
	})

	toolCalls = append(toolCalls, ToolCall{
		ID:   generateToolCallID(),
		Name: "search_code",
		Arguments: map[string]interface{}{
			"query":    extractSearchQuery(userInput),
			"fileType": extractFileType(userInput),
		},
	})

	toolCalls = append(toolCalls, ToolCall{
		ID:   generateToolCallID(),
		Name: "list_files",
		Arguments: map[string]interface{}{
			"path": ".",
		},
	})

	toolCalls = append(toolCalls, ToolCall{
		ID:   generateToolCallID(),
		Name: "run_tests",
		Arguments: map[string]interface{}{
			"testPath": extractTestPath(userInput),
		},
	})

	return toolCalls
}

// 生成工具调用ID
func generateToolCallID() string {
	return fmt.Sprintf("call_%d_%s", time.Now().UnixNano(), randomString(6))
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// 辅助函数: 提取参数
func extractFilePath(input string) string {
	// 简单的路径提取逻辑
	if strings.Contains(input, "main.go") {
		return "main.go"
	}
	if strings.Contains(input, "config") {
		return "config.json"
	}
	return "example.go"
}

func extractCommand(input string) string {
	if strings.Contains(input, "go build") {
		return "go build ./..."
	}
	if strings.Contains(input, "go test") {
		return "go test -v ./..."
	}
	return "go run main.go"
}

func extractSearchQuery(input string) string {
	// 提取搜索关键词
	words := strings.Fields(input)
	for i, word := range words {
		if strings.Contains(word, "search") || strings.Contains(word, "查找") {
			if i+1 < len(words) {
				return words[i+1]
			}
		}
	}
	return "function"
}

func extractFileType(input string) string {
	if strings.Contains(input, ".go") {
		return ".go"
	}
	if strings.Contains(input, ".js") {
		return ".js"
	}
	return ""
}

func extractTestPath(input string) string {
	if strings.Contains(input, "unit") {
		return "./internal/..."
	}
	return "./..."
}

func generateSampleCode(topic string) string {
	return fmt.Sprintf(`// 根据需求 "%s" 生成的代码
package main

import (
    "fmt"
    "log"
)

func main() {
    fmt.Println("Hello from Alan Code!")
    
    // 在这里实现你的逻辑
    result, err := ProcessRequest()
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Result: %%v\n", result)
}

func ProcessRequest() (interface{}, error) {
    // TODO: 实现业务逻辑
    return nil, nil
}
`, topic)
}

// 根据用户输入生成模拟响应 (支持工具调用)
func (g *AIResponseGenerator) GenerateResponse(userInput string, toolResults []ToolResult) string {
	input := strings.ToLower(userInput)
	switch {
	case strings.Contains(input, "hello") || strings.Contains(input, "hi"):
		return "Hello! I'm Alan Code, your AI coding assistant. I can help you with:\n\n" +
			"• Writing and reviewing code\n" +
			"• Debugging issues\n" +
			"• Explaining complex concepts\n" +
			"• Running tests and commands\n" +
			"• **Reading and writing files** (try saying 'read file' or 'write file')\n" +
			"• **Searching code** (try 'search for X')\n\n" +
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
		return "I understand you're asking about **" + userInput + "**.\n\nHere's what I can do:\n\n🔄 **Searching codebase...**\n📖 **Reading documentation...**\n💡 **Generating solution...**\n\nBased on my analysis, I recommend:\n\n```go\n// Example implementation\nfunc ProcessRequest(req *Request) (*Response, error) {\n    // TODO: Implement based on requirements\n    return &Response{}, nil\n}\n```\n\nCould you provide more context about your specific use case? This will help me give you a more accurate solution.\n\n💡 **Tip**: Try saying \"search for main function\" or \"read main.go\" to test my tool capabilities!"
	}
}

// 基于工具调用结果生成响应
func (g *AIResponseGenerator) generateResponseWithTools(userInput string, toolResults []ToolResult) string {
	var response strings.Builder

	response.WriteString("🛠️ **我已经执行了以下操作:**\n\n")

	for i, result := range toolResults {
		response.WriteString(fmt.Sprintf("%d. ", i+1))
		if result.Error != nil {
			response.WriteString(fmt.Sprintf("❌ **工具执行失败**: %s\n", result.Error.Error()))
		} else {
			response.WriteString("✅ **工具执行成功**\n")
			response.WriteString(fmt.Sprintf("\n```\n%s\n```\n\n", result.Content))
		}
	}

	// 根据用户输入和工具结果生成总结
	response.WriteString("📋 **总结:**\n\n")

	userInputLower := strings.ToLower(userInput)
	if strings.Contains(userInputLower, "read") {
		response.WriteString("我已经读取了请求的文件。你可以：\n")
		response.WriteString("• 让我分析文件中的代码\n")
		response.WriteString("• 询问具体的函数实现\n")
		response.WriteString("• 请求修改文件内容\n")
	} else if strings.Contains(userInputLower, "write") || strings.Contains(userInputLower, "创建") {
		response.WriteString("文件创建成功！你可以：\n")
		response.WriteString("• 查看写入的内容\n")
		response.WriteString("• 请求修改文件\n")
		response.WriteString("• 运行文件测试\n")
	} else if strings.Contains(userInputLower, "search") || strings.Contains(userInputLower, "查找") {
		response.WriteString("搜索完成！我可以帮你：\n")
		response.WriteString("• 解释搜索结果中的代码\n")
		response.WriteString("• 定位到具体的文件位置\n")
		response.WriteString("• 基于搜索结果进行修改\n")
	} else if strings.Contains(userInputLower, "test") {
		response.WriteString("测试执行完成！我可以帮你：\n")
		response.WriteString("• 修复失败的测试\n")
		response.WriteString("• 增加测试覆盖率\n")
		response.WriteString("• 优化测试代码\n")
	} else {
		response.WriteString("工具已执行完毕。还有什么我可以帮你的吗？\n")
	}

	return response.String()
}

// 流式输出
func (g *AIResponseGenerator) StreamOutput(text string) {
	if text == "TOOL_CALLS_NEEDED" {
		return // 工具调用由主循环处理
	}

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
		"  /tools         " + dimStyle.Render("Show available tools"),
		"  " + dimStyle.Render("Tip: Just type your question naturally!"),
	}

	helpContent := lipgloss.JoinVertical(
		lipgloss.Left,
		primaryStyle.Render("📖 Available Commands & Tools"),
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

func createInputDivider(width int) string {
	return dimStyle.Render(strings.Repeat("─", width))
}

func getStyledInput() string {
	width := getTerminalWidth()
	divider := dimStyle.Render(strings.Repeat("─", width))

	// 打印上分隔线和提示符
	fmt.Println(divider)
	fmt.Print(inputPromptStyle.Render("> "))

	// 读取输入
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	deletePreviousLine()

	deletePreviousLine()

	return input
}

// ========== 主程序 ==========
func main() {
	rand.Seed(time.Now().UnixNano())

	printHeader()

	ai := NewAIResponseGenerator()
	// reader := bufio.NewReader(os.Stdin)
	conversation := &Conversation{
		Messages:    []Message{},
		CurrentFile: "",
		ProjectPath: ".",
	}

	for {
		width := getTerminalWidth()
		input := getStyledInput()

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

		fmt.Println(dimStyle.Background(lipgloss.Color("#282828")).Width(width - 2).Render(getUserMessageHeader() + " " + input))

		// 记录用户消息
		conversation.Messages = append(conversation.Messages, Message{
			Role:      "user",
			Content:   input,
			Timestamp: time.Now(),
		})

		// 显示普通动画
		fmt.Println()
		showThinkingAnimation()

		// 检测是否需要工具调用
		toolCalls := ai.DetectToolCalls(input)

		// 如果有工具调用，执行它们
		if len(toolCalls) > 0 {
			fmt.Println()
			fmt.Println(toolStyle.Render("🔧 Detected tool calls..."))
			time.Sleep(300 * time.Millisecond)

		}

		// 生成响应
		var response string
		response = ai.GenerateResponse(input, nil)

		// 流式输出响应
		ai.StreamOutput(response)

		// 记录助手消息
		conversation.Messages = append(conversation.Messages, Message{
			Role:      "assistant",
			Content:   response,
			Timestamp: time.Now(),
			ToolCalls: toolCalls,
		})

		fmt.Println()
	}
}
