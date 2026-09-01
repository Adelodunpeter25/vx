package syntax

import (
	"strings"
	"sync"

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
	mu             sync.RWMutex
	highlighting   bool
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
	if buf.LineCount() > MaxHighlightLines || buf.IsLazy() {
		e.mu.Lock()
		e.tooLarge = true
		e.cache = nil // Free memory
		e.mu.Unlock()
		return nil
	}

	currentVersion := buf.ModVersion()

	// Fast path: cache ready for this version
	e.mu.RLock()
	if e.cachedVersion == currentVersion && e.cache != nil {
		if styled, ok := e.cache[lineNum]; ok {
			e.mu.RUnlock()
			return styled
		}
		e.mu.RUnlock()
		// Cache ready but line not in cache (edge case) — fallback to per-line
		return e.highlighter.HighlightLine(line)
	}
	highlighting := e.highlighting
	e.mu.RUnlock()

	// Cache not ready — trigger async full-buffer highlight if not already running
	if !highlighting {
		e.mu.Lock()
		// Double-check under write lock
		if !e.highlighting && e.cachedVersion != currentVersion {
			e.highlighting = true
			// Snapshot full text synchronously (fast, no lex) to avoid data race
			var fullText strings.Builder
			// Pre-grow roughly: avg 80 chars per line
			fullText.Grow(buf.LineCount() * 80)
			for i := 0; i < buf.LineCount(); i++ {
				if i > 0 {
					fullText.WriteString("\n")
				}
				fullText.WriteString(buf.Line(i))
			}
			text := fullText.String()
			version := currentVersion
			go func(v int, t string) {
				lines := e.highlighter.HighlightText(t)
				m := make(map[int][]highlight.StyledRune, len(lines))
				for i, l := range lines {
					m[i] = l
				}
				e.mu.Lock()
				// Only commit if buffer hasn't changed again
				if v == buf.ModVersion() {
					e.cache = m
					e.cachedVersion = v
				}
				e.highlighting = false
				e.mu.Unlock()
			}(version, text)
		}
		e.mu.Unlock()
	}

	// While async highlight is pending, return per-line highlight (cheap, non-blocking)
	// This gives immediate coloring without waiting for full-buffer lex.
	// For true plain-first-frame, return nil instead.
	return e.highlighter.HighlightLine(line)
}

func (e *Engine) highlightBuffer(buf *buffer.Buffer) {
	// Kept for synchronous callers (e.g., tests); builds cache inline.
	m := make(map[int][]highlight.StyledRune)
	var fullText strings.Builder
	fullText.Grow(buf.LineCount() * 80)
	for i := 0; i < buf.LineCount(); i++ {
		if i > 0 {
			fullText.WriteString("\n")
		}
		fullText.WriteString(buf.Line(i))
	}
	lines := e.highlighter.HighlightText(fullText.String())
	for i, line := range lines {
		m[i] = line
	}
	e.mu.Lock()
	e.cache = m
	e.cachedVersion = buf.ModVersion()
	e.mu.Unlock()
}

func (e *Engine) InvalidateCache() {
	e.mu.Lock()
	e.cachedVersion = -1
	e.cache = make(map[int][]highlight.StyledRune)
	e.mu.Unlock()
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
