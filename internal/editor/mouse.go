package editor

import (
	"github.com/Adelodunpeter25/vx/internal/terminal"
	"github.com/gdamore/tcell/v2"
)

func (e *Editor) handleMouseEvent(ev *terminal.Event) {
	// Don't handle mouse in preview mode
	if e.preview.IsEnabled() {
		return
	}
	
	// Only handle left click for now
	if ev.Button != tcell.Button1 {
		return
	}
	
	// Only handle button press, not release or motion
	if ev.Button == tcell.ButtonNone {
		return
	}
	
	mouseX, mouseY := ev.MouseX, ev.MouseY
	
	// Convert screen coordinates to buffer coordinates
	bufferY := mouseY + e.offsetY
	bufferX := mouseX + e.offsetX
	
	// Ensure we're not clicking below the content area (status line)
	contentHeight := e.height - 1
	if mouseY >= contentHeight {
		return
	}
	
	// Ensure click is within buffer bounds
	if bufferY >= e.buffer.LineCount() {
		bufferY = e.buffer.LineCount() - 1
	}
	
	// Move cursor to clicked position
	e.cursorY = bufferY
	e.cursorX = bufferX
	
	// Clamp cursor to valid position
	e.clampCursor()
	e.adjustScroll()
}
