package editor

import "unicode"

// getClosingChar returns the closing character for an opening bracket/quote
func getClosingChar(r rune) (rune, bool) {
	switch r {
	case '(':
		return ')', true
	case '[':
		return ']', true
	case '{':
		return '}', true
	case '"':
		return '"', true
	case '\'':
		return '\'', true
	case '`':
		return '`', true
	default:
		return 0, false
	}
}

// shouldAutoClose checks if we should auto-close the bracket/quote
func (e *Editor) shouldAutoClose(r rune) bool {
	if e.cursorY >= e.buffer.LineCount() {
		return false
	}
	
	line := e.buffer.Line(e.cursorY)
	
	// If at end of line, auto-close
	if e.cursorX >= len(line) {
		return true
	}
	
	// Check next character
	nextChar := rune(line[e.cursorX])
	
	// Auto-close if next char is whitespace or closing bracket
	return unicode.IsSpace(nextChar) || 
		nextChar == ')' || 
		nextChar == ']' || 
		nextChar == '}' || 
		nextChar == '"' || 
		nextChar == '\'' || 
		nextChar == '`'
}
