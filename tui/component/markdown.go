package component

import (
	"sync"

	"charm.land/glamour/v2"
)

var (
	mdCacheMu sync.Mutex
	mdCache   = map[int]*glamour.TermRenderer{}
)

func getRenderer(width int) *glamour.TermRenderer {
	mdCacheMu.Lock()
	defer mdCacheMu.Unlock()

	if r, ok := mdCache[width]; ok {
		return r
	}

	r, err := glamour.NewTermRenderer(
		glamour.WithStandardStyle("dark"),
		glamour.WithWordWrap(width),
	)
	if err != nil {
		return nil
	}

	mdCache[width] = r
	return r
}

func RenderMarkdown(content string, width int) string {
	r := getRenderer(width)
	if r == nil {
		return content
	}

	out, err := r.Render(content)
	if err != nil {
		return content
	}

	return out
}
