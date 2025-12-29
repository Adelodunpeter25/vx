package editor

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
	
	// Vertical scroll
	if e.cursorY < e.offsetY {
		e.offsetY = e.cursorY
	}
	if e.cursorY >= e.offsetY+contentHeight {
		e.offsetY = e.cursorY - contentHeight + 1
	}
	
	// Horizontal scroll - keep cursor visible with margin
	if e.cursorX < e.offsetX {
		e.offsetX = e.cursorX
	}
	if e.cursorX >= e.offsetX+e.width-1 {
		e.offsetX = e.cursorX - e.width + 2
	}
	
	// Ensure offsetX doesn't go negative
	if e.offsetX < 0 {
		e.offsetX = 0
	}
}
