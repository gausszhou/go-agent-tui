package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

func main() {
	// 获取终端色域配置（自动检测）
	profile := lipgloss.ColorProfile()

	switch profile {
	case termenv.TrueColor:
		fmt.Println("✅ 终端支持真彩色 (16.7 million colors)")
	case termenv.ANSI256:
		fmt.Println("🟡 终端支持 256 色")
	case termenv.ANSI:
		fmt.Println("🟠 终端仅支持 16 色")
	case termenv.Ascii:
		fmt.Println("🔴 终端不支持颜色")
	}

}
