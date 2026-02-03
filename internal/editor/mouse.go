package editor

import (
	"github.com/Adelodunpeter25/vx/internal/terminal"
	"github.com/Adelodunpeter25/vx/internal/wrap"
	"github.com/gdamore/tcell/v2"
)

func (e *Editor) handleMouseEvent(ev *terminal.Event) {
	p := e.active()
	paneWidth := p.viewWidth
	paneHeight := p.viewHeight
	if paneWidth == 0 {
		paneWidth = e.width
	}
	if paneHeight == 0 {
		paneHeight = e.height - 1
	}
	if ev.MouseX < 0 || ev.MouseX >= paneWidth || ev.MouseY < 0 || ev.MouseY >= paneHeight {
		return
	}
	// Handle scroll wheel
	if ev.Button == tcell.WheelUp {
		if p.preview.IsEnabled() {
			p.preview.Scroll(-1)
		} else {
			// Scroll view up by one visual row
			if p.visualOffsetY > 0 {
				p.visualOffsetY--
				// Update offsetY to match
				gutterWidth := e.getGutterWidthFor(p)
				maxWidth := paneWidth - gutterWidth
				p.offsetY = e.findLineAtVisualRow(p.visualOffsetY, maxWidth)
			}
		}
		return
	}

	if ev.Button == tcell.WheelDown {
		if p.preview.IsEnabled() {
			p.preview.Scroll(1)
		} else {
			// Scroll view down by one visual row
			gutterWidth := e.getGutterWidthFor(p)
			maxWidth := paneWidth - gutterWidth
			contentHeight := paneHeight

			// Calculate total visual rows
			totalVisualRows := 0
			for i := 0; i < p.buffer.LineCount(); i++ {
				line := p.buffer.Line(i)
				totalVisualRows += wrap.VisualLineCount(line, maxWidth)
			}

			// Only scroll if there's more content below
			if p.visualOffsetY+contentHeight < totalVisualRows {
				p.visualOffsetY++
				// Update offsetY to match
				p.offsetY = e.findLineAtVisualRow(p.visualOffsetY, maxWidth)
			}
		}
		return
	}

	// Don't handle clicks in preview mode
	if p.preview.IsEnabled() {
		return
	}

	// Only handle left click for positioning and selection
	if ev.Button != tcell.Button1 && ev.Button != tcell.ButtonNone {
		return
	}

	// Detect button state change
	buttonPressed := ev.Button == tcell.Button1
	buttonReleased := ev.Button == tcell.ButtonNone && p.mouseDragging

	mouseX, mouseY := ev.MouseX, ev.MouseY

	// Ensure we're not clicking below the content area (status line)
	contentHeight := paneHeight
	if mouseY >= contentHeight {
		return
	}

	// Convert screen coordinates to buffer coordinates (accounting for wrapped lines and visual offset)
	gutterWidth := e.getGutterWidthFor(p)
	maxWidth := paneWidth - gutterWidth

	if mouseX < gutterWidth {
		return
	}

	bufferY, bufferX := e.bufferPosFromScreen(mouseX, mouseY, gutterWidth, maxWidth)

	// Check if button is pressed or released
	if buttonPressed {
		// Button is pressed/held
		if !p.mouseDragging && !p.selection.IsActive() {
			// First press - record position but don't start selection yet
			p.mouseDownX = mouseX
			p.mouseDownY = mouseY
			p.mouseDragging = true
		}

		// Check if mouse moved enough to start selection
		if p.mouseDragging && !p.selection.IsActive() && (abs(mouseX-p.mouseDownX) > 1 || abs(mouseY-p.mouseDownY) > 0) {
			// Mouse moved - start selection from original down position
			startBufferY, startBufferX := e.bufferPosFromScreen(p.mouseDownX, p.mouseDownY, gutterWidth, maxWidth)
			p.selection.Start(startBufferY, startBufferX)
		}

		// Update selection if active
		if p.selection.IsActive() {
			p.selection.Update(bufferY, bufferX)
		}

		// Auto-scroll if dragging near edges
		if p.selection.IsActive() {
			contentHeight := paneHeight
			if mouseY < 2 && p.visualOffsetY > 0 {
				p.visualOffsetY--
				p.offsetY = e.findLineAtVisualRow(p.visualOffsetY, maxWidth)
			} else if mouseY > contentHeight-3 {
				totalVisualRows := 0
				for i := 0; i < p.buffer.LineCount(); i++ {
					line := p.buffer.Line(i)
					totalVisualRows += wrap.VisualLineCount(line, maxWidth)
				}

				if p.visualOffsetY+contentHeight < totalVisualRows {
					p.visualOffsetY++
					p.offsetY = e.findLineAtVisualRow(p.visualOffsetY, maxWidth)
				}
			}
		}

		p.cursorY = bufferY
		p.cursorX = bufferX
		e.clampCursor()
	} else if buttonReleased {
		// Button released
		if !p.selection.IsActive() {
			// Was just a click, not a drag - move cursor
			p.cursorY = bufferY
			p.cursorX = bufferX
			e.clampCursor()
		}
		// Reset drag state
		p.mouseDragging = false
	}

	// Don't call adjustScroll here - let cursor stay where it is
	// Only adjust scroll when cursor moves via keyboard
}

func (e *Editor) bufferPosFromScreen(mouseX, mouseY, gutterWidth, maxWidth int) (bufferY, bufferX int) {
	p := e.active()
	clickedVisualRow := p.visualOffsetY + mouseY
	currentVisualRow := 0
	bufferY = 0
	bufferX = mouseX - gutterWidth

	for bufferY < p.buffer.LineCount() {
		line := p.buffer.Line(bufferY)
		segments := wrap.WrapLine(line, bufferY, maxWidth)

		for _, seg := range segments {
			if currentVisualRow == clickedVisualRow {
				bufferX = seg.StartCol + (mouseX - gutterWidth)
				segLen := len([]rune(seg.Text))
				if bufferX > seg.StartCol+segLen {
					bufferX = seg.StartCol + segLen
				}
				return bufferY, bufferX
			}
			currentVisualRow++
		}
		bufferY++
	}

	if bufferY >= p.buffer.LineCount() {
		bufferY = p.buffer.LineCount() - 1
	}
	return bufferY, bufferX
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
