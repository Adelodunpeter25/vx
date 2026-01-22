package editor

import (
	"strings"

	"github.com/Adelodunpeter25/vx/internal/clipboard"
	"github.com/Adelodunpeter25/vx/internal/terminal"
	"github.com/Adelodunpeter25/vx/internal/utils"
	"github.com/gdamore/tcell/v2"
)

func (e *Editor) handleNormalMode(ev *terminal.Event) {
	// Clear temporary messages on any key
	if e.message == "Top of file" || e.message == "End of file" || 
	   e.message == "Buffer closed" || e.message == "Buffer saved and closed" || 
	   e.message == "Buffer closed without saving" || e.message == "Cannot close last buffer" {
		e.message = ""
	}
	
	// Clear file info messages on any key
	if strings.Contains(e.message, " lines") && (strings.Contains(e.message, "KB") || strings.Contains(e.message, "MB") || strings.Contains(e.message, "GB") || strings.Contains(e.message, " B,")) {
		e.message = ""
	}
	
	// If in preview mode, handle preview-specific keys
	if e.preview.IsEnabled() {
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
		if e.buffer.Filename() == "" {
			e.message = "No filename specified"
		} else {
			if err := e.buffer.Save(); err != nil {
				e.message = utils.FormatSaveError(e.buffer.Filename(), err)
			} else {
				size, _ := e.buffer.GetFileSize()
				e.message = utils.FormatFileInfo(e.buffer.Filename(), size, e.buffer.LineCount())
			}
		}
		return
	}
	
	// Ctrl+F search
	if ev.Key == tcell.KeyCtrlF {
		e.mode = ModeSearch
		e.searchBuf = ""
		e.message = ""
		e.lastKey = 0
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
		e.mode = ModeInsert
		e.lastKey = 0
	case ':':
		e.mode = ModeCommand
		e.commandBuf = ""
		e.message = ""
		e.lastKey = 0
	case '/':
		e.mode = ModeSearch
		e.searchBuf = ""
		e.message = ""
		e.lastKey = 0
	case 'H':
		// Ctrl+H for replace
		e.mode = ModeReplace
		e.replace.Start()
		e.message = ""
		e.lastKey = 0
	case 'n':
		e.searchNext()
		e.lastKey = 0
	case 'N':
		e.searchPrevious()
		e.lastKey = 0
	case 'c':
		// Copy selection if active, otherwise copy current line
		if e.selection.IsActive() {
			e.copySelection()
		} else {
			e.copyCurrentLine()
		}
		e.lastKey = 0
	case 'x':
		// Cut selection if active, otherwise delete character
		if e.selection.IsActive() {
			e.cutSelection()
		} else {
			e.deleteCharacter()
		}
		e.lastKey = 0
	case 'd':
		// Handle dd (delete line)
		if e.lastKey == 'd' {
			e.deleteCurrentLine()
			e.lastKey = 0
		} else {
			e.lastKey = 'd'
		}
	case 'p':
		// Check if this is a markdown file
		if strings.HasSuffix(e.buffer.Filename(), ".md") {
			e.togglePreview()
		} else {
			e.pasteFromClipboard()
		}
		e.lastKey = 0
	case 'u':
		e.performUndo()
		e.lastKey = 0
	case 'r':
		e.performRedo()
		e.lastKey = 0
	case 'g':
		// Handle gg (go to start of file)
		if e.lastKey == 'g' {
			e.jumpToStart()
			e.lastKey = 0
		} else {
			e.lastKey = 'g'
		}
	case 'G':
		// Go to end of file
		e.jumpToEnd()
		e.lastKey = 0
	case 'w':
		e.moveWordForward()
		e.lastKey = 0
	case 'b':
		e.moveWordBackward()
		e.lastKey = 0
	case 'h':
		if e.cursorX > 0 {
			e.cursorX--
		}
		e.adjustScroll()
		e.selection.Clear()
		e.lastKey = 0
	case 'j':
		if e.cursorY < e.buffer.LineCount()-1 {
			e.cursorY++
			e.adjustScroll()
			e.clampCursor()
		} else {
			e.message = "End of file"
		}
		e.selection.Clear()
		e.lastKey = 0
	case 'k':
		if e.cursorY > 0 {
			e.cursorY--
			e.adjustScroll()
			e.clampCursor()
		} else {
			e.message = "Top of file"
		}
		e.selection.Clear()
		e.lastKey = 0
	case 'l':
		line := e.buffer.Line(e.cursorY)
		if e.cursorX < len(line) {
			e.cursorX++
		}
		e.adjustScroll()
		e.selection.Clear()
		e.lastKey = 0
	default:
		// Clear lastKey if any other key is pressed
		e.lastKey = 0
	}
	
	switch ev.Key {
	case tcell.KeyEscape:
		e.selection.Clear()
		e.lastKey = 0
	case tcell.KeyLeft:
		if e.cursorX > 0 {
			e.cursorX--
		}
		e.adjustScroll()
		e.selection.Clear()
		e.lastKey = 0
	case tcell.KeyRight:
		line := e.buffer.Line(e.cursorY)
		if e.cursorX < len(line) {
			e.cursorX++
		}
		e.adjustScroll()
		e.selection.Clear()
		e.lastKey = 0
	case tcell.KeyUp:
		if e.cursorY > 0 {
			e.cursorY--
			e.adjustScroll()
			e.clampCursor()
		} else {
			e.message = "Top of file"
		}
		e.selection.Clear()
		e.lastKey = 0
	case tcell.KeyDown:
		if e.cursorY < e.buffer.LineCount()-1 {
			e.cursorY++
			e.adjustScroll()
			e.clampCursor()
		} else {
			e.message = "End of file"
		}
		e.selection.Clear()
		e.lastKey = 0
	}
}

func (e *Editor) jumpToStart() {
	e.cursorY = 0
	e.cursorX = 0
	e.offsetY = 0
	e.message = ""
}

func (e *Editor) jumpToEnd() {
	e.cursorY = e.buffer.LineCount() - 1
	e.cursorX = 0
	e.adjustScroll()
	e.message = ""
}

func (e *Editor) togglePreview() {
	e.preview.Toggle()
	if e.preview.IsEnabled() {
		e.preview.Update(e.buffer)
		e.message = "Preview enabled"
	} else {
		e.message = "Preview disabled"
	}
	e.renderCache.invalidate()
}

func (e *Editor) handlePreviewKeys(ev *terminal.Event) {
	switch ev.Rune {
	case 'p':
		// Toggle preview off
		e.togglePreview()
	case 'j':
		e.preview.Scroll(1)
	case 'k':
		e.preview.Scroll(-1)
	case 'q':
		e.quit = true
	}
	
	switch ev.Key {
	case tcell.KeyDown:
		e.preview.Scroll(1)
	case tcell.KeyUp:
		e.preview.Scroll(-1)
	case tcell.KeyCtrlC:
		e.quit = true
	}
	
	// Ensure render is triggered after preview key handling
	e.renderCache.invalidate()
}

func (e *Editor) copyCurrentLine() {
	line := e.buffer.Line(e.cursorY)
	err := clipboard.Copy(line)
	if err != nil {
		e.message = "Failed to copy to clipboard"
	} else {
		e.message = "Line copied to clipboard"
	}
}

func (e *Editor) pasteFromClipboard() {
	text, err := clipboard.Paste()
	if err != nil {
		e.message = "Failed to paste from clipboard"
		return
	}
	
	if text == "" {
		e.message = "Clipboard is empty"
		return
	}
	
	// Insert text at cursor position
	for _, r := range text {
		if r == '\n' {
			e.buffer.SplitLine(e.cursorY, e.cursorX)
			e.cursorY++
			e.cursorX = 0
		} else {
			e.buffer.InsertRune(e.cursorY, e.cursorX, r)
			e.cursorX++
		}
	}
	
	e.adjustScroll()
	e.message = "Pasted from clipboard"
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

func (e *Editor) performUndo() {
	if e.buffer.Undo() {
		e.message = "Undo"
		e.clampCursor()
		e.adjustScroll()
	} else {
		e.message = "Nothing to undo"
	}
}

func (e *Editor) performRedo() {
	if e.buffer.Redo() {
		e.message = "Redo"
		e.clampCursor()
		e.adjustScroll()
	} else {
		e.message = "Nothing to redo"
	}
}
