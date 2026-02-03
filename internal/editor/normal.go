package editor

import (
	"strings"

	"github.com/Adelodunpeter25/vx/internal/clipboard"
	"github.com/Adelodunpeter25/vx/internal/terminal"
	"github.com/Adelodunpeter25/vx/internal/utils"
	"github.com/gdamore/tcell/v2"
)

func (e *Editor) handleNormalMode(ev *terminal.Event) {
	p := e.active()
	// Clear transient messages on any key
	p.msgManager.ClearIfTransient()

	// If in preview mode, handle preview-specific keys
	if p.preview.IsEnabled() {
		e.handlePreviewKeys(ev)
		return
	}

	// Ctrl+C force quit
	if ev.Key == tcell.KeyCtrlC {
		e.quit = true
		return
	}

	// Ctrl+S save
	if ev.Key == tcell.KeyCtrlS {
		if p.buffer.Filename() == "" {
			p.msgManager.SetError("No filename specified")
		} else {
			if err := p.buffer.Save(); err != nil {
				p.msgManager.SetError(utils.FormatSaveError(p.buffer.Filename(), err))
			} else {
				size, _ := p.buffer.GetFileSize()
				p.msgManager.SetPersistent(utils.FormatFileInfo(p.buffer.Filename(), size, p.buffer.LineCount()))
			}
		}
		return
	}

	// Ctrl+F search
	if ev.Key == tcell.KeyCtrlF {
		p.mode = ModeSearch
		p.searchBuf = ""
		p.msgManager.Clear()
		p.lastKey = 0
		return
	}

	// Ctrl+N next buffer
	if ev.Key == tcell.KeyCtrlN {
		e.nextBuffer()
		return
	}

	// Ctrl+P previous buffer
	if ev.Key == tcell.KeyCtrlP {
		e.previousBuffer()
		return
	}

	switch ev.Rune {
	case 'q':
		e.quit = true
	case 'i':
		p.mode = ModeInsert
		p.lastKey = 0
	case ':':
		p.mode = ModeCommand
		p.commandBuf = ""
		p.msgManager.Clear()
		p.lastKey = 0
	case '/':
		p.mode = ModeSearch
		p.searchBuf = ""
		p.msgManager.Clear()
		p.lastKey = 0
	case 'H':
		// Ctrl+H for replace
		p.mode = ModeReplace
		p.replace.Start()
		p.msgManager.Clear()
		p.lastKey = 0
	case 'n':
		e.searchNext()
		p.lastKey = 0
	case 'N':
		e.searchPrevious()
		p.lastKey = 0
	case 'c':
		// Copy selection if active, otherwise copy current line
		if p.selection.IsActive() {
			e.copySelection()
		} else {
			e.copyCurrentLine()
		}
		p.lastKey = 0
	case 'x':
		// Cut selection if active, otherwise delete character
		if p.selection.IsActive() {
			e.cutSelection()
		} else {
			e.deleteCharacter()
		}
		p.lastKey = 0
	case 'd':
		// Handle dd (delete line)
		if p.lastKey == 'd' {
			e.deleteCurrentLine()
			p.lastKey = 0
		} else {
			p.lastKey = 'd'
		}
	case 'p':
		// Check if this is a markdown file
		if strings.HasSuffix(p.buffer.Filename(), ".md") {
			e.togglePreview()
		} else {
			e.pasteFromClipboard()
		}
		p.lastKey = 0
	case 'u':
		e.performUndo()
		p.lastKey = 0
	case 'r':
		e.performRedo()
		p.lastKey = 0
	case 'g':
		// Handle gg (go to start of file)
		if p.lastKey == 'g' {
			e.jumpToStart()
			p.lastKey = 0
		} else {
			p.lastKey = 'g'
		}
	case 'G':
		// Go to end of file
		e.jumpToEnd()
		p.lastKey = 0
	case 'w':
		e.moveWordForward()
		p.lastKey = 0
	case 'b':
		e.moveWordBackward()
		p.lastKey = 0
	case 'h':
		if p.cursorX > 0 {
			p.cursorX--
		}
		e.adjustScroll()
		p.selection.Clear()
		p.lastKey = 0
	case 'j':
		if p.cursorY < p.buffer.LineCount()-1 {
			p.cursorY++
			e.adjustScroll()
			e.clampCursor()
		} else {
			p.msgManager.SetTransient("End of file")
		}
		p.selection.Clear()
		p.lastKey = 0
	case 'k':
		if p.cursorY > 0 {
			p.cursorY--
			e.adjustScroll()
			e.clampCursor()
		} else {
			p.msgManager.SetTransient("Top of file")
		}
		p.selection.Clear()
		p.lastKey = 0
	case 'l':
		line := p.buffer.Line(p.cursorY)
		if p.cursorX < lineRuneCount(line) {
			p.cursorX++
		}
		e.adjustScroll()
		p.selection.Clear()
		p.lastKey = 0
	default:
		// Clear lastKey if any other key is pressed
		p.lastKey = 0
	}

	switch ev.Key {
	case tcell.KeyEscape:
		p.selection.Clear()
		p.lastKey = 0
	case tcell.KeyBackspace, tcell.KeyBackspace2:
		if p.selection.IsActive() {
			e.deleteSelection()
		}
		p.lastKey = 0
	case tcell.KeyLeft:
		if p.cursorX > 0 {
			p.cursorX--
		}
		e.adjustScroll()
		p.selection.Clear()
		p.lastKey = 0
	case tcell.KeyRight:
		line := p.buffer.Line(p.cursorY)
		if p.cursorX < lineRuneCount(line) {
			p.cursorX++
		}
		e.adjustScroll()
		p.selection.Clear()
		p.lastKey = 0
	case tcell.KeyUp:
		if p.cursorY > 0 {
			p.cursorY--
			e.adjustScroll()
			e.clampCursor()
		} else {
			p.msgManager.SetTransient("Top of file")
		}
		p.selection.Clear()
		p.lastKey = 0
	case tcell.KeyDown:
		if p.cursorY < p.buffer.LineCount()-1 {
			p.cursorY++
			e.adjustScroll()
			e.clampCursor()
		} else {
			p.msgManager.SetTransient("End of file")
		}
		p.selection.Clear()
		p.lastKey = 0
	}
}

func (e *Editor) jumpToStart() {
	p := e.active()
	p.cursorY = 0
	p.cursorX = 0
	p.offsetY = 0
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
		p.preview.Update(p.buffer)
		p.msgManager.SetTransient("Preview enabled")
	} else {
		p.msgManager.SetTransient("Preview disabled")
	}
	p.renderCache.invalidate()
}

func (e *Editor) handlePreviewKeys(ev *terminal.Event) {
	p := e.active()
	switch ev.Rune {
	case 'p':
		// Toggle preview off
		e.togglePreview()
	case 'j':
		p.preview.Scroll(1)
	case 'k':
		p.preview.Scroll(-1)
	case 'q':
		e.quit = true
	}

	switch ev.Key {
	case tcell.KeyDown:
		p.preview.Scroll(1)
	case tcell.KeyUp:
		p.preview.Scroll(-1)
	case tcell.KeyCtrlC:
		e.quit = true
	}

	// Ensure render is triggered after preview key handling
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

	// Insert text at cursor position
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
