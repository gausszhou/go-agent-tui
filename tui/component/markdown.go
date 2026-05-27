package component

import (
	"github.com/gausszhou/gruff"
)

func RenderMarkdown(content string, width int, _ string) string {
	out, err := gruff.Render(content, gruff.WithDark(), gruff.WithWordWrap(width))
	if err != nil {
		return content
	}
	return out
}
