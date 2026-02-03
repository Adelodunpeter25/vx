package terminalpane

import (
	"github.com/Adelodunpeter25/vx/internal/terminal"
)

type Renderer struct {
	last []CellSnapshot
	cols int
	rows int
}

func NewRenderer(cols, rows int) *Renderer {
	return &Renderer{cols: cols, rows: rows, last: make([]CellSnapshot, cols*rows)}
}

func (r *Renderer) Resize(cols, rows int) {
	if cols == r.cols && rows == r.rows {
		return
	}
	r.cols = cols
	r.rows = rows
	r.last = make([]CellSnapshot, cols*rows)
}

func (r *Renderer) Render(term *terminal.Terminal, emu *Emulator, x, y, cols, rows int) {
	if term == nil || emu == nil || cols <= 0 || rows <= 0 {
		return
	}
	if cols != r.cols || rows != r.rows {
		r.Resize(cols, rows)
	}
	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			cell := emu.CellAt(col, row)
			ch, style := encodeCell(cell)
			idx := row*cols + col
			if r.last[idx].Ch == ch && r.last[idx].Style == style {
				continue
			}
			r.last[idx] = CellSnapshot{Ch: ch, Style: style}
			term.SetCell(x+col, y+row, ch, toTcellStyle(style))
		}
	}
}
