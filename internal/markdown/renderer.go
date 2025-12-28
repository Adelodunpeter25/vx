package markdown

import (
	"github.com/gdamore/tcell/v2"
)

// RenderElement converts a markdown element to styled text for terminal
func RenderElement(elem Element) []RenderedSegment {
	segments := make([]RenderedSegment, 0)
	base := tcell.StyleDefault
	
	switch elem.Type {
	case TypeHeader:
		// Headers: bold and colored based on level
		style := base.Bold(true)
		switch elem.Level {
		case 1:
			style = style.Foreground(tcell.NewRGBColor(100, 200, 255))
		case 2:
			style = style.Foreground(tcell.NewRGBColor(150, 220, 255))
		default:
			style = style.Foreground(tcell.NewRGBColor(200, 230, 255))
		}
		
		// Render inline segments with header style as base
		for _, seg := range elem.Segments {
			segStyle := style
			if seg.Bold {
				segStyle = segStyle.Bold(true)
			}
			if seg.Italic {
				segStyle = segStyle.Italic(true)
			}
			if seg.Code {
				segStyle = base.Foreground(tcell.NewRGBColor(255, 200, 100)).
					Background(tcell.NewRGBColor(40, 40, 40))
			}
			segments = append(segments, RenderedSegment{Text: seg.Text, Style: segStyle})
		}
		
	case TypeCodeBlock:
		style := base.Foreground(tcell.NewRGBColor(255, 200, 100)).
			Background(tcell.NewRGBColor(40, 40, 40))
		segments = append(segments, RenderedSegment{Text: elem.Content, Style: style})
		
	case TypeList:
		segments = append(segments, RenderedSegment{Text: "• ", Style: base})
		for _, seg := range elem.Segments {
			segments = append(segments, renderSegment(seg, base))
		}
		
	case TypeBlockquote:
		style := base.Foreground(tcell.ColorGray).Italic(true)
		segments = append(segments, RenderedSegment{Text: "│ ", Style: style})
		for _, seg := range elem.Segments {
			segments = append(segments, renderSegment(seg, style))
		}
		
	default: // TypeText
		for _, seg := range elem.Segments {
			segments = append(segments, renderSegment(seg, base))
		}
	}
	
	return segments
}

func renderSegment(seg Segment, baseStyle tcell.Style) RenderedSegment {
	style := baseStyle
	
	if seg.Bold {
		style = style.Bold(true)
	}
	if seg.Italic {
		style = style.Italic(true)
	}
	if seg.Code {
		style = tcell.StyleDefault.Foreground(tcell.NewRGBColor(255, 200, 100)).
			Background(tcell.NewRGBColor(40, 40, 40))
	}
	
	return RenderedSegment{Text: seg.Text, Style: style}
}

type RenderedSegment struct {
	Text  string
	Style tcell.Style
}
