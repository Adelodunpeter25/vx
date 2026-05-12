package preview

import (
	"github.com/Adelodunpeter25/vx/internal/buffer"
	"github.com/Adelodunpeter25/vx/internal/markdown"
	"github.com/Adelodunpeter25/vx/internal/terminal"
	"github.com/gdamore/tcell/v2"
)

type Preview struct {
	enabled     bool
	elements    []markdown.Element
	visualLines []markdown.VisualLine
	offsetY     int
	modVersion  int
	lastWidth   int
	viewHeight  int
}

func New() *Preview {
	return &Preview{enabled: false}
}

func (p *Preview) Toggle() {
	p.enabled = !p.enabled
}

func (p *Preview) IsEnabled() bool {
	return p.enabled
}

func (p *Preview) Update(buf *buffer.Buffer, width int) {
	if width <= 0 {
		return
	}
	modVer := buf.ModVersion()
	if modVer == p.modVersion && width == p.lastWidth && len(p.visualLines) > 0 {
		return
	}

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
	p.lastWidth = width

	maxWidth := width
	if maxWidth < 1 {
		maxWidth = 1
	}

	p.visualLines = nil
	for _, elem := range p.elements {
		vls := markdown.ElementVisualLines(elem, maxWidth)
		p.visualLines = append(p.visualLines, vls...)
	}

	p.modVersion = modVer
	p.clampOffset()
}

func (p *Preview) TotalVisualLines() int {
	return len(p.visualLines)
}

func (p *Preview) Render(term *terminal.Terminal, startX, startY, height, width int) {
	if !p.enabled || width <= 0 || height <= 0 {
		return
	}

	p.viewHeight = height
	p.clampOffset()

	for y := 0; y < height; y++ {
		vi := p.offsetY + y
		for x := 0; x < width; x++ {
			term.SetCell(startX+x, startY+y, ' ', tcell.StyleDefault)
		}
		if vi < len(p.visualLines) {
			vline := p.visualLines[vi]
			dx := startX
			for _, seg := range vline.Segments {
				for _, r := range seg.Text {
					if dx < startX+width {
						term.SetCell(dx, startY+y, r, seg.Style)
						dx++
					}
				}
			}
		}
	}

}

func (p *Preview) Scroll(delta int) {
	p.offsetY += delta
	p.clampOffset()
}

func (p *Preview) ScrollPage(delta int, viewHeight int) {
	p.offsetY += delta * viewHeight
	p.clampOffset()
}

func (p *Preview) ScrollToStart() {
	p.offsetY = 0
}

func (p *Preview) ScrollToEnd() {
	p.offsetY = len(p.visualLines) - 1
	p.clampOffset()
}

func (p *Preview) Offset() int {
	return p.offsetY
}

func (p *Preview) clampOffset() {
	maxOffset := 0
	total := len(p.visualLines)
	if p.viewHeight > 0 && total > p.viewHeight {
		maxOffset = total - p.viewHeight
	} else if total > 0 {
		maxOffset = total - 1
	}
	if p.offsetY < 0 {
		p.offsetY = 0
	}
	if p.offsetY > maxOffset {
		p.offsetY = maxOffset
	}
}
