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
			// Calculate total visual rows and max scroll position
			gutterWidth := e.getGutterWidth()
			maxWidth := e.width - gutterWidth
			contentHeight := e.height - 1
			
			totalVisualRows := 0
			for i := 0; i < e.buffer.LineCount(); i++ {
				line := e.buffer.Line(i)
				totalVisualRows += wrap.VisualLineCount(line, maxWidth)
			}
			
			// Calculate current visual position of offsetY
			currentVisualOffset := 0
			for i := 0; i < e.offsetY && i < e.buffer.LineCount(); i++ {
				line := e.buffer.Line(i)
				currentVisualOffset += wrap.VisualLineCount(line, maxWidth)
			}
			
			// Only scroll if there's more content below
			if currentVisualOffset + contentHeight < totalVisualRows {
				e.offsetY++
			}
		}
		return
	}
	
	// Don't handle clicks in preview mode
	if e.preview.IsEnabled() {
		return
	}
	
	// Only handle left click for positioning and selection
	if ev.Button != tcell.Button1 && ev.Button != tcell.ButtonNone {
		return
	}
	
	// Detect button state change
	buttonPressed := ev.Button == tcell.Button1
	buttonReleased := ev.Button == tcell.ButtonNone && e.mouseDragging
	
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
	
	// Check if button is pressed or released
	if buttonPressed {
		// Button is pressed/held
		if !e.mouseDragging && !e.selection.IsActive() {
			// First press - record position but don't start selection yet
			e.mouseDownX = mouseX
			e.mouseDownY = mouseY
			e.mouseDragging = true
		}
		
		// Check if mouse moved enough to start selection
		if e.mouseDragging && !e.selection.IsActive() && (abs(mouseX-e.mouseDownX) > 1 || abs(mouseY-e.mouseDownY) > 0) {
			// Mouse moved - start selection from original down position
			// Convert mouseDown position to buffer coordinates
			screenRow := 0
			startBufferY := e.offsetY
			startBufferX := e.mouseDownX - gutterWidth
			
			for startBufferY < e.buffer.LineCount() {
				line := e.buffer.Line(startBufferY)
				lineVisualRows := wrap.VisualLineCount(line, maxWidth)
				
				if screenRow+lineVisualRows > e.mouseDownY {
					segmentIndex := e.mouseDownY - screenRow
					startBufferX = segmentIndex*maxWidth + (e.mouseDownX - gutterWidth)
					break
				}
				
				screenRow += lineVisualRows
				startBufferY++
			}
			
			if startBufferY >= e.buffer.LineCount() {
				startBufferY = e.buffer.LineCount() - 1
			}
			
			e.selection.Start(startBufferY, startBufferX)
		}
		
		// Update selection if active
		if e.selection.IsActive() {
			e.selection.Update(bufferY, bufferX)
		}
		
		// Auto-scroll if dragging near edges
		if e.selection.IsActive() {
			contentHeight := e.height - 1
			if mouseY < 2 && e.offsetY > 0 {
				e.offsetY--
			} else if mouseY > contentHeight - 3 {
				totalVisualRows := 0
				for i := 0; i < e.buffer.LineCount(); i++ {
					line := e.buffer.Line(i)
					totalVisualRows += wrap.VisualLineCount(line, maxWidth)
				}
				
				currentVisualOffset := 0
				for i := 0; i < e.offsetY && i < e.buffer.LineCount(); i++ {
					line := e.buffer.Line(i)
					currentVisualOffset += wrap.VisualLineCount(line, maxWidth)
				}
				
				if currentVisualOffset + contentHeight < totalVisualRows {
					e.offsetY++
				}
			}
		}
		
		e.cursorY = bufferY
		e.cursorX = bufferX
		e.clampCursor()
	} else if buttonReleased {
		// Button released
		if !e.selection.IsActive() {
			// Was just a click, not a drag - move cursor
			e.cursorY = bufferY
			e.cursorX = bufferX
			e.clampCursor()
		}
		// Reset drag state
		e.mouseDragging = false
	}
	
	// Don't call adjustScroll here - let cursor stay where it is
	// Only adjust scroll when cursor moves via keyboard
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
