package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"charm.land/glamour/v2"
)

func getCurrentFilePath() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		fmt.Println("无法获取当前文件路径")
		return "无法获取当前文件路径"
	}
	dir := filepath.Dir(filename)
	return dir
}

func main() {
	// 读取 markdown 文件
	buffer, err := os.ReadFile(filepath.Join(getCurrentFilePath(), "data.md"))
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

	// 直接输出到 stdout
	fmt.Print(rendered)

	// 确保输出末尾有换行
	if !strings.HasSuffix(rendered, "\n") {
		fmt.Println()
	}
}
