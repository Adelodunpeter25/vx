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
	cursorVisualLine := 0
	for lineNum := 0; lineNum < e.cursorY && lineNum < e.buffer.LineCount(); lineNum++ {
		line := e.buffer.Line(lineNum)
		cursorVisualLine += wrap.VisualLineCount(line, maxWidth)
	}
	
	// Find which wrapped segment contains the cursor
	currentLine := e.buffer.Line(e.cursorY)
	segments := wrap.WrapLine(currentLine, e.cursorY, maxWidth)
	for i, seg := range segments {
		segEndCol := seg.StartCol + len([]rune(seg.Text))
		if e.cursorX >= seg.StartCol && e.cursorX <= segEndCol {
			cursorVisualLine += i
			break
		}
	}
	
	// Adjust visual offset to keep cursor visible
	if cursorVisualLine < e.visualOffsetY {
		// Cursor above viewport - scroll up
		e.visualOffsetY = cursorVisualLine
	}
	if cursorVisualLine >= e.visualOffsetY + contentHeight {
		// Cursor below viewport - scroll down
		e.visualOffsetY = cursorVisualLine - contentHeight + 1
	}
	
	// Ensure visual offset doesn't go negative
	if e.visualOffsetY < 0 {
		e.visualOffsetY = 0
	}
	
	// Convert visual offset to buffer line offset for rendering
	e.offsetY = e.findLineAtVisualRow(e.visualOffsetY, maxWidth)
	
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
