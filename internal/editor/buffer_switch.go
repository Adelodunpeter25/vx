package editor

import "github.com/Adelodunpeter25/vx/internal/buffer"

// saveCurrentBufferState saves cursor and scroll position to buffer manager
func (e *Editor) saveCurrentBufferState() {
	p := e.active()
	p.bufferMgr.SaveState(p.cursorX, p.cursorY, p.offsetX, p.offsetY)
}

// switchToBuffer switches to a different buffer and restores its state
func (e *Editor) switchToBuffer() {
	p := e.active()
	current := p.bufferMgr.Current()
	if current == nil {
		return
	}

	p.buffer = current.Buffer
	p.syntax = current.Syntax
	p.cursorX, p.cursorY, p.offsetX, p.offsetY = p.bufferMgr.RestoreState()

	// Clamp cursor to valid position
	e.clampCursor()
	e.adjustScroll()
	p.renderCache.invalidate()
}

// addBuffer adds a new buffer and switches to it
func (e *Editor) addBuffer(buf *buffer.Buffer, filename string) {
	e.saveCurrentBufferState()
	e.active().bufferMgr.Add(buf, filename)
	e.switchToBuffer()
}

// nextBuffer switches to the next buffer
func (e *Editor) nextBuffer() {
	if e.active().bufferMgr.Count() <= 1 {
		return
	}
	e.saveCurrentBufferState()
	e.active().bufferMgr.Next()
	e.switchToBuffer()
}

// previousBuffer switches to the previous buffer
func (e *Editor) previousBuffer() {
	if e.active().bufferMgr.Count() <= 1 {
		return
	}
	e.saveCurrentBufferState()
	e.active().bufferMgr.Previous()
	e.switchToBuffer()
}

// deleteCurrentBuffer closes the current buffer with save prompt if modified
func (e *Editor) deleteCurrentBuffer() {
	p := e.active()
	if p.bufferMgr.Count() == 1 {
		p.msgManager.SetTransient("Cannot close last buffer")
		return
	}

	if p.buffer.IsModified() {
		// Enter a special prompt mode
		p.mode = ModeBufferPrompt
		p.msgManager.SetPersistent("Save changes? [y/n]")
		p.renderCache.invalidate()
		return
	}

	// Not modified, just delete
	p.bufferMgr.Delete()
	e.switchToBuffer()
	p.msgManager.SetTransient("Buffer closed")
}
