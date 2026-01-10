package visual

// GetNormalizedSelection returns selection with start before end
func (m *Manager) GetNormalizedSelection() (startLine, startCol, endLine, endCol int) {
	sel := m.selection
	
	// Ensure start is before end
	if sel.StartLine < sel.EndLine || 
	   (sel.StartLine == sel.EndLine && sel.StartCol <= sel.EndCol) {
		return sel.StartLine, sel.StartCol, sel.EndLine, sel.EndCol
	}
	
	return sel.EndLine, sel.EndCol, sel.StartLine, sel.StartCol
}

// IsPositionSelected returns true if the given position is within the selection
func (m *Manager) IsPositionSelected(line, col int) bool {
	if !m.selection.Active {
		return false
	}
	
	startLine, startCol, endLine, endCol := m.GetNormalizedSelection()
	
	// Single line selection
	if startLine == endLine {
		return line == startLine && col >= startCol && col < endCol
	}
	
	// Multi-line selection
	if line == startLine {
		return col >= startCol
	}
	if line == endLine {
		return col < endCol
	}
	if line > startLine && line < endLine {
		return true
	}
	
	return false
}

// GetSelectedText returns the selected text from buffer lines
func (m *Manager) GetSelectedText(lines []string) string {
	if !m.HasSelection() {
		return ""
	}
	
	startLine, startCol, endLine, endCol := m.GetNormalizedSelection()
	
	// Bounds check
	if startLine >= len(lines) || endLine >= len(lines) {
		return ""
	}
	
	// Single line selection
	if startLine == endLine {
		line := lines[startLine]
		if startCol >= len(line) {
			return ""
		}
		if endCol > len(line) {
			endCol = len(line)
		}
		return line[startCol:endCol]
	}
	
	// Multi-line selection
	var result string
	
	// First line
	firstLine := lines[startLine]
	if startCol < len(firstLine) {
		result += firstLine[startCol:] + "\n"
	}
	
	// Middle lines
	for i := startLine + 1; i < endLine; i++ {
		result += lines[i] + "\n"
	}
	
	// Last line
	if endLine < len(lines) {
		lastLine := lines[endLine]
		if endCol > len(lastLine) {
			endCol = len(lastLine)
		}
		result += lastLine[:endCol]
	}
	
	return result
}
