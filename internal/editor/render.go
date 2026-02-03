package editor

import (
	"github.com/Adelodunpeter25/vx/internal/wrap"
	"github.com/gdamore/tcell/v2"
)

func (e *Editor) render() {
	e.term.Clear()

	// If preview is enabled, show full-screen preview
	if e.preview.IsEnabled() {
		previewHeight := e.height - 1 // Reserve 1 line for status
		e.preview.Update(e.buffer)
		e.preview.Render(e.term, 0, previewHeight, e.width)
		e.renderStatusLine()
		e.term.Show()
		return
	}

	// Normal editor rendering with line wrapping
	contentHeight := e.height - 1
	gutterWidth := e.getGutterWidth()
	maxWidth := e.width - gutterWidth
	matchLine, matchCol := e.findMatchingBracket(e.cursorY, e.cursorX)

	// Calculate cursor screen position
	cursorScreenY, cursorScreenX := e.getCursorScreenPos(gutterWidth, maxWidth)

	// Render wrapped lines starting from visualOffsetY
	screenRow := 0
	lineNum := e.offsetY

	// Calculate how many visual rows to skip in the first line
	visualRowsBeforeOffset := 0
	for i := 0; i < e.offsetY; i++ {
		line := e.buffer.Line(i)
		visualRowsBeforeOffset += wrap.VisualLineCount(line, maxWidth)
	}
	skipRows := e.visualOffsetY - visualRowsBeforeOffset

	for screenRow < contentHeight && lineNum < e.buffer.LineCount() {
		line := e.buffer.Line(lineNum)
		segments := wrap.WrapLine(line, lineNum, maxWidth)

		for segIdx, seg := range segments {
			// Skip rows if we're in the first line and need to offset
			if lineNum == e.offsetY && segIdx < skipRows {
				continue
			}

			if screenRow >= contentHeight {
				break
			}

			// Render line number only on first visible segment of each line
			if segIdx == skipRows && lineNum == e.offsetY {
				e.renderLineNumber(screenRow, lineNum, gutterWidth)
			} else if !seg.IsWrapped && lineNum > e.offsetY {
				e.renderLineNumber(screenRow, lineNum, gutterWidth)
			}

			// Render the segment
			e.renderWrappedSegment(screenRow, lineNum, seg, gutterWidth)

			// Highlight selection if active
			if e.selection.IsActive() {
				e.highlightSelection(screenRow, lineNum, seg, gutterWidth)
			}

			// Highlight matching bracket if on this line
			if matchLine == lineNum && matchCol >= seg.StartCol && matchCol < seg.StartCol+len([]rune(seg.Text)) {
				e.highlightBracketWrapped(matchCol-seg.StartCol, screenRow, gutterWidth, line, matchCol)
			}

			screenRow++
		}
		lineNum++
		skipRows = 0 // Only skip rows in the first line
	}

	// Fill remaining rows with ~
	for screenRow < contentHeight {
		e.term.DrawText(gutterWidth, screenRow, "~", tcell.StyleDefault.Foreground(tcell.ColorBlue))
		screenRow++
	}

	e.renderStatusLine()

	// Position cursor (but not in search mode - cursor position is shown in status bar)
	if e.mode != ModeSearch && cursorScreenY >= 0 && cursorScreenY < contentHeight && cursorScreenX >= gutterWidth && cursorScreenX < e.width {
		e.term.SetCell(cursorScreenX, cursorScreenY, ' ', tcell.StyleDefault.Reverse(true))

		currentLine := []rune(e.buffer.Line(e.cursorY))
		if e.cursorX < len(currentLine) && isBracket(currentLine[e.cursorX]) {
			style := tcell.StyleDefault.Background(tcell.NewRGBColor(255, 200, 0)).Foreground(tcell.ColorBlack).Bold(true)
			e.term.SetCell(cursorScreenX, cursorScreenY, currentLine[e.cursorX], style)
		}
	}

	e.term.Show()
}

func (e *Editor) getCursorScreenPos(gutterWidth, maxWidth int) (screenY, screenX int) {
	// Calculate cursor's visual row position
	cursorVisualLine := 0
	for lineNum := 0; lineNum < e.cursorY && lineNum < e.buffer.LineCount(); lineNum++ {
		line := e.buffer.Line(lineNum)
		cursorVisualLine += wrap.VisualLineCount(line, maxWidth)
	}

	// Find which wrapped segment contains the cursor
	currentLine := e.buffer.Line(e.cursorY)
	segments := wrap.WrapLine(currentLine, e.cursorY, maxWidth)

	for i, seg := range segments {
		segEndCol := seg.StartCol + len([]rune(seg.Text))
		if e.cursorX >= seg.StartCol && e.cursorX <= segEndCol {
			cursorVisualLine += i
			screenX = (e.cursorX - seg.StartCol) + gutterWidth
			break
		}
	}

	// Convert to screen position relative to visual offset
	screenY = cursorVisualLine - e.visualOffsetY

	return screenY, screenX
}

func (e *Editor) renderWrappedSegment(screenRow, lineNum int, seg wrap.Line, gutterWidth int) {
	styledRunes := e.syntax.HighlightLine(lineNum, e.buffer.Line(lineNum), e.buffer)

	runes := []rune(seg.Text)
	for i, r := range runes {
		bufferCol := seg.StartCol + i
		style := tcell.StyleDefault

		// Apply syntax highlighting if available
		if styledRunes != nil && bufferCol < len(styledRunes) {
			style = styledRunes[bufferCol].Style
		}

		e.term.SetCell(gutterWidth+i, screenRow, r, style)
	}

	// Highlight search matches
	e.highlightSearchMatchesWrapped(screenRow, lineNum, seg, gutterWidth)
}

func (e *Editor) highlightBracketWrapped(screenCol, screenRow, gutterWidth int, line string, bufferCol int) {
	runes := []rune(line)
	if bufferCol < len(runes) {
		style := tcell.StyleDefault.Background(tcell.NewRGBColor(255, 200, 0)).Foreground(tcell.ColorBlack).Bold(true)
		e.term.SetCell(gutterWidth+screenCol, screenRow, runes[bufferCol], style)
	}
}

func (e *Editor) highlightSearchMatchesWrapped(screenRow, lineNum int, seg wrap.Line, gutterWidth int) {
	if !e.search.HasMatches() {
		return
	}

	line := []rune(e.buffer.Line(lineNum))

	highlightStyle := tcell.StyleDefault.
		Background(tcell.NewRGBColor(80, 80, 80)).
		Foreground(tcell.ColorWhite).
		Bold(true)

	currentStyle := tcell.StyleDefault.
		Background(tcell.NewRGBColor(0, 200, 200)).
		Foreground(tcell.ColorBlack).
		Bold(true)

	for _, match := range e.search.GetMatches() {
		if match.Line != lineNum {
			continue
		}

		isCurrent := e.search.Current() != nil &&
			match.Line == e.search.Current().Line &&
			match.Col == e.search.Current().Col

		style := highlightStyle
		if isCurrent {
			style = currentStyle
		}

		// Check if match overlaps with this segment
		for i := 0; i < match.Len && match.Col+i < len(line); i++ {
			bufferCol := match.Col + i
			if bufferCol >= seg.StartCol && bufferCol < seg.StartCol+len([]rune(seg.Text)) {
				screenCol := bufferCol - seg.StartCol
				e.term.SetCell(gutterWidth+screenCol, screenRow, line[bufferCol], style)
			}
		}
	}
}

func (e *Editor) highlightBracket(x, y, gutterWidth int) {
	// Legacy function - kept for compatibility
	line := e.buffer.Line(e.offsetY + y)
	runes := []rune(line)
	if x < len(runes) {
		style := tcell.StyleDefault.Background(tcell.NewRGBColor(255, 200, 0)).Foreground(tcell.ColorBlack).Bold(true)
		e.term.SetCell(gutterWidth+x, y, runes[x], style)
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

			lineRunes := []rune(line)
			for i := 0; i < match.Len && match.Col+i < len(lineRunes); i++ {
				screenX := match.Col + i - e.offsetX + gutterWidth
				if screenX >= gutterWidth && screenX < e.width {
					e.term.SetCell(screenX, y, lineRunes[match.Col+i], style)
				}
			}
		}
	}
}
