package editor

// deleteCharacter deletes the character under the cursor
func (e *Editor) deleteCharacter() {
	p := e.active()
	if p.cursorY >= p.buffer.LineCount() {
		return
	}

	line := p.buffer.Line(p.cursorY)
	if p.cursorX >= lineRuneCount(line) {
		// At end of line, join with next line
		if p.cursorY < p.buffer.LineCount()-1 {
			nextLine := p.buffer.Line(p.cursorY + 1)
			p.buffer.DeleteLine(p.cursorY + 1)
			// Append next line content to current line
			for _, r := range []rune(nextLine) {
				p.buffer.InsertRune(p.cursorY, lineRuneCount(p.buffer.Line(p.cursorY)), r)
			}
		}
		return
	}

	// Delete character at cursor
	p.buffer.DeleteRune(p.cursorY, p.cursorX+1)
	e.clampCursor()
	e.adjustScroll()
}

// deleteCurrentLine deletes the entire current line
func (e *Editor) deleteCurrentLine() {
	p := e.active()
	if p.buffer.LineCount() == 1 {
		// Last line, just clear it
		line := p.buffer.Line(0)
		for i := lineRuneCount(line); i > 0; i-- {
			p.buffer.DeleteRune(0, i)
		}
		p.cursorX = 0
		return
	}

	p.buffer.DeleteLine(p.cursorY)

	// Adjust cursor position
	if p.cursorY >= p.buffer.LineCount() {
		p.cursorY = p.buffer.LineCount() - 1
	}
	p.cursorX = 0

	e.clampCursor()
	e.adjustScroll()
}
