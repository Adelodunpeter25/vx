package editor

import "github.com/Adelodunpeter25/vx/internal/clipboard"

// copySelection copies the selected text to clipboard
func (e *Editor) copySelection() {
	text := e.selection.GetSelectedText(e.buffer)
	if text == "" {
		return
	}
	
	err := clipboard.Copy(text)
	if err != nil {
		e.message = "Failed to copy selection"
	} else {
		e.message = "Selection copied"
	}
	e.selection.Clear()
}

// cutSelection copies and deletes the selected text
func (e *Editor) cutSelection() {
	text := e.selection.GetSelectedText(e.buffer)
	if text == "" {
		return
	}
	
	// Copy to clipboard
	err := clipboard.Copy(text)
	if err != nil {
		e.message = "Failed to cut selection"
		return
	}
	
	// Get selection range for cursor positioning
	startLine, startCol, _, _, ok := e.selection.GetRange()
	if !ok {
		return
	}
	
	// Delete the selected text
	e.selection.DeleteSelectedText(e.buffer)
	
	// Position cursor at start of selection
	e.cursorY = startLine
	e.cursorX = startCol
	e.clampCursor()
	
	e.message = "Selection cut"
	e.selection.Clear()
}
