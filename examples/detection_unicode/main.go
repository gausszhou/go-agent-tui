package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func main() {
	// 测试用例
	testStrings := []string{
		"Hello",    // 纯英文
		"你好",       // 纯中文
		"Hello 👋!", // 混合英文和Emoji
		"Hello 世界", // 混合英文和中文
		lipgloss.NewStyle().Foreground(lipgloss.Color(1)).Render("Hello"), // 带ANSI颜色的文本
	}

	for _, s := range testStrings {
		width := lipgloss.Width(s)
		fmt.Printf("字符串: %-15q | 显示宽度: %d 列 | 原始长度: %d 字符\n", s, width, len([]rune(s)))
	}
}
