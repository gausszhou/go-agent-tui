package main

import (
	"fmt"
	"time"
)

type Spinner struct {
	Frames   []string
	delay    time.Duration
	stopChan chan bool
	message  string
}

func NewSpinner() *Spinner {
	return &Spinner{
		Frames: []string{"|", "/", "-", "\\"}, // √
		// Frames: []string{"⣾ ", "⣽ ", "⣻ ", "⢿ ", "⡿ ", "⣟ ", "⣯ ", "⣷ "}, // √
		// Frames: []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}, // √
		// Frames: []string{"⢄", "⢂", "⢁", "⡁", "⡈", "⡐", "⡠"}, // √
		// Frames: []string{"█", "▓", "▒", "░"}, // x
		// Frames: []string{"∙∙∙", "●∙∙", "∙●∙", "∙∙●"}, // x
		// Frames: []string{"🌍", "🌎", "🌏"},  // √
		// Frames: []string{"🌑", "🌒", "🌓", "🌔", "🌕", "🌖", "🌗", "🌘"}, // √
		// Frames: []string{"☱", "☲", "☴", "☲"}, // √
		delay: 80 * time.Millisecond,
	}
}

// Start 启动 spinner（在 goroutine 中运行）
func (s *Spinner) Start(message string) {
	s.stopChan = make(chan bool)

	go func() {
		i := 0
		for {
			select {
			case <-s.stopChan:
				fmt.Print("\r\033[K")
				return
			default:
				frame := s.Frames[i%len(s.Frames)]
				fmt.Printf("\r%s %s", frame, message)
				i++
				time.Sleep(s.delay)
			}
		}
	}()
}

// Stop 停止 spinner
func (s *Spinner) Stop() {
	if s.stopChan != nil {
		s.stopChan <- true
		close(s.stopChan)
	}
}
