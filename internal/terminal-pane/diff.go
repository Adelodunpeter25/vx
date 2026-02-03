package terminalpane

import (
	"image/color"

	uv "github.com/charmbracelet/ultraviolet"
	"github.com/gdamore/tcell/v2"
)

func encodeCell(cell *uv.Cell) (rune, uint64) {
	if cell == nil {
		return ' ', 0
	}
	style := uint64(0)
	if cell.Style.Fg != nil {
		style ^= colorToUint(cell.Style.Fg)
	}
	if cell.Style.Bg != nil {
		style ^= colorToUint(cell.Style.Bg) << 24
	}
	style ^= uint64(cell.Style.Attrs) << 48
	return cell.Rune, style
}

func colorToUint(c color.Color) uint64 {
	r, g, b, _ := c.RGBA()
	return (uint64(r>>8) << 16) | (uint64(g>>8) << 8) | uint64(b>>8)
}

func toTcellStyle(style uint64) tcell.Style {
	fg := int64((style >> 0) & 0xFFFFFF)
	bg := int64((style >> 24) & 0xFFFFFF)
	attrs := int((style >> 48) & 0xFFFF)
	st := tcell.StyleDefault.Foreground(tcell.NewRGBColor((fg>>16)&0xFF, (fg>>8)&0xFF, fg&0xFF))
	st = st.Background(tcell.NewRGBColor((bg>>16)&0xFF, (bg>>8)&0xFF, bg&0xFF))
	if attrs&int(uv.AttrBold) != 0 {
		st = st.Bold(true)
	}
	if attrs&int(uv.AttrItalic) != 0 {
		st = st.Italic(true)
	}
	if attrs&int(uv.AttrReverse) != 0 {
		st = st.Reverse(true)
	}
	if uv.UnderlineStyle(attrs) != uv.UnderlineStyleNone {
		st = st.Underline(true)
	}
	return st
}
