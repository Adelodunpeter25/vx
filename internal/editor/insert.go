package editor

import (
	"github.com/Adelodunpeter25/vx/internal/terminal"
	"github.com/gdamore/tcell/v2"
)

func (e *Editor) handleInsertMode(ev *terminal.Event) {
	p := e.active()
	// Clear transient messages on any key
	p.msgManager.ClearIfTransient()

	// Ctrl+C force quit
	if ev.Key == tcell.KeyCtrlC {
		e.quit = true
		return
	}

	if ev.Key == tcell.KeyEscape {
		p.mode = ModeNormal
		if p.cursorX > 0 {
			p.cursorX--
		}
		return
	}

	if ev.Key == tcell.KeyTab {
		p.buffer.InsertRune(p.cursorY, p.cursorX, '\t')
		p.cursorX++
		return
	}

	if ev.Key == tcell.KeyEnter {
		// Get current line indentation
		currentLine := p.buffer.Line(p.cursorY)
		indent := getIndentation(currentLine)

		p.buffer.SplitLine(p.cursorY, p.cursorX)
		p.cursorY++
		p.cursorX = 0

		// Auto-indent: insert same indentation on new line
		for _, r := range indent {
			p.buffer.InsertRune(p.cursorY, p.cursorX, r)
			p.cursorX++
		}

		e.adjustScroll()
		return
	}

	if ev.Key == tcell.KeyBackspace || ev.Key == tcell.KeyBackspace2 {
		if p.cursorX > 0 {
			p.buffer.DeleteRune(p.cursorY, p.cursorX)
			p.cursorX--
			e.adjustScroll()
		} else if p.cursorY > 0 {
			prevLen := lineRuneCount(p.buffer.Line(p.cursorY - 1))
			p.buffer.JoinLine(p.cursorY - 1)
			p.cursorY--
			p.cursorX = prevLen
			e.adjustScroll()
		}
		return
	}

	switch ev.Key {
	case tcell.KeyLeft:
		if p.cursorX > 0 {
			p.cursorX--
			e.adjustScroll()
		}
		return
	case tcell.KeyRight:
		line := p.buffer.Line(p.cursorY)
		if p.cursorX < lineRuneCount(line) {
			p.cursorX++
			e.adjustScroll()
		}
		return
	case tcell.KeyUp:
		if p.cursorY > 0 {
			p.cursorY--
			e.adjustScroll()
			e.clampCursor()
		}
		return
	case tcell.KeyDown:
		if p.cursorY < p.buffer.LineCount()-1 {
			p.cursorY++
			e.adjustScroll()
			e.clampCursor()
		}
		return
	}

	if ev.Rune != 0 {
		p.buffer.InsertRune(p.cursorY, p.cursorX, ev.Rune)
		p.cursorX++
		e.adjustScroll()
	}
}

func getIndentation(line string) string {
	indent := ""
	for _, r := range line {
		if r == ' ' || r == '\t' {
			indent += string(r)
		} else {
			break
		}
	}
	return indent
}
