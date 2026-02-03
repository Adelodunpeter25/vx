package editor

import (
	"unicode/utf8"

	"github.com/Adelodunpeter25/vx/internal/wrap"
)

func lineRuneCount(line string) int {
	return utf8.RuneCountInString(line)
}

func (e *Editor) clampCursor() {
	p := e.active()
	line := p.buffer.Line(p.cursorY)
	maxX := lineRuneCount(line)
	if p.mode == ModeNormal && maxX > 0 {
		maxX--
	}
	if p.cursorX > maxX {
		p.cursorX = maxX
	}
	if p.cursorX < 0 {
		p.cursorX = 0
	}
}

func (e *Editor) adjustScroll() {
	p := e.active()
	contentHeight := e.height - 1
	gutterWidth := e.getGutterWidth()
	maxWidth := e.width - gutterWidth

	// Calculate visual line position of cursor
	cursorVisualLine := 0
	for lineNum := 0; lineNum < p.cursorY && lineNum < p.buffer.LineCount(); lineNum++ {
		line := p.buffer.Line(lineNum)
		cursorVisualLine += wrap.VisualLineCount(line, maxWidth)
	}

	// Find which wrapped segment contains the cursor
	currentLine := p.buffer.Line(p.cursorY)
	segments := wrap.WrapLine(currentLine, p.cursorY, maxWidth)
	for i, seg := range segments {
		segEndCol := seg.StartCol + len([]rune(seg.Text))
		if p.cursorX >= seg.StartCol && p.cursorX <= segEndCol {
			cursorVisualLine += i
			break
		}
	}

	// Adjust visual offset to keep cursor visible
	if cursorVisualLine < p.visualOffsetY {
		// Cursor above viewport - scroll up
		p.visualOffsetY = cursorVisualLine
	}
	if cursorVisualLine >= p.visualOffsetY+contentHeight {
		// Cursor below viewport - scroll down
		p.visualOffsetY = cursorVisualLine - contentHeight + 1
	}

	// Ensure visual offset doesn't go negative
	if p.visualOffsetY < 0 {
		p.visualOffsetY = 0
	}

	// Convert visual offset to buffer line offset for rendering
	p.offsetY = e.findLineAtVisualRow(p.visualOffsetY, maxWidth)

	// No horizontal scroll needed with wrapping
	p.offsetX = 0
}

// findLineAtVisualRow finds which buffer line contains the given visual row
func (e *Editor) findLineAtVisualRow(targetVisual, maxWidth int) int {
	p := e.active()
	visualLine := 0
	for lineNum := 0; lineNum < p.buffer.LineCount(); lineNum++ {
		line := p.buffer.Line(lineNum)
		lineVisualCount := wrap.VisualLineCount(line, maxWidth)
		if visualLine+lineVisualCount > targetVisual {
			return lineNum
		}
		visualLine += lineVisualCount
	}
	return p.buffer.LineCount() - 1
}
