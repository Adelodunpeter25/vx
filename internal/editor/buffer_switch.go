package editor

import "github.com/Adelodunpeter25/vx/internal/buffer"

// saveCurrentBufferState saves cursor and scroll position to buffer manager
func (e *Editor) saveCurrentBufferState() {
	e.bufferMgr.SaveState(e.cursorX, e.cursorY, e.offsetX, e.offsetY)
}

// switchToBuffer switches to a different buffer and restores its state
func (e *Editor) switchToBuffer() {
	current := e.bufferMgr.Current()
	if current == nil {
		return
	}

	e.buffer = current.Buffer
	e.syntax = current.Syntax
	e.cursorX, e.cursorY, e.offsetX, e.offsetY = e.bufferMgr.RestoreState()

	// Clamp cursor to valid position
	e.clampCursor()
	e.adjustScroll()
	e.renderCache.invalidate()
}

// addBuffer adds a new buffer and switches to it
func (e *Editor) addBuffer(buf *buffer.Buffer, filename string) {
	e.saveCurrentBufferState()
	e.bufferMgr.Add(buf, filename)
	e.switchToBuffer()
}

// nextBuffer switches to the next buffer
func (e *Editor) nextBuffer() {
	if e.bufferMgr.Count() <= 1 {
		return
	}
	e.saveCurrentBufferState()
	e.bufferMgr.Next()
	e.switchToBuffer()
}

// previousBuffer switches to the previous buffer
func (e *Editor) previousBuffer() {
	if e.bufferMgr.Count() <= 1 {
		return
	}
	e.saveCurrentBufferState()
	e.bufferMgr.Previous()
	e.switchToBuffer()
}

// deleteCurrentBuffer closes the current buffer with save prompt if modified
func (e *Editor) deleteCurrentBuffer() {
	if e.bufferMgr.Count() == 1 {
		e.message = "Cannot close last buffer"
		return
	}

	if e.buffer.IsModified() {
		// Enter a special prompt mode
		e.mode = ModeBufferPrompt
		e.message = "Save changes? [y/n]"
		e.renderCache.invalidate()
		return
	}

	// Not modified, just delete
	e.bufferMgr.Delete()
	e.switchToBuffer()
	e.message = "Buffer closed"
}
