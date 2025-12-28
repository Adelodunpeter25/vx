package markdown

import (
	"strings"

	"github.com/gdamore/tcell/v2"
)

// RenderElement converts a markdown element to styled text
func RenderElement(elem Element) (string, tcell.Style) {
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
		prefix := strings.Repeat("#", elem.Level) + " "
		return prefix + elem.Content, style
		
	case TypeBold:
		return elem.Content, base.Bold(true)
		
	case TypeItalic:
		return elem.Content, base.Italic(true)
		
	case TypeCode, TypeCodeBlock:
		style := base.Foreground(tcell.NewRGBColor(255, 200, 100)).
			Background(tcell.NewRGBColor(40, 40, 40))
		return elem.Content, style
		
	case TypeList:
		return "• " + elem.Content, base
		
	case TypeBlockquote:
		style := base.Foreground(tcell.ColorGray).Italic(true)
		return "│ " + elem.Content, style
		
	case TypeLink:
		style := base.Foreground(tcell.NewRGBColor(100, 150, 255)).Underline(true)
		return elem.Content, style
		
	default:
		return elem.Content, base
	}
}
