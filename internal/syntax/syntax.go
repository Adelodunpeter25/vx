package syntax

import (
	"strings"

	"github.com/Adelodunpeter25/vx/internal/buffer"
	"github.com/Adelodunpeter25/vx/pkg/highlight"
)

const MaxHighlightLines = 10000 // Don't highlight files larger than this

type Engine struct {
	highlighter    *highlight.Highlighter
	enabled        bool
	cache          map[int][]highlight.StyledRune
	cachedVersion  int
	tooLarge       bool
}

func New(filename string) *Engine {
	return &Engine{
		highlighter:   highlight.New(filename),
		enabled:       true,
		cache:         make(map[int][]highlight.StyledRune),
		cachedVersion: -1,
	}
}

func (e *Engine) HighlightLine(lineNum int, line string, buf *buffer.Buffer) []highlight.StyledRune {
	if !e.enabled || e.tooLarge {
		return nil
	}
	
	// Check if buffer is too large for highlighting
	if buf.LineCount() > MaxHighlightLines {
		e.tooLarge = true
		e.cache = nil // Free memory
		return nil
	}
	
	// Only re-highlight if buffer version changed
	currentVersion := buf.ModVersion()
	if e.cachedVersion != currentVersion {
		e.highlightBuffer(buf)
		e.cachedVersion = currentVersion
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
	e.cachedVersion = -1
	e.cache = make(map[int][]highlight.StyledRune)
}

func (e *Engine) Toggle() {
	e.enabled = !e.enabled
}

func (e *Engine) IsEnabled() bool {
	return e.enabled
}

func (e *Engine) IsTooLarge() bool {
	return e.tooLarge
}
