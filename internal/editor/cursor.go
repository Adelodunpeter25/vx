package editor

import "github.com/Adelodunpeter25/vx/internal/wrap"

func (e *Editor) clampCursor() {
	line := e.buffer.Line(e.cursorY)
	maxX := len(line)
	if e.mode == ModeNormal && maxX > 0 {
		maxX--
	}
	if e.cursorX > maxX {
		e.cursorX = maxX
	}
	if e.cursorX < 0 {
		e.cursorX = 0
	}
}

func (e *Editor) adjustScroll() {
	contentHeight := e.height - 1
	gutterWidth := e.getGutterWidth()
	maxWidth := e.width - gutterWidth
	
	// Calculate visual line position of cursor
	visualLine := 0
	for lineNum := 0; lineNum < e.cursorY && lineNum < e.buffer.LineCount(); lineNum++ {
		line := e.buffer.Line(lineNum)
		visualLine += wrap.VisualLineCount(line, maxWidth)
	}
	// Add wrapped rows within current line
	if maxWidth > 0 {
		visualLine += e.cursorX / maxWidth
	}
	
	// Calculate visual line of offsetY
	offsetVisual := 0
	for lineNum := 0; lineNum < e.offsetY && lineNum < e.buffer.LineCount(); lineNum++ {
		line := e.buffer.Line(lineNum)
		offsetVisual += wrap.VisualLineCount(line, maxWidth)
	}
	
	// Vertical scroll - adjust offsetY to keep cursor visible
	if visualLine < offsetVisual {
		// Cursor above viewport - scroll up
		e.offsetY = e.findLineAtVisualRow(visualLine, maxWidth)
	}
	if visualLine >= offsetVisual+contentHeight {
		// Cursor below viewport - scroll down
		targetVisual := visualLine - contentHeight + 1
		e.offsetY = e.findLineAtVisualRow(targetVisual, maxWidth)
	}
	
	// Ensure offsetY doesn't go negative
	if e.offsetY < 0 {
		e.offsetY = 0
	}
	
	// No horizontal scroll needed with wrapping
	e.offsetX = 0
}

// findLineAtVisualRow finds which buffer line contains the given visual row
func (e *Editor) findLineAtVisualRow(targetVisual, maxWidth int) int {
	visualLine := 0
	for lineNum := 0; lineNum < e.buffer.LineCount(); lineNum++ {
		line := e.buffer.Line(lineNum)
		lineVisualCount := wrap.VisualLineCount(line, maxWidth)
		if visualLine+lineVisualCount > targetVisual {
			return lineNum
		}
		visualLine += lineVisualCount
	}
	return e.buffer.LineCount() - 1
}
