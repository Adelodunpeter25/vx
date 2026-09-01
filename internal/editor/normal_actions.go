package editor

import (
	"github.com/Adelodunpeter25/vx/internal/clipboard"
	"github.com/Adelodunpeter25/vx/internal/terminal"
	"github.com/gdamore/tcell/v2"
)

func (e *Editor) jumpToStart() {
	p := e.active()
	p.cursorY = 0
	p.cursorX = 0
	p.offsetY = 0
	p.visualOffsetY = 0
	p.offsetX = 0
	p.renderCache.invalidate()
	p.msgManager.Clear()
}

func (e *Editor) jumpToEnd() {
	p := e.active()
	p.cursorY = p.buffer.LineCount() - 1
	p.cursorX = 0
	e.adjustScroll()
	p.msgManager.Clear()
}

func (e *Editor) togglePreview() {
	p := e.active()
	p.preview.Toggle()
	if p.preview.IsEnabled() {
		p.preview.Update(p.buffer, p.viewWidth)
		p.msgManager.SetTransient("Preview enabled")
	} else {
		p.msgManager.SetTransient("Preview disabled")
	}
	p.renderCache.invalidate()
}

func (e *Editor) handlePreviewKeys(ev *terminal.Event) {
	p := e.active()

	if ev.Rune != 'g' {
		p.lastKey = 0
	}

	switch ev.Rune {
	case 'p':
		e.togglePreview()
	case 'j':
		p.preview.Scroll(1)
	case 'k':
		p.preview.Scroll(-1)
	case 'q':
		e.quit = true
	case 'g':
		if p.lastKey == 'g' {
			p.preview.ScrollToStart()
			p.lastKey = 0
		} else {
			p.lastKey = 'g'
		}
		return
	case 'G':
		p.preview.ScrollToEnd()
	}

	switch ev.Key {
	case tcell.KeyDown:
		p.preview.Scroll(1)
	case tcell.KeyUp:
		p.preview.Scroll(-1)
	case tcell.KeyCtrlC:
		e.quit = true
	case tcell.KeyCtrlU:
		p.preview.ScrollPage(-1, p.viewHeight)
	case tcell.KeyCtrlD:
		p.preview.ScrollPage(1, p.viewHeight)
	}

	p.renderCache.invalidate()
}

func (e *Editor) copyCurrentLine() {
	p := e.active()
	line := p.buffer.Line(p.cursorY)
	err := clipboard.Copy(line)
	if err != nil {
		p.msgManager.SetError("Failed to copy to clipboard")
	} else {
		p.msgManager.SetTransient("Line copied to clipboard")
	}
}

func (e *Editor) pasteFromClipboard() {
	p := e.active()
	text, err := clipboard.Paste()
	if err != nil {
		p.msgManager.SetError("Failed to paste from clipboard")
		return
	}

	if text == "" {
		p.msgManager.SetTransient("Clipboard is empty")
		return
	}

	for _, r := range text {
		if r == '\n' {
			p.buffer.SplitLine(p.cursorY, p.cursorX)
			p.cursorY++
			p.cursorX = 0
		} else {
			p.buffer.InsertRune(p.cursorY, p.cursorX, r)
			p.cursorX++
		}
	}

	e.adjustScroll()
	p.msgManager.SetTransient("Pasted from clipboard")
}

func (e *Editor) searchNext() {
	p := e.active()
	if !p.search.HasMatches() {
		p.msgManager.SetTransient("No search results")
		return
	}

	match := p.search.Next()
	if match != nil {
		p.cursorY = match.Line
		p.cursorX = match.Col
		e.adjustScroll()
		p.msgManager.Clear()
	}
}

func (e *Editor) searchPrevious() {
	p := e.active()
	if !p.search.HasMatches() {
		p.msgManager.SetTransient("No search results")
		return
	}

	match := p.search.Previous()
	if match != nil {
		p.cursorY = match.Line
		p.cursorX = match.Col
		e.adjustScroll()
		p.msgManager.Clear()
	}
}

func (e *Editor) performUndo() {
	p := e.active()
	if p.buffer.Undo() {
		p.msgManager.SetTransient("Undo")
		e.clampCursor()
		e.adjustScroll()
	} else {
		p.msgManager.SetTransient("Nothing to undo")
	}
}

func (e *Editor) performRedo() {
	p := e.active()
	if p.buffer.Redo() {
		p.msgManager.SetTransient("Redo")
		e.clampCursor()
		e.adjustScroll()
	} else {
		p.msgManager.SetTransient("Nothing to redo")
	}
}
