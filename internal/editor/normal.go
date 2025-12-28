package editor

import (
	"github.com/Adelodunpeter25/vx/internal/terminal"
	"github.com/gdamore/tcell/v2"
)

func (e *Editor) handleNormalMode(ev *terminal.Event) {
	// Clear temporary messages on any key
	if e.message == "Top of file" || e.message == "End of file" {
		e.message = ""
	}
	
	// Ctrl+C force quit
	if ev.Key == tcell.KeyCtrlC {
		e.quit = true
		return
	}
	
	switch ev.Rune {
	case 'q':
		e.quit = true
	case 'i':
		e.mode = ModeInsert
	case ':':
		e.mode = ModeCommand
		e.commandBuf = ""
		e.message = ""
	case '/':
		e.mode = ModeSearch
		e.searchBuf = ""
		e.message = ""
	case 'n':
		e.searchNext()
	case 'N':
		e.searchPrevious()
	case 'h':
		if e.cursorX > 0 {
			e.cursorX--
		}
	case 'j':
		if e.cursorY < e.buffer.LineCount()-1 {
			e.cursorY++
			e.adjustScroll()
			e.clampCursor()
		} else {
			e.message = "End of file"
		}
	case 'k':
		if e.cursorY > 0 {
			e.cursorY--
			e.adjustScroll()
			e.clampCursor()
		} else {
			e.message = "Top of file"
		}
	case 'l':
		line := e.buffer.Line(e.cursorY)
		if e.cursorX < len(line) {
			e.cursorX++
		}
	}
	
	switch ev.Key {
	case tcell.KeyLeft:
		if e.cursorX > 0 {
			e.cursorX--
		}
	case tcell.KeyRight:
		line := e.buffer.Line(e.cursorY)
		if e.cursorX < len(line) {
			e.cursorX++
		}
	case tcell.KeyUp:
		if e.cursorY > 0 {
			e.cursorY--
			e.adjustScroll()
			e.clampCursor()
		} else {
			e.message = "Top of file"
		}
	case tcell.KeyDown:
		if e.cursorY < e.buffer.LineCount()-1 {
			e.cursorY++
			e.adjustScroll()
			e.clampCursor()
		} else {
			e.message = "End of file"
		}
	}
}

func (e *Editor) searchNext() {
	if !e.search.HasMatches() {
		e.message = "No search results"
		return
	}
	
	match := e.search.Next()
	if match != nil {
		e.cursorY = match.Line
		e.cursorX = match.Col
		e.adjustScroll()
		e.message = ""
	}
}

func (e *Editor) searchPrevious() {
	if !e.search.HasMatches() {
		e.message = "No search results"
		return
	}
	
	match := e.search.Previous()
	if match != nil {
		e.cursorY = match.Line
		e.cursorX = match.Col
		e.adjustScroll()
		e.message = ""
	}
}
