package editor

import (
	"fmt"

	splitpane "github.com/Adelodunpeter25/vx/internal/split-pane"
	"github.com/gdamore/tcell/v2"
)

func (e *Editor) getGutterWidth() int {
	return e.getGutterWidthFor(e.active())
}

func (e *Editor) getGutterWidthFor(p *Pane) int {
	if p == nil {
		return 2
	}
	lineCount := p.buffer.LineCount()
	return len(fmt.Sprintf("%d", lineCount)) + 2 // diff marker + spacing
}

func (e *Editor) renderLineNumbers(contentHeight int) {
	gutterWidth := e.getGutterWidth()
	style := tcell.StyleDefault.Foreground(tcell.NewRGBColor(100, 100, 100))

	for i := 0; i < contentHeight; i++ {
		lineNum := e.active().offsetY + i
		if lineNum >= e.active().buffer.LineCount() {
			break
		}

		numStr := fmt.Sprintf("%*d ", gutterWidth-1, lineNum+1)
		for x, r := range numStr {
			e.term.SetCell(x, i, r, style)
		}
	}
}

// renderLineNumber renders a single line number at the given screen row
func (e *Editor) renderLineNumber(screenRow, lineNum, gutterWidth int) {
	e.renderLineNumberAt(splitpane.Rect{X: 0, Y: 0, Width: e.width, Height: e.height - 1}, screenRow, lineNum, gutterWidth)
}
