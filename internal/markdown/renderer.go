package markdown

import (
	"github.com/gdamore/tcell/v2"
)

type RenderedSegment struct {
	Text  string
	Style tcell.Style
}

type VisualLine struct {
	Segments []RenderedSegment
}

func ElementVisualLines(elem Element, maxWidth int) []VisualLine {
	switch elem.Type {
	case TypeEmptyLine:
		return []VisualLine{{Segments: []RenderedSegment{{Text: "", Style: tcell.StyleDefault}}}}

	case TypeHorizontalRule:
		line := ""
		for i := 0; i < maxWidth; i++ {
			line += "─"
		}
		style := tcell.StyleDefault.Foreground(tcell.ColorGray)
		return []VisualLine{{Segments: []RenderedSegment{{Text: line, Style: style}}}}

	case TypeHeader:
		style := tcell.StyleDefault.Bold(true)
		switch elem.Level {
		case 1:
			style = style.Foreground(tcell.NewRGBColor(100, 200, 255))
		case 2:
			style = style.Foreground(tcell.NewRGBColor(150, 220, 255))
		default:
			style = style.Foreground(tcell.NewRGBColor(200, 230, 255))
		}
		base := style
		lines := renderSegmentsWithWrap(elem.Segments, base, maxWidth)
		if elem.Level == 1 {
			underline := tcell.StyleDefault.Foreground(tcell.NewRGBColor(60, 60, 60))
			ulText := ""
			for i := 0; i < maxWidth; i++ {
				ulText += "═"
			}
			lines = append(lines, VisualLine{Segments: []RenderedSegment{{Text: ulText, Style: underline}}})
		}
		return lines

	case TypeCodeBlock:
		style := tcell.StyleDefault.
			Foreground(tcell.NewRGBColor(255, 200, 100)).
			Background(tcell.NewRGBColor(40, 40, 40))
		var lines []VisualLine
		codeText := elem.Content
		if codeText == "" {
			return lines
		}
		runes := []rune(codeText)
		for len(runes) > 0 {
			chunk := runes
			if len(chunk) > maxWidth {
				chunk = chunk[:maxWidth]
			}
			lines = append(lines, VisualLine{
				Segments: []RenderedSegment{{Text: string(chunk), Style: style}},
			})
			runes = runes[len(chunk):]
		}
		return lines

	case TypeBlockquote:
		baseStyle := tcell.StyleDefault.Foreground(tcell.ColorGray).Italic(true)
		prefix := RenderedSegment{Text: "│ ", Style: baseStyle}
		return renderSegmentsWithWrapPrefixed(elem.Segments, baseStyle, maxWidth, prefix)

	case TypeList:
		base := tcell.StyleDefault
		prefix := RenderedSegment{Text: "• ", Style: base.Foreground(tcell.NewRGBColor(255, 200, 100))}
		return renderSegmentsWithWrapPrefixed(elem.Segments, base, maxWidth, prefix)

	case TypeOrderedList:
		base := tcell.StyleDefault
		numText := ""
		if elem.Number > 0 {
			numText = itoa(elem.Number) + ". "
		}
		prefix := RenderedSegment{Text: numText, Style: base.Foreground(tcell.NewRGBColor(255, 200, 100))}
		return renderSegmentsWithWrapPrefixed(elem.Segments, base, maxWidth, prefix)

	case TypeTaskList:
		base := tcell.StyleDefault
		checkMark := "[ ] "
		if elem.Checked != nil && *elem.Checked {
			checkMark = "[✓] "
		}
		prefix := RenderedSegment{Text: checkMark, Style: base}
		return renderSegmentsWithWrapPrefixed(elem.Segments, base, maxWidth, prefix)

	default:
		base := tcell.StyleDefault
		return renderSegmentsWithWrap(elem.Segments, base, maxWidth)
	}
}

func renderSegmentsWithWrap(segs []Segment, baseStyle tcell.Style, maxWidth int) []VisualLine {
	return renderSegmentsWithWrapPrefixed(segs, baseStyle, maxWidth, RenderedSegment{})
}

func renderSegmentsWithWrapPrefixed(segs []Segment, baseStyle tcell.Style, maxWidth int, prefix RenderedSegment) []VisualLine {
	var lines []VisualLine
	var current []RenderedSegment
	currentWidth := runeWidth(prefix.Text)

	if currentWidth > 0 {
		current = append(current, prefix)
	}

	needsPrefix := func() {
		if len(current) > 0 {
			lines = append(lines, VisualLine{Segments: current})
		}
		current = nil
		currentWidth = 0
		if prefix.Text != "" {
			current = append(current, prefix)
			currentWidth = runeWidth(prefix.Text)
		}
	}

	for _, seg := range segs {
		rs := segToRendered(seg, baseStyle)
		segRunes := []rune(rs.Text)

		if currentWidth+len(segRunes) <= maxWidth {
			current = append(current, rs)
			currentWidth += len(segRunes)
			continue
		}

		spaceLeft := maxWidth - currentWidth
		if spaceLeft > 0 {
			current = append(current, RenderedSegment{
				Text:  string(segRunes[:spaceLeft]),
				Style: rs.Style,
			})
			segRunes = segRunes[spaceLeft:]
		}

		needsPrefix()

		for len(segRunes) > 0 {
			chunk := segRunes
			if len(chunk) > maxWidth {
				chunk = chunk[:maxWidth]
			}
			current = append(current, RenderedSegment{Text: string(chunk), Style: rs.Style})
			segRunes = segRunes[len(chunk):]
			if len(segRunes) > 0 {
				needsPrefix()
			}
		}
		currentWidth = runeWidthPrefix(current)
	}

	if len(current) > 0 {
		lines = append(lines, VisualLine{Segments: current})
	}

	if len(lines) == 0 {
		lines = append(lines, VisualLine{Segments: []RenderedSegment{{Text: "", Style: baseStyle}}})
	}

	return lines
}

func segToRendered(seg Segment, baseStyle tcell.Style) RenderedSegment {
	style := baseStyle
	if seg.Bold {
		style = style.Bold(true)
	}
	if seg.Italic {
		style = style.Italic(true)
	}
	if seg.Code {
		style = tcell.StyleDefault.
			Foreground(tcell.NewRGBColor(255, 200, 100)).
			Background(tcell.NewRGBColor(40, 40, 40))
	}
	if seg.Strikethrough {
		style = style.Foreground(tcell.ColorGray)
	}
	if seg.Link {
		style = style.Foreground(tcell.NewRGBColor(80, 180, 255)).Underline(true)
	}
	return RenderedSegment{Text: seg.Text, Style: style}
}

func runeWidth(text string) int {
	return len([]rune(text))
}

func runeWidthPrefix(lines []RenderedSegment) int {
	n := 0
	for _, s := range lines {
		n += runeWidth(s.Text)
	}
	return n
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	s := ""
	for n > 0 {
		s = string(byte('0'+n%10)) + s
		n /= 10
	}
	return s
}
