package editor

import (
	"fmt"
	"strings"

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
	maxIndent := e.getMaxIndentInView()
	matchLine, matchCol := e.findMatchingBracket(e.cursorY, e.cursorX)
	
	for i := 0; i < contentHeight; i++ {
		lineNum := e.offsetY + i
		if lineNum >= e.buffer.LineCount() {
			e.term.DrawText(0, i, "~", tcell.StyleDefault.Foreground(tcell.ColorBlue))
		} else {
			line := e.buffer.Line(lineNum)
			e.renderLine(i, line)
			e.drawIndentGuides(i, line, maxIndent)
			
			if matchLine == lineNum && matchCol >= 0 {
				e.highlightBracket(matchCol, i)
			}
		}
	}
	
	e.renderStatusLine()
	
	// Position cursor
	screenY := e.cursorY - e.offsetY
	screenX := e.cursorX - e.offsetX
	if screenY < contentHeight && screenX >= 0 && screenX < e.width {
		e.term.SetCell(screenX, screenY, ' ', tcell.StyleDefault.Reverse(true))
		
		currentLine := e.buffer.Line(e.cursorY)
		if e.cursorX < len(currentLine) && isBracket(rune(currentLine[e.cursorX])) {
			style := tcell.StyleDefault.Background(tcell.NewRGBColor(100, 100, 150))
			e.term.SetCell(screenX, screenY, rune(currentLine[e.cursorX]), style)
		}
	}
	
	e.term.Show()
}

func (e *Editor) highlightBracket(x, y int) {
	line := e.buffer.Line(e.offsetY + y)
	if x >= e.offsetX && x < e.offsetX+e.width && x < len(line) {
		screenX := x - e.offsetX
		style := tcell.StyleDefault.Background(tcell.NewRGBColor(100, 100, 150))
		e.term.SetCell(screenX, y, rune(line[x]), style)
	}
}

func (e *Editor) renderLine(y int, line string) {
	lineNum := e.offsetY + y
	styledRunes := e.syntax.HighlightLine(lineNum, line, e.buffer)
	
	// Apply horizontal offset
	visibleStart := e.offsetX
	visibleEnd := e.offsetX + e.width
	
	if styledRunes == nil || len(styledRunes) == 0 {
		// Plain text rendering with horizontal scroll
		runes := []rune(line)
		for x := 0; x < e.width && visibleStart+x < len(runes); x++ {
			e.term.SetCell(x, y, runes[visibleStart+x], tcell.StyleDefault)
		}
		e.highlightSearchMatches(y, lineNum, line)
		return
	}
	
	// Render styled runes with horizontal scroll
	for i, sr := range styledRunes {
		if i >= visibleStart && i < visibleEnd {
			screenX := i - visibleStart
			e.term.SetCell(screenX, y, sr.Rune, sr.Style)
		}
	}
	
	// Highlight search matches on top of syntax highlighting
	e.highlightSearchMatches(y, lineNum, line)
}

func (e *Editor) highlightSearchMatches(y, lineNum int, line string) {
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
				screenX := match.Col + i - e.offsetX
				if screenX >= 0 && screenX < e.width {
					e.term.SetCell(screenX, y, rune(line[match.Col+i]), style)
				}
			}
		}
	}
}

func (e *Editor) renderStatusLine() {
	y := e.height - 1
	style := tcell.StyleDefault.Reverse(true)
	
	for x := 0; x < e.width; x++ {
		e.term.SetCell(x, y, ' ', style)
	}
	
	if e.mode == ModeCommand {
		cmd := ":" + e.commandBuf
		e.term.DrawText(0, y, cmd, style)
		return
	}
	
	if e.mode == ModeSearch {
		search := "/" + e.searchBuf
		e.term.DrawText(0, y, search, style)
		return
	}
	
	// Show PREVIEW mode when preview is active
	var mode string
	if e.preview.IsEnabled() {
		mode = "PREVIEW"
	} else {
		mode = e.mode.String()
	}
	e.term.DrawText(0, y, " "+mode+" ", style)
	modeWidth := len(mode) + 2
	
	if e.message != "" {
		// Check if message is a file info message (contains KB/MB and "lines")
		if strings.Contains(e.message, " lines") && (strings.Contains(e.message, "KB") || strings.Contains(e.message, "MB") || strings.Contains(e.message, "GB") || strings.Contains(e.message, " B,")) {
			e.renderFileInfoMessage(y, style, modeWidth)
		} else {
			e.term.DrawText(modeWidth+1, y, e.message, style)
		}
		return
	}
	
	filename := e.buffer.Filename()
	if filename == "" {
		filename = "[No Name]"
	}
	modified := ""
	if e.buffer.IsModified() {
		modified = " [+]"
	}
	info := filename + modified
	e.term.DrawText(modeWidth+1, y, info, style)
	
	// Don't show cursor position in preview mode
	if !e.preview.IsEnabled() {
		pos := fmt.Sprintf(" %d,%d ", e.cursorY+1, e.cursorX+1)
		e.term.DrawText(e.width-len(pos), y, pos, style)
	}
}

func (e *Editor) renderFileInfoMessage(y int, style tcell.Style, modeWidth int) {
	// Parse message: "filename" size, lines
	parts := strings.SplitN(e.message, "\"", 3)
	if len(parts) < 3 {
		e.term.DrawText(modeWidth+1, y, e.message, style)
		return
	}
	
	filename := parts[1]
	rest := strings.TrimSpace(parts[2])
	
	// Draw filename after mode
	e.term.DrawText(modeWidth+1, y, "\""+filename+"\"", style)
	
	// Draw size and lines on right
	e.term.DrawText(e.width-len(rest)-1, y, rest, style)
}
