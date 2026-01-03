package editor

import (
	"github.com/Adelodunpeter25/vx/internal/terminal"
	"github.com/gdamore/tcell/v2"
)

func (e *Editor) handleMouseEvent(ev *terminal.Event) {
	// Handle scroll wheel
	if ev.Button == tcell.WheelUp {
		if e.preview.IsEnabled() {
			e.preview.Scroll(-1)
		} else if e.cursorY > 0 {
			e.cursorY--
			e.adjustScroll()
			e.clampCursor()
		}
		return
	}
	
	if ev.Button == tcell.WheelDown {
		if e.preview.IsEnabled() {
			e.preview.Scroll(1)
		} else if e.cursorY < e.buffer.LineCount()-1 {
			e.cursorY++
			e.adjustScroll()
			e.clampCursor()
		}
		return
	}
	
	// Don't handle clicks in preview mode
	if e.preview.IsEnabled() {
		return
	}
	
	// Only handle left click for positioning
	if ev.Button != tcell.Button1 {
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
