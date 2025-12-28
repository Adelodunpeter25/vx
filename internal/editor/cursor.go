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
	if e.cursorY < e.offsetY {
		e.offsetY = e.cursorY
	}
	if e.cursorY >= e.offsetY+contentHeight {
		e.offsetY = e.cursorY - contentHeight + 1
	}
}
