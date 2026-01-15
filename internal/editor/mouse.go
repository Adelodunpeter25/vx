package editor

import (
	"github.com/Adelodunpeter25/vx/internal/terminal"
	"github.com/Adelodunpeter25/vx/internal/wrap"
	"github.com/gdamore/tcell/v2"
)

func (e *Editor) handleMouseEvent(ev *terminal.Event) {
	// Handle scroll wheel
	if ev.Button == tcell.WheelUp {
		if e.preview.IsEnabled() {
			e.preview.Scroll(-1)
		} else {
			// Scroll view up (decrease offsetY)
			if e.offsetY > 0 {
				e.offsetY--
			}
		}
		return
	}
	
	if ev.Button == tcell.WheelDown {
		if e.preview.IsEnabled() {
			e.preview.Scroll(1)
		} else {
			// Scroll view down (increase offsetY)
			maxOffset := e.buffer.LineCount() - (e.height - 1)
			if maxOffset < 0 {
				maxOffset = 0
			}
			if e.offsetY < maxOffset {
				e.offsetY++
			}
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
	
	// Ensure we're not clicking below the content area (status line)
	contentHeight := e.height - 1
	if mouseY >= contentHeight {
		return
	}
	
	// Convert screen coordinates to buffer coordinates (accounting for wrapped lines)
	gutterWidth := e.getGutterWidth()
	maxWidth := e.width - gutterWidth
	
	if mouseX < gutterWidth {
		return
	}
	
	// Find which buffer line and column the click corresponds to
	screenRow := 0
	bufferY := e.offsetY
	bufferX := mouseX - gutterWidth
	
	for bufferY < e.buffer.LineCount() {
		line := e.buffer.Line(bufferY)
		lineVisualRows := wrap.VisualLineCount(line, maxWidth)
		
		if screenRow+lineVisualRows > mouseY {
			// Click is on this buffer line
			// Calculate which wrapped segment
			segmentIndex := mouseY - screenRow
			bufferX = segmentIndex*maxWidth + (mouseX - gutterWidth)
			break
		}
		
		screenRow += lineVisualRows
		bufferY++
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
