package editor

import (
	"unicode/utf8"

	"github.com/Adelodunpeter25/vx/internal/wrap"
)

func lineRuneCount(line string) int {
	return utf8.RuneCountInString(line)
}

func (e *Editor) clampCursor() {
	if p := e.active(); p != nil {
		e.clampCursorForPane(p)
	}
}

func (e *Editor) clampCursorForPane(p *Pane) {
	if p == nil {
		return
	}
	if p.cursorY < 0 {
		p.cursorY = 0
	}
	if p.cursorY >= p.buffer.LineCount() {
		p.cursorY = max(0, p.buffer.LineCount()-1)
	}
	line := p.buffer.Line(p.cursorY)
	maxX := lineRuneCount(line)
	if p.mode == ModeNormal && maxX > 0 {
		maxX--
	}
	if p.cursorX > maxX {
		p.cursorX = maxX
	}
	if p.cursorX < 0 {
		p.cursorX = 0
	}
}

func (e *Editor) getPaneContentHeight(p *Pane) int {
	if p != nil && p.viewHeight > 0 {
		return p.viewHeight
	}
	h := e.height - 1
	if h < 1 {
		h = 1
	}
	return h
}

func (e *Editor) adjustScroll() {
	p := e.active()
	if p != nil {
		e.adjustScrollForPane(p)
	}
}

func (e *Editor) adjustScrollForPane(p *Pane) {
	if p == nil {
		return
	}
	// Use per-pane view dimensions when available (set in renderPane), fallback to global.
	contentHeight := p.viewHeight
	if contentHeight <= 0 {
		contentHeight = e.height - 1
		if contentHeight < 1 {
			contentHeight = 1
		}
	}
	gutterWidth := e.getGutterWidthFor(p)
	maxWidth := p.viewWidth - gutterWidth
	if maxWidth <= 0 {
		// Fallback: account for file-browser and split before first render
		contentWidth := e.width
		if e.fileBrowser != nil && e.fileBrowser.Open {
			fbWidth := e.fileBrowser.Width
			if fbWidth < 10 {
				fbWidth = 10
			}
			if fbWidth > e.width-10 {
				fbWidth = e.width - 10
			}
			contentWidth = e.width - fbWidth - 1
		}
		if len(e.panes) > 1 {
			// Approximate per-pane width in split layout
			contentWidth = (contentWidth - 1) / len(e.panes)
			if contentWidth < 10 {
				contentWidth = 10
			}
		}
		maxWidth = contentWidth - gutterWidth
		if maxWidth < 1 {
			maxWidth = 1
		}
	}

	// Calculate visual line position of cursor
	cursorVisualLine := 0
	for lineNum := 0; lineNum < p.cursorY && lineNum < p.buffer.LineCount(); lineNum++ {
		line := p.buffer.Line(lineNum)
		cursorVisualLine += wrap.VisualLineCount(line, maxWidth)
	}

	// Find which wrapped segment contains the cursor
	currentLine := p.buffer.Line(p.cursorY)
	segments := wrap.WrapLine(currentLine, p.cursorY, maxWidth)
	for i, seg := range segments {
		segEndCol := seg.StartCol + len([]rune(seg.Text))
		if p.cursorX >= seg.StartCol && p.cursorX <= segEndCol {
			cursorVisualLine += i
			break
		}
	}

	// Adjust visual offset to keep cursor visible
	if cursorVisualLine < p.visualOffsetY {
		// Cursor above viewport - scroll up
		p.visualOffsetY = cursorVisualLine
	}
	if cursorVisualLine >= p.visualOffsetY+contentHeight {
		// Cursor below viewport - scroll down
		p.visualOffsetY = cursorVisualLine - contentHeight + 1
	}

	// Ensure visual offset within bounds
	if p.visualOffsetY < 0 {
		p.visualOffsetY = 0
	}
	// Clamp upper bound based on total visual rows
	totalVisualRows := 0
	for i := 0; i < p.buffer.LineCount(); i++ {
		totalVisualRows += wrap.VisualLineCount(p.buffer.Line(i), maxWidth)
	}
	maxOffset := totalVisualRows - contentHeight
	if maxOffset < 0 {
		maxOffset = 0
	}
	if p.visualOffsetY > maxOffset {
		p.visualOffsetY = maxOffset
	}

	// Convert visual offset to buffer line offset for rendering
	p.offsetY = e.findLineAtVisualRowForPane(p, p.visualOffsetY, maxWidth)

	// No horizontal scroll needed with wrapping
	p.offsetX = 0
}

// findLineAtVisualRow finds which buffer line contains the given visual row for active pane
func (e *Editor) findLineAtVisualRow(targetVisual, maxWidth int) int {
	p := e.active()
	return e.findLineAtVisualRowForPane(p, targetVisual, maxWidth)
}

func (e *Editor) findLineAtVisualRowForPane(p *Pane, targetVisual, maxWidth int) int {
	if p == nil {
		return 0
	}
	visualLine := 0
	for lineNum := 0; lineNum < p.buffer.LineCount(); lineNum++ {
		line := p.buffer.Line(lineNum)
		lineVisualCount := wrap.VisualLineCount(line, maxWidth)
		if visualLine+lineVisualCount > targetVisual {
			return lineNum
		}
		visualLine += lineVisualCount
	}
	return max(0, p.buffer.LineCount()-1)
}
