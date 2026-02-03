package editor

// FindMatchingBracket finds the matching bracket for the one at the given position
// Returns (-1, -1) if no match found
func (e *Editor) findMatchingBracket(line, col int) (int, int) {
	p := e.active()
	if line < 0 || line >= p.buffer.LineCount() {
		return -1, -1
	}

	currentLine := p.buffer.Line(line)
	runes := []rune(currentLine)
	if col < 0 || col >= len(runes) {
		return -1, -1
	}

	char := runes[col]

	// Check if current char is a bracket
	var opening, closing rune
	var forward bool

	switch char {
	case '(':
		opening, closing, forward = '(', ')', true
	case ')':
		opening, closing, forward = '(', ')', false
	case '[':
		opening, closing, forward = '[', ']', true
	case ']':
		opening, closing, forward = '[', ']', false
	case '{':
		opening, closing, forward = '{', '}', true
	case '}':
		opening, closing, forward = '{', '}', false
	default:
		return -1, -1
	}

	if forward {
		return e.findForwardBracket(line, col, opening, closing)
	}
	return e.findBackwardBracket(line, col, opening, closing)
}

func (e *Editor) findForwardBracket(startLine, startCol int, opening, closing rune) (int, int) {
	p := e.active()
	depth := 1

	// Start from next character
	for l := startLine; l < p.buffer.LineCount(); l++ {
		line := []rune(p.buffer.Line(l))
		start := 0
		if l == startLine {
			start = startCol + 1
		}

		for c := start; c < len(line); c++ {
			char := line[c]
			if char == opening {
				depth++
			} else if char == closing {
				depth--
				if depth == 0 {
					return l, c
				}
			}
		}
	}

	return -1, -1
}

func (e *Editor) findBackwardBracket(startLine, startCol int, opening, closing rune) (int, int) {
	p := e.active()
	depth := 1

	// Start from previous character
	for l := startLine; l >= 0; l-- {
		line := []rune(p.buffer.Line(l))
		end := len(line) - 1
		if l == startLine {
			end = startCol - 1
		}

		for c := end; c >= 0; c-- {
			char := line[c]
			if char == closing {
				depth++
			} else if char == opening {
				depth--
				if depth == 0 {
					return l, c
				}
			}
		}
	}

	return -1, -1
}

func isBracket(r rune) bool {
	return r == '(' || r == ')' || r == '[' || r == ']' || r == '{' || r == '}'
}
