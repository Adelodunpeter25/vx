package visual

import "github.com/Adelodunpeter25/vx/internal/buffer"

// GetSelectedText extracts the selected text from the buffer
func (s *Selection) GetSelectedText(buf *buffer.Buffer) string {
	startLine, startCol, endLine, endCol, ok := s.GetRange()
	if !ok {
		return ""
	}

	// Single line selection
	if startLine == endLine {
		line := []rune(buf.Line(startLine))
		if startCol >= len(line) {
			return ""
		}
		if endCol > len(line) {
			endCol = len(line)
		}
		return string(line[startCol:endCol])
	}

	// Multi-line selection
	var result string

	// First line
	firstLine := []rune(buf.Line(startLine))
	if startCol < len(firstLine) {
		result += string(firstLine[startCol:]) + "\n"
	}

	// Middle lines
	for i := startLine + 1; i < endLine; i++ {
		result += buf.Line(i) + "\n"
	}

	// Last line
	lastLine := []rune(buf.Line(endLine))
	if endCol > len(lastLine) {
		endCol = len(lastLine)
	}
	if endCol > 0 {
		result += string(lastLine[:endCol])
	}

	return result
}

// DeleteSelectedText removes the selected text from the buffer
func (s *Selection) DeleteSelectedText(buf *buffer.Buffer) {
	startLine, startCol, endLine, endCol, ok := s.GetRange()
	if !ok {
		return
	}

	// Single line deletion
	if startLine == endLine {
		// Delete characters from startCol to endCol
		for i := startCol; i < endCol; i++ {
			buf.DeleteRune(startLine, startCol)
		}
		return
	}

	// Multi-line deletion
	// Delete from startCol to end of first line
	firstLine := []rune(buf.Line(startLine))
	for i := startCol; i < len(firstLine); i++ {
		buf.DeleteRune(startLine, startCol)
	}
	
	// Delete middle lines (from end to start to avoid index shifting)
	for i := endLine - 1; i > startLine; i-- {
		buf.DeleteLine(i)
	}
	
	// Delete from start of last line to endCol (now it's startLine+1)
	lastLine := []rune(buf.Line(startLine + 1))
	for i := 0; i < endCol && i < len(lastLine); i++ {
		buf.DeleteRune(startLine + 1, 0)
	}
	
	// Join the lines
	buf.JoinLine(startLine)
}
