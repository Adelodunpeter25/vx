package syntax

import "github.com/Adelodunpeter25/vx/pkg/highlight"

type Engine struct {
	highlighter *highlight.Highlighter
	enabled     bool
}

func New(filename string) *Engine {
	return &Engine{
		highlighter: highlight.New(filename),
		enabled:     true,
	}
}

func (e *Engine) HighlightLine(line string) []highlight.StyledRune {
	if !e.enabled {
		return nil
	}
	return e.highlighter.HighlightLine(line)
}

func (e *Engine) Toggle() {
	e.enabled = !e.enabled
}

func (e *Engine) IsEnabled() bool {
	return e.enabled
}
