package editor

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
)

func (e *Editor) getGutterWidth() int {
	lineCount := e.buffer.LineCount()
	return len(fmt.Sprintf("%d", lineCount)) + 1 // +1 for spacing
}

func (e *Editor) renderLineNumbers(contentHeight int) {
	gutterWidth := e.getGutterWidth()
	style := tcell.StyleDefault.Foreground(tcell.NewRGBColor(100, 100, 100))
	
	for i := 0; i < contentHeight; i++ {
		lineNum := e.offsetY + i
		if lineNum >= e.buffer.LineCount() {
			break
		}
		
		numStr := fmt.Sprintf("%*d ", gutterWidth-1, lineNum+1)
		for x, r := range numStr {
			e.term.SetCell(x, i, r, style)
		}
	}
}
