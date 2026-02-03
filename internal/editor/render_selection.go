package editor

import (
	"github.com/Adelodunpeter25/vx/internal/wrap"
	"github.com/gdamore/tcell/v2"
)

// highlightSelection highlights the selected text on the given screen row
func (e *Editor) highlightSelection(screenRow, lineNum int, seg wrap.Line, gutterWidth int) {
	p := e.active()
	startLine, startCol, endLine, endCol, ok := p.selection.GetRange()
	if !ok {
		return
	}

	// Check if this line is within selection range
	if lineNum < startLine || lineNum > endLine {
		return
	}

	// Calculate which part of this segment is selected
	segStart := seg.StartCol
	segEnd := seg.StartCol + len([]rune(seg.Text))

	var highlightStart, highlightEnd int

	if lineNum == startLine && lineNum == endLine {
		// Selection is on single line
		highlightStart = max(startCol, segStart)
		highlightEnd = min(endCol, segEnd)
	} else if lineNum == startLine {
		// First line of multi-line selection
		highlightStart = max(startCol, segStart)
		highlightEnd = segEnd
	} else if lineNum == endLine {
		// Last line of multi-line selection
		highlightStart = segStart
		highlightEnd = min(endCol, segEnd)
	} else {
		// Middle line of multi-line selection
		highlightStart = segStart
		highlightEnd = segEnd
	}

	// Only highlight if there's overlap with this segment
	if highlightStart >= segEnd || highlightEnd <= segStart {
		return
	}

	// Apply highlight style to selected characters
	selectionStyle := tcell.StyleDefault.Background(tcell.ColorGray).Foreground(tcell.ColorBlack)
	line := []rune(p.buffer.Line(lineNum))

	for col := highlightStart; col < highlightEnd && col < len(line); col++ {
		screenX := gutterWidth + (col - segStart)
		if screenX >= gutterWidth && screenX < e.width {
			e.term.SetCell(screenX, screenRow, line[col], selectionStyle)
		}
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
