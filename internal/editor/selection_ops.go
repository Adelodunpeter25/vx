package editor

import "github.com/Adelodunpeter25/vx/internal/clipboard"

// copySelection copies the selected text to clipboard
func (e *Editor) copySelection() {
	p := e.active()
	text := p.selection.GetSelectedText(p.buffer)
	if text == "" {
		return
	}

	err := clipboard.Copy(text)
	if err != nil {
		p.msgManager.SetError("Failed to copy selection")
	} else {
		p.msgManager.SetTransient("Selection copied")
	}
	p.selection.Clear()
}

// cutSelection copies and deletes the selected text
func (e *Editor) cutSelection() {
	p := e.active()
	text := p.selection.GetSelectedText(p.buffer)
	if text == "" {
		return
	}

	// Copy to clipboard
	err := clipboard.Copy(text)
	if err != nil {
		p.msgManager.SetError("Failed to cut selection")
		return
	}

	// Get selection range for cursor positioning
	startLine, startCol, _, _, ok := p.selection.GetRange()
	if !ok {
		return
	}

	// Delete the selected text
	p.selection.DeleteSelectedText(p.buffer)

	// Position cursor at start of selection
	p.cursorY = startLine
	p.cursorX = startCol
	e.clampCursor()

	p.msgManager.SetTransient("Selection cut")
	p.selection.Clear()
}

// deleteSelection deletes the selected text without touching the clipboard
func (e *Editor) deleteSelection() {
	p := e.active()
	startLine, startCol, _, _, ok := p.selection.GetRange()
	if !ok {
		return
	}
	p.selection.DeleteSelectedText(p.buffer)
	p.cursorY = startLine
	p.cursorX = startCol
	e.clampCursor()
	e.adjustScroll()
	p.selection.Clear()
	p.msgManager.SetTransient("Selection deleted")
}
