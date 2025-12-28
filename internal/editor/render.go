package editor

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
)

func (e *Editor) render() {
	e.term.Clear()
	
	contentHeight := e.height - 1
	maxIndent := e.getMaxIndentInView()
	
	// Find matching bracket if cursor is on one
	matchLine, matchCol := e.findMatchingBracket(e.cursorY, e.cursorX)
	
	for i := 0; i < contentHeight; i++ {
		lineNum := e.offsetY + i
		if lineNum >= e.buffer.LineCount() {
			e.term.DrawText(0, i, "~", tcell.StyleDefault.Foreground(tcell.ColorBlue))
		} else {
			line := e.buffer.Line(lineNum)
			e.renderLine(i, line)
			
			// Draw indent guides after content
			e.drawIndentGuides(i, line, maxIndent)
			
			// Highlight matching bracket if on this line
			if matchLine == lineNum && matchCol >= 0 {
				e.highlightBracket(matchCol, i)
			}
		}
	}
	
	e.renderStatusLine()
	
	// Position cursor and highlight bracket under cursor
	screenY := e.cursorY - e.offsetY
	e.term.SetCell(e.cursorX, screenY, ' ', tcell.StyleDefault.Reverse(true))
	
	// Highlight current bracket if cursor is on one
	currentLine := e.buffer.Line(e.cursorY)
	if e.cursorX < len(currentLine) && isBracket(rune(currentLine[e.cursorX])) {
		style := tcell.StyleDefault.Background(tcell.NewRGBColor(100, 100, 150))
		e.term.SetCell(e.cursorX, screenY, rune(currentLine[e.cursorX]), style)
	}
	
	e.term.Show()
}

func (e *Editor) highlightBracket(x, y int) {
	line := e.buffer.Line(e.offsetY + y)
	if x < len(line) {
		style := tcell.StyleDefault.Background(tcell.NewRGBColor(100, 100, 150))
		e.term.SetCell(x, y, rune(line[x]), style)
	}
}

func (e *Editor) renderLine(y int, line string) {
	lineNum := e.offsetY + y
	styledRunes := e.syntax.HighlightLine(lineNum, line, e.buffer)
	
	if styledRunes == nil || len(styledRunes) == 0 {
		e.term.DrawText(0, y, line, tcell.StyleDefault)
		e.highlightSearchMatches(y, lineNum, line)
		return
	}
	
	for x, sr := range styledRunes {
		e.term.SetCell(x, y, sr.Rune, sr.Style)
	}
	
	// Highlight search matches on top of syntax highlighting
	e.highlightSearchMatches(y, lineNum, line)
}

func (e *Editor) highlightSearchMatches(y, lineNum int, line string) {
	if !e.search.HasMatches() {
		return
	}
	
	// All matches: light blue/cyan background
	highlightStyle := tcell.StyleDefault.
		Background(tcell.NewRGBColor(100, 180, 255)).
		Foreground(tcell.ColorBlack).
		Bold(true)
	
	// Current match: bright green background
	currentStyle := tcell.StyleDefault.
		Background(tcell.NewRGBColor(100, 255, 100)).
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
				e.term.SetCell(match.Col+i, y, rune(line[match.Col+i]), style)
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
	
	// Always show mode
	mode := e.mode.String()
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
	
	pos := fmt.Sprintf(" %d,%d ", e.cursorY+1, e.cursorX+1)
	e.term.DrawText(e.width-len(pos), y, pos, style)
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
