package editor

// deleteCharacter deletes the character under the cursor
func (e *Editor) deleteCharacter() {
	if e.cursorY >= e.buffer.LineCount() {
		return
	}
	
	line := e.buffer.Line(e.cursorY)
	if e.cursorX >= len(line) {
		// At end of line, join with next line
		if e.cursorY < e.buffer.LineCount()-1 {
			nextLine := e.buffer.Line(e.cursorY + 1)
			e.buffer.DeleteLine(e.cursorY + 1)
			// Append next line content to current line
			for _, r := range nextLine {
				e.buffer.InsertRune(e.cursorY, len(e.buffer.Line(e.cursorY)), r)
			}
		}
		return
	}
	
	// Delete character at cursor
	e.buffer.DeleteRune(e.cursorY, e.cursorX+1)
	e.clampCursor()
	e.adjustScroll()
}

// deleteCurrentLine deletes the entire current line
func (e *Editor) deleteCurrentLine() {
	if e.buffer.LineCount() == 1 {
		// Last line, just clear it
		line := e.buffer.Line(0)
		for i := len(line); i > 0; i-- {
			e.buffer.DeleteRune(0, i)
		}
		e.cursorX = 0
		return
	}
	
	e.buffer.DeleteLine(e.cursorY)
	
	// Adjust cursor position
	if e.cursorY >= e.buffer.LineCount() {
		e.cursorY = e.buffer.LineCount() - 1
	}
	e.cursorX = 0
	
	e.clampCursor()
	e.adjustScroll()
}
