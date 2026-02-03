package editor

import "unicode"

// moveWordForward moves cursor to the start of the next word
func (e *Editor) moveWordForward() {
	p := e.active()
	if p.cursorY >= p.buffer.LineCount() {
		return
	}

	line := p.buffer.Line(p.cursorY)
	runes := []rune(line)

	// If at end of line, move to next line
	if p.cursorX >= len(runes) {
		if p.cursorY < p.buffer.LineCount()-1 {
			p.cursorY++
			p.cursorX = 0
			e.adjustScroll()
		}
		return
	}

	// Skip current word
	for p.cursorX < len(runes) && !unicode.IsSpace(runes[p.cursorX]) {
		p.cursorX++
	}

	// Skip whitespace
	for p.cursorX < len(runes) && unicode.IsSpace(runes[p.cursorX]) {
		p.cursorX++
	}

	// If we reached end of line, stay there
	if p.cursorX >= len(runes) {
		p.cursorX = len(runes)
	}

	e.adjustScroll()
}

// moveWordBackward moves cursor to the start of the previous word
func (e *Editor) moveWordBackward() {
	p := e.active()
	if p.cursorY >= p.buffer.LineCount() {
		return
	}

	line := p.buffer.Line(p.cursorY)
	runes := []rune(line)

	// If at start of line, move to previous line
	if p.cursorX == 0 {
		if p.cursorY > 0 {
			p.cursorY--
			p.cursorX = lineRuneCount(p.buffer.Line(p.cursorY))
			e.adjustScroll()
		}
		return
	}

	// Move back one position
	p.cursorX--

	// Skip whitespace
	for p.cursorX > 0 && unicode.IsSpace(runes[p.cursorX]) {
		p.cursorX--
	}

	// Skip to start of word
	for p.cursorX > 0 && !unicode.IsSpace(runes[p.cursorX-1]) {
		p.cursorX--
	}

	e.adjustScroll()
}
