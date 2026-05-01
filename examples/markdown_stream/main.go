package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
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

	// 流式输出原始文本
	for _, ch := range markdown {
		fmt.Print(string(ch))
		time.Sleep(10 * time.Millisecond) // 模拟流式效果
	}

	// 确保末尾换行
	if !strings.HasSuffix(markdown, "\n") {
		fmt.Println()
	}
}
