package editor

import (
	"github.com/gdamore/tcell/v2"
)

func (e *Editor) render() {
	e.term.Clear()
	
	// If preview is enabled, show full-screen preview
	if e.preview.IsEnabled() {
		previewHeight := e.height - 1 // Reserve 1 line for status
		e.preview.Update(e.buffer)
		e.preview.Render(e.term, 0, previewHeight)
		e.renderStatusLine()
		e.term.Show()
		return
	}
	
	// Normal editor rendering
	contentHeight := e.height - 1
	gutterWidth := e.getGutterWidth()
	maxIndent := e.getMaxIndentInView()
	matchLine, matchCol := e.findMatchingBracket(e.cursorY, e.cursorX)
	
	// Render line numbers
	e.renderLineNumbers(contentHeight)
	
	for i := 0; i < contentHeight; i++ {
		lineNum := e.offsetY + i
		if lineNum >= e.buffer.LineCount() {
			e.term.DrawText(gutterWidth, i, "~", tcell.StyleDefault.Foreground(tcell.ColorBlue))
		} else {
			line := e.buffer.Line(lineNum)
			e.renderLine(i, line, gutterWidth)
			e.drawIndentGuides(i, line, maxIndent, gutterWidth)
			
			if matchLine == lineNum && matchCol >= 0 {
				e.highlightBracket(matchCol, i, gutterWidth)
			}
		}
	}
	
	e.renderStatusLine()
	
	// Position cursor
	screenY := e.cursorY - e.offsetY
	screenX := e.cursorX - e.offsetX + gutterWidth
	
	// Debug: ensure cursor is always visible
	if screenX < gutterWidth || screenX >= e.width {
		// Cursor would be off-screen, force adjust scroll
		e.adjustScroll()
		screenX = e.cursorX - e.offsetX + gutterWidth
	}
	
	if screenY >= 0 && screenY < contentHeight && screenX >= gutterWidth && screenX < e.width {
		e.term.SetCell(screenX, screenY, ' ', tcell.StyleDefault.Reverse(true))
		
		currentLine := e.buffer.Line(e.cursorY)
		if e.cursorX < len(currentLine) && isBracket(rune(currentLine[e.cursorX])) {
			// Bright yellow background for cursor bracket
			style := tcell.StyleDefault.Background(tcell.NewRGBColor(255, 200, 0)).Foreground(tcell.ColorBlack).Bold(true)
			e.term.SetCell(screenX, screenY, rune(currentLine[e.cursorX]), style)
		}
	}
	
	e.term.Show()
}

func (e *Editor) highlightBracket(x, y, gutterWidth int) {
	line := e.buffer.Line(e.offsetY + y)
	if x >= e.offsetX && x < e.offsetX+e.width-gutterWidth && x < len(line) {
		screenX := x - e.offsetX + gutterWidth
		// Bright yellow background for matching bracket
		style := tcell.StyleDefault.Background(tcell.NewRGBColor(255, 200, 0)).Foreground(tcell.ColorBlack).Bold(true)
		e.term.SetCell(screenX, y, rune(line[x]), style)
	}
}

func (e *Editor) renderLine(y int, line string, gutterWidth int) {
	lineNum := e.offsetY + y
	styledRunes := e.syntax.HighlightLine(lineNum, line, e.buffer)
	
	// Apply horizontal offset
	visibleStart := e.offsetX
	visibleEnd := e.offsetX + e.width - gutterWidth
	
	if styledRunes == nil || len(styledRunes) == 0 {
		// Plain text rendering with horizontal scroll
		runes := []rune(line)
		for x := 0; x < e.width-gutterWidth && visibleStart+x < len(runes); x++ {
			e.term.SetCell(x+gutterWidth, y, runes[visibleStart+x], tcell.StyleDefault)
		}
		e.highlightSearchMatches(y, lineNum, line, gutterWidth)
		return
	}
	
	// Render styled runes with horizontal scroll
	for i, sr := range styledRunes {
		if i >= visibleStart && i < visibleEnd {
			screenX := i - visibleStart + gutterWidth
			e.term.SetCell(screenX, y, sr.Rune, sr.Style)
		}
	}
	
	// Highlight search matches on top of syntax highlighting
	e.highlightSearchMatches(y, lineNum, line, gutterWidth)
}

func (e *Editor) highlightSearchMatches(y, lineNum int, line string, gutterWidth int) {
	if !e.search.HasMatches() {
		return
	}
	
	// All matches: dark gray background with white text
	highlightStyle := tcell.StyleDefault.
		Background(tcell.NewRGBColor(80, 80, 80)).
		Foreground(tcell.ColorWhite).
		Bold(true)
	
	// Current match: orange background with black text (better contrast than pure yellow)
	currentStyle := tcell.StyleDefault.
		Background(tcell.NewRGBColor(255, 180, 0)).
		Foreground(tcell.ColorBlack).
		Bold(true)
	
	for _, match := range e.search.GetMatches() {
		if match.Line == lineNum {
			isCurrent := e.search.Current() != nil && 
				match.Line == e.search.Current().Line && 
				match.Col == e.search.Current().Col
			
			style := highlightStyle
			if isCurrent {
				style = currentStyle
			}
			
			for i := 0; i < match.Len && match.Col+i < len(line); i++ {
				screenX := match.Col + i - e.offsetX + gutterWidth
				if screenX >= gutterWidth && screenX < e.width {
					e.term.SetCell(screenX, y, rune(line[match.Col+i]), style)
				}
			}
		}
	}
}
