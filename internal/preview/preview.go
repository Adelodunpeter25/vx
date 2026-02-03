package preview

import (
	"github.com/Adelodunpeter25/vx/internal/buffer"
	"github.com/Adelodunpeter25/vx/internal/markdown"
	"github.com/Adelodunpeter25/vx/internal/terminal"
	"github.com/gdamore/tcell/v2"
)

// Preview manages the markdown preview pane
type Preview struct {
	enabled  bool
	elements []markdown.Element
	offsetY  int
}

func New() *Preview {
	return &Preview{
		enabled: false,
	}
}

// Toggle enables/disables preview
func (p *Preview) Toggle() {
	p.enabled = !p.enabled
}

// IsEnabled returns preview state
func (p *Preview) IsEnabled() bool {
	return p.enabled
}

// Update re-parses the buffer content
func (p *Preview) Update(buf *buffer.Buffer) {
	// Build full text from buffer
	lines := make([]string, buf.LineCount())
	for i := 0; i < buf.LineCount(); i++ {
		lines[i] = buf.Line(i)
	}

	text := ""
	for i, line := range lines {
		if i > 0 {
			text += "\n"
		}
		text += line
	}

	p.elements = markdown.Parse(text)
}

// Render draws the preview pane
func (p *Preview) Render(term *terminal.Terminal, startX, startY, height, width int) {
	if !p.enabled {
		return
	}

	y := startY
	for i := p.offsetY; i < len(p.elements) && y < startY+height; i++ {
		elem := p.elements[i]
		segments := markdown.RenderElement(elem)

		// Draw all segments on the same line
		x := startX
		for _, seg := range segments {
			for _, r := range seg.Text {
				if y < startY+height && x < startX+width {
					term.SetCell(x, y, r, seg.Style)
					x++
				}
			}
		}
		y++
	}

	// Fill remaining space
	for y < startY+height {
		for x := startX; x < startX+width; x++ {
			term.SetCell(x, y, ' ', tcell.StyleDefault)
		}
		y++
	}
}

// Scroll adjusts preview scroll position
func (p *Preview) Scroll(delta int) {
	p.offsetY += delta
	if p.offsetY < 0 {
		p.offsetY = 0
	}
	maxOffset := len(p.elements) - 1
	if maxOffset < 0 {
		maxOffset = 0
	}
	if p.offsetY > maxOffset {
		p.offsetY = maxOffset
	}
}
