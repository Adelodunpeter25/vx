package editor

import "github.com/gdamore/tcell/v2"

func (e *Editor) handleInsertMode(ev *terminal.Event) {
	if ev.Key == tcell.KeyEscape {
		e.mode = ModeNormal
		if e.cursorX > 0 {
			e.cursorX--
		}
		return
	}
	
	if ev.Key == tcell.KeyEnter {
		e.buffer.SplitLine(e.cursorY, e.cursorX)
		e.cursorY++
		e.cursorX = 0
		e.adjustScroll()
		return
	}
	
	if ev.Key == tcell.KeyBackspace || ev.Key == tcell.KeyBackspace2 {
		if e.cursorX > 0 {
			e.buffer.DeleteRune(e.cursorY, e.cursorX)
			e.cursorX--
		} else if e.cursorY > 0 {
			prevLen := len(e.buffer.Line(e.cursorY - 1))
			e.buffer.JoinLine(e.cursorY - 1)
			e.cursorY--
			e.cursorX = prevLen
			e.adjustScroll()
		}
		return
	}
	
	switch ev.Key {
	case tcell.KeyLeft:
		if e.cursorX > 0 {
			e.cursorX--
		}
		return
	case tcell.KeyRight:
		line := e.buffer.Line(e.cursorY)
		if e.cursorX < len(line) {
			e.cursorX++
		}
		return
	case tcell.KeyUp:
		if e.cursorY > 0 {
			e.cursorY--
			e.adjustScroll()
			e.clampCursor()
		}
		return
	case tcell.KeyDown:
		if e.cursorY < e.buffer.LineCount()-1 {
			e.cursorY++
			e.adjustScroll()
			e.clampCursor()
		}
		return
	}
	
	if ev.Rune != 0 {
		e.buffer.InsertRune(e.cursorY, e.cursorX, ev.Rune)
		e.cursorX++
	}
}
