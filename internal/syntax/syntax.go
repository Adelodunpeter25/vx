package syntax

import (
	"strings"

	"github.com/Adelodunpeter25/vx/internal/buffer"
	"github.com/Adelodunpeter25/vx/pkg/highlight"
)

type Engine struct {
	highlighter *highlight.Highlighter
	enabled     bool
	cache       map[int][]highlight.StyledRune
	lastVersion int
}

func New(filename string) *Engine {
	return &Engine{
		highlighter: highlight.New(filename),
		enabled:     true,
		cache:       make(map[int][]highlight.StyledRune),
	}
}

func (e *Engine) HighlightLine(lineNum int, line string, buf *buffer.Buffer) []highlight.StyledRune {
	if !e.enabled {
		return nil
	}
	
	// Check if we need to re-highlight entire buffer
	if len(e.cache) == 0 || buf.IsModified() {
		e.highlightBuffer(buf)
	}
	
	if styled, ok := e.cache[lineNum]; ok {
		return styled
	}
	
	return e.highlighter.HighlightLine(line)
}

func (e *Engine) highlightBuffer(buf *buffer.Buffer) {
	e.cache = make(map[int][]highlight.StyledRune)
	
	// Build full text
	var fullText strings.Builder
	for i := 0; i < buf.LineCount(); i++ {
		if i > 0 {
			fullText.WriteString("\n")
		}
		fullText.WriteString(buf.Line(i))
	}
	
	// Highlight entire buffer
	lines := e.highlighter.HighlightText(fullText.String())
	for i, line := range lines {
		e.cache[i] = line
	}
}

func (e *Engine) InvalidateCache() {
	e.cache = make(map[int][]highlight.StyledRune)
}

func (e *Engine) Toggle() {
	e.enabled = !e.enabled
}

func (e *Engine) IsEnabled() bool {
	return e.enabled
}
