package editor

import "unicode"

// moveWordForward moves cursor to the start of the next word
func (e *Editor) moveWordForward() {
	if e.cursorY >= e.buffer.LineCount() {
		return
	}

	line := e.buffer.Line(e.cursorY)
	runes := []rune(line)

	// If at end of line, move to next line
	if e.cursorX >= len(runes) {
		if e.cursorY < e.buffer.LineCount()-1 {
			e.cursorY++
			e.cursorX = 0
			e.adjustScroll()
		}
		return
	}

	// Skip current word
	for e.cursorX < len(runes) && !unicode.IsSpace(runes[e.cursorX]) {
		e.cursorX++
	}

	// Skip whitespace
	for e.cursorX < len(runes) && unicode.IsSpace(runes[e.cursorX]) {
		e.cursorX++
	}

	// If we reached end of line, stay there
	if e.cursorX >= len(runes) {
		e.cursorX = len(runes)
	}

	e.adjustScroll()
}

// moveWordBackward moves cursor to the start of the previous word
func (e *Editor) moveWordBackward() {
	if e.cursorY >= e.buffer.LineCount() {
		return
	}

	line := e.buffer.Line(e.cursorY)
	runes := []rune(line)

	// If at start of line, move to previous line
	if e.cursorX == 0 {
		if e.cursorY > 0 {
			e.cursorY--
			e.cursorX = lineRuneCount(e.buffer.Line(e.cursorY))
			e.adjustScroll()
		}
		return
	}

	// Move back one position
	e.cursorX--

	// Skip whitespace
	for e.cursorX > 0 && unicode.IsSpace(runes[e.cursorX]) {
		e.cursorX--
	}

	// Skip to start of word
	for e.cursorX > 0 && !unicode.IsSpace(runes[e.cursorX-1]) {
		e.cursorX--
	}

	e.adjustScroll()
}
