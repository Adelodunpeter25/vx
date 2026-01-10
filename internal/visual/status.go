package visual

import "fmt"

// GetStatusInfo returns status bar information for visual mode
func (m *Manager) GetStatusInfo(lines []string) string {
	if !m.HasSelection() {
		return ""
	}
	
	startLine, startCol, endLine, endCol := m.GetNormalizedSelection()
	
	// Count lines and characters
	lineCount := endLine - startLine + 1
	charCount := 0
	
	if startLine == endLine {
		// Single line selection
		charCount = endCol - startCol
	} else {
		// Multi-line selection
		// First line
		if startLine < len(lines) {
			charCount += len(lines[startLine]) - startCol + 1 // +1 for newline
		}
		
		// Middle lines
		for i := startLine + 1; i < endLine && i < len(lines); i++ {
			charCount += len(lines[i]) + 1 // +1 for newline
		}
		
		// Last line
		if endLine < len(lines) {
			charCount += endCol
		}
	}
	
	// Format status message
	if lineCount == 1 {
		return fmt.Sprintf("VISUAL - %d chars", charCount)
	} else {
		return fmt.Sprintf("VISUAL - %d lines, %d chars", lineCount, charCount)
	}
}

// IsVisualMode returns true if in visual selection mode
func (m *Manager) IsVisualMode() bool {
	return m.selection.Active
}
