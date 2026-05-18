package component

import (
	"sync"

	"charm.land/glamour/v2"
	"charm.land/glamour/v2/styles"
)

var (
	mdCacheMu sync.Mutex
	mdCache   = map[int]*glamour.TermRenderer{}
)

func getRenderer(width int, bgColor string) *glamour.TermRenderer {
	mdCacheMu.Lock()
	defer mdCacheMu.Unlock()

	if r, ok := mdCache[width]; ok {
		return r
	}

	cfg := styles.DarkStyleConfig
	bg := bgColor
	cfg.Document.StylePrimitive.BackgroundColor = &bg
	cfg.H1.StylePrimitive.BackgroundColor = &bg
	cfg.Code.StylePrimitive.BackgroundColor = &bg
	cfg.CodeBlock.StylePrimitive.BackgroundColor = &bg
	if cfg.CodeBlock.Chroma != nil {
		cfg.CodeBlock.Chroma.Background.BackgroundColor = &bg
		cfg.CodeBlock.Chroma.Error.BackgroundColor = &bg
	}

	r, err := glamour.NewTermRenderer(
		glamour.WithStyles(cfg),
		glamour.WithWordWrap(width),
	)
	if err != nil {
		return nil
	}

	mdCache[width] = r
	return r
}

func RenderMarkdown(content string, width int, bgColor string) string {
	r := getRenderer(width, bgColor)
	if r == nil {
		return content
	}

	out, err := r.Render(content)
	if err != nil {
		return content
	}

	return out
}
