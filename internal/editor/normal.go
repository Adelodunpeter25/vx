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
