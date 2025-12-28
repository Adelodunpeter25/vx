package editor

import "github.com/gdamore/tcell/v2"

// GetIndentLevel returns the indentation level of a line (number of leading tabs/spaces)
func GetIndentLevel(line string) int {
	level := 0
	for _, r := range line {
		if r == '\t' {
			level++
		} else if r == ' ' {
			// Count 2 spaces as 1 indent level (configurable)
			level++
			if level%2 == 0 {
				continue
			}
		} else {
			break
		}
	}
	return level / 2 // Normalize spaces to tab equivalents
}

// DrawIndentGuides draws vertical lines at indentation boundaries
func (e *Editor) drawIndentGuides(y int, line string, maxIndent int) {
	// Draw guides for each indent level
	for i := 0; i < maxIndent; i++ {
		x := i * 2 // 2 spaces per indent level
		if x < len(line) {
			char := line[x]
			// Only draw guide if this position is whitespace
			if char == ' ' || char == '\t' {
				style := tcell.StyleDefault.Foreground(tcell.ColorGray).Dim(true)
				e.term.SetCell(x, y, '│', style)
			}
		}
	}
}

// GetMaxIndentInView returns the maximum indent level in visible lines
func (e *Editor) getMaxIndentInView() int {
	maxIndent := 0
	contentHeight := e.height - 1
	
	for i := 0; i < contentHeight; i++ {
		lineNum := e.offsetY + i
		if lineNum >= e.buffer.LineCount() {
			break
		}
		line := e.buffer.Line(lineNum)
		indent := GetIndentLevel(line)
		if indent > maxIndent {
			maxIndent = indent
		}
	}
	
	return maxIndent
}
