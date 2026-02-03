package editor

import (
	"fmt"

	splitpane "github.com/Adelodunpeter25/vx/internal/split-pane"
	"github.com/Adelodunpeter25/vx/internal/wrap"
	"github.com/gdamore/tcell/v2"
)

func (e *Editor) render() {
	e.term.Clear()

	contentHeight := e.height - 1
	contentX := 0
	contentWidth := e.width
	if e.fileBrowser != nil && e.fileBrowser.Open {
		if e.fileBrowser.Width < 10 {
			e.fileBrowser.Width = 10
		}
		if e.fileBrowser.Width > e.width-10 {
			e.fileBrowser.Width = e.width - 10
		}
		e.fileBrowser.Render(e.term, 0, 0, e.fileBrowser.Width, contentHeight)
		contentX = e.fileBrowser.Width
		contentWidth = e.width - contentX
	}
	rects, dividerX := splitpane.LayoutSideBySide(contentWidth, contentHeight, len(e.panes), e.splitRatio)
	for i, rect := range rects {
		if i < len(e.panes) {
			isActive := i == e.activePane
			rect.X += contentX
			e.renderPane(e.panes[i], rect, isActive)
		}
	}

	// Draw divider for 2-pane layout
	if dividerX >= 0 {
		style := tcell.StyleDefault.Foreground(tcell.ColorGray)
		if e.dragSplit {
			style = style.Bold(true)
		}
		for y := 0; y < contentHeight; y++ {
			e.term.SetCell(contentX+dividerX, y, '│', style)
		}
	}

	e.renderStatusLine()
	e.term.Show()
}

func (e *Editor) renderPane(p *Pane, rect splitpane.Rect, isActive bool) {
	if p == nil {
		return
	}

	p.viewX = rect.X
	p.viewY = rect.Y
	p.viewWidth = rect.Width
	p.viewHeight = rect.Height

	// If preview is enabled, show preview within pane rect
	if p.preview.IsEnabled() {
		p.preview.Update(p.buffer)
		e.renderPreviewPane(p, rect)
		return
	}

	contentHeight := rect.Height
	gutterWidth := e.getGutterWidthFor(p)
	maxWidth := rect.Width - gutterWidth
	if maxWidth < 1 {
		return
	}

	matchLine, matchCol := e.findMatchingBracket(p.cursorY, p.cursorX)
	cursorScreenY, cursorScreenX := e.getCursorScreenPosFor(p, gutterWidth, maxWidth)

	screenRow := 0
	lineNum := p.offsetY
	visualRowsBeforeOffset := 0
	for i := 0; i < p.offsetY; i++ {
		line := p.buffer.Line(i)
		visualRowsBeforeOffset += wrap.VisualLineCount(line, maxWidth)
	}
	skipRows := p.visualOffsetY - visualRowsBeforeOffset

	for screenRow < contentHeight && lineNum < p.buffer.LineCount() {
		line := p.buffer.Line(lineNum)
		segments := wrap.WrapLine(line, lineNum, maxWidth)

		for segIdx, seg := range segments {
			if lineNum == p.offsetY && segIdx < skipRows {
				continue
			}
			if screenRow >= contentHeight {
				break
			}

			// Line numbers
			if segIdx == skipRows && lineNum == p.offsetY {
				e.renderLineNumberAt(rect, screenRow, lineNum, gutterWidth)
			} else if !seg.IsWrapped && lineNum > p.offsetY {
				e.renderLineNumberAt(rect, screenRow, lineNum, gutterWidth)
			}

			e.renderWrappedSegmentAt(rect, p, screenRow, lineNum, seg, gutterWidth)

			if p.selection.IsActive() {
				e.highlightSelectionAt(rect, p, screenRow, lineNum, seg, gutterWidth)
			}

			if matchLine == lineNum && matchCol >= seg.StartCol && matchCol < seg.StartCol+len([]rune(seg.Text)) {
				e.highlightBracketWrappedAt(rect, matchCol-seg.StartCol, screenRow, gutterWidth, line, matchCol)
			}

			screenRow++
		}
		lineNum++
		skipRows = 0
	}

	// Fill remaining rows with ~
	for screenRow < contentHeight {
		e.drawTextAt(rect, gutterWidth, screenRow, "~", tcell.StyleDefault.Foreground(tcell.ColorBlue))
		screenRow++
	}

	// Cursor (only active pane)
	if isActive && p.mode != ModeSearch && cursorScreenY >= 0 && cursorScreenY < contentHeight && cursorScreenX >= gutterWidth && cursorScreenX < rect.Width {
		e.setCellAt(rect, cursorScreenX, cursorScreenY, ' ', tcell.StyleDefault.Reverse(true))
		currentLine := []rune(p.buffer.Line(p.cursorY))
		if p.cursorX < len(currentLine) && isBracket(currentLine[p.cursorX]) {
			style := tcell.StyleDefault.Background(tcell.NewRGBColor(255, 200, 0)).Foreground(tcell.ColorBlack).Bold(true)
			e.setCellAt(rect, cursorScreenX, cursorScreenY, currentLine[p.cursorX], style)
		}
	}
}

func (e *Editor) renderPreviewPane(p *Pane, rect splitpane.Rect) {
	p.preview.Render(e.term, rect.X, rect.Y, rect.Height, rect.Width)
}

func (e *Editor) setCellAt(rect splitpane.Rect, x, y int, r rune, style tcell.Style) {
	e.term.SetCell(rect.X+x, rect.Y+y, r, style)
}

func (e *Editor) drawTextAt(rect splitpane.Rect, x, y int, text string, style tcell.Style) {
	for i, r := range text {
		e.setCellAt(rect, x+i, y, r, style)
	}
}

func (e *Editor) renderLineNumberAt(rect splitpane.Rect, screenRow, lineNum, gutterWidth int) {
	style := tcell.StyleDefault.Foreground(tcell.NewRGBColor(100, 100, 100))
	numStr := fmt.Sprintf("%*d ", gutterWidth-1, lineNum+1)
	for x, r := range numStr {
		e.setCellAt(rect, x, screenRow, r, style)
	}
}

func (e *Editor) getCursorScreenPos(gutterWidth, maxWidth int) (screenY, screenX int) {
	return e.getCursorScreenPosFor(e.active(), gutterWidth, maxWidth)
}

func (e *Editor) getCursorScreenPosFor(p *Pane, gutterWidth, maxWidth int) (screenY, screenX int) {
	if p == nil {
		return 0, 0
	}
	// Calculate cursor's visual row position
	cursorVisualLine := 0
	for lineNum := 0; lineNum < p.cursorY && lineNum < p.buffer.LineCount(); lineNum++ {
		line := p.buffer.Line(lineNum)
		cursorVisualLine += wrap.VisualLineCount(line, maxWidth)
	}

	// Find which wrapped segment contains the cursor
	currentLine := p.buffer.Line(p.cursorY)
	segments := wrap.WrapLine(currentLine, p.cursorY, maxWidth)

	for i, seg := range segments {
		segEndCol := seg.StartCol + len([]rune(seg.Text))
		if p.cursorX >= seg.StartCol && p.cursorX <= segEndCol {
			cursorVisualLine += i
			screenX = (p.cursorX - seg.StartCol) + gutterWidth
			break
		}
	}

	// Convert to screen position relative to visual offset
	screenY = cursorVisualLine - p.visualOffsetY

	return screenY, screenX
}

func (e *Editor) renderWrappedSegment(screenRow, lineNum int, seg wrap.Line, gutterWidth int) {
	p := e.active()
	e.renderWrappedSegmentAt(splitpane.Rect{X: 0, Y: 0, Width: e.width, Height: e.height - 1}, p, screenRow, lineNum, seg, gutterWidth)
}

func (e *Editor) renderWrappedSegmentAt(rect splitpane.Rect, p *Pane, screenRow, lineNum int, seg wrap.Line, gutterWidth int) {
	styledRunes := p.syntax.HighlightLine(lineNum, p.buffer.Line(lineNum), p.buffer)

	runes := []rune(seg.Text)
	for i, r := range runes {
		bufferCol := seg.StartCol + i
		style := tcell.StyleDefault

		// Apply syntax highlighting if available
		if styledRunes != nil && bufferCol < len(styledRunes) {
			style = styledRunes[bufferCol].Style
		}

		e.setCellAt(rect, gutterWidth+i, screenRow, r, style)
	}

	// Highlight search matches
	e.highlightSearchMatchesWrappedAt(rect, p, screenRow, lineNum, seg, gutterWidth)
}

func (e *Editor) highlightBracketWrapped(screenCol, screenRow, gutterWidth int, line string, bufferCol int) {
	e.highlightBracketWrappedAt(splitpane.Rect{X: 0, Y: 0, Width: e.width, Height: e.height - 1}, screenCol, screenRow, gutterWidth, line, bufferCol)
}

func (e *Editor) highlightBracketWrappedAt(rect splitpane.Rect, screenCol, screenRow, gutterWidth int, line string, bufferCol int) {
	runes := []rune(line)
	if bufferCol < len(runes) {
		style := tcell.StyleDefault.Background(tcell.NewRGBColor(255, 200, 0)).Foreground(tcell.ColorBlack).Bold(true)
		e.setCellAt(rect, gutterWidth+screenCol, screenRow, runes[bufferCol], style)
	}
}

func (e *Editor) highlightSearchMatchesWrapped(screenRow, lineNum int, seg wrap.Line, gutterWidth int) {
	p := e.active()
	e.highlightSearchMatchesWrappedAt(splitpane.Rect{X: 0, Y: 0, Width: e.width, Height: e.height - 1}, p, screenRow, lineNum, seg, gutterWidth)
}

func (e *Editor) highlightSearchMatchesWrappedAt(rect splitpane.Rect, p *Pane, screenRow, lineNum int, seg wrap.Line, gutterWidth int) {
	if !p.search.HasMatches() {
		return
	}

	line := []rune(p.buffer.Line(lineNum))

	highlightStyle := tcell.StyleDefault.
		Background(tcell.NewRGBColor(80, 80, 80)).
		Foreground(tcell.ColorWhite).
		Bold(true)

	currentStyle := tcell.StyleDefault.
		Background(tcell.NewRGBColor(0, 200, 200)).
		Foreground(tcell.ColorBlack).
		Bold(true)

	for _, match := range p.search.GetMatches() {
		if match.Line != lineNum {
			continue
		}

		isCurrent := p.search.Current() != nil &&
			match.Line == p.search.Current().Line &&
			match.Col == p.search.Current().Col

		style := highlightStyle
		if isCurrent {
			style = currentStyle
		}

		// Check if match overlaps with this segment
		for i := 0; i < match.Len && match.Col+i < len(line); i++ {
			bufferCol := match.Col + i
			if bufferCol >= seg.StartCol && bufferCol < seg.StartCol+len([]rune(seg.Text)) {
				screenCol := bufferCol - seg.StartCol
				e.setCellAt(rect, gutterWidth+screenCol, screenRow, line[bufferCol], style)
			}
		}
	}
}

func (e *Editor) highlightBracket(x, y, gutterWidth int) {
	// Legacy function - kept for compatibility
	p := e.active()
	line := p.buffer.Line(p.offsetY + y)
	runes := []rune(line)
	if x < len(runes) {
		style := tcell.StyleDefault.Background(tcell.NewRGBColor(255, 200, 0)).Foreground(tcell.ColorBlack).Bold(true)
		e.term.SetCell(gutterWidth+x, y, runes[x], style)
	}
}

func (e *Editor) renderLine(y int, line string, gutterWidth int) {
	p := e.active()
	lineNum := p.offsetY + y
	styledRunes := p.syntax.HighlightLine(lineNum, line, p.buffer)

	// Apply horizontal offset
	visibleStart := p.offsetX
	visibleEnd := p.offsetX + e.width - gutterWidth

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
	p := e.active()
	if !p.search.HasMatches() {
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

	for _, match := range p.search.GetMatches() {
		if match.Line == lineNum {
			isCurrent := p.search.Current() != nil &&
				match.Line == p.search.Current().Line &&
				match.Col == p.search.Current().Col

			style := highlightStyle
			if isCurrent {
				style = currentStyle
			}

			lineRunes := []rune(line)
			for i := 0; i < match.Len && match.Col+i < len(lineRunes); i++ {
				screenX := match.Col + i - p.offsetX + gutterWidth
				if screenX >= gutterWidth && screenX < e.width {
					e.term.SetCell(screenX, y, lineRunes[match.Col+i], style)
				}
			}
		}
	}
}
