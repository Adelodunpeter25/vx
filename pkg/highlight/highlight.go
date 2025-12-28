package highlight

import (
	"strings"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/gdamore/tcell/v2"
)

type Highlighter struct {
	lexer chroma.Lexer
}

func New(filename string) *Highlighter {
	lexer := lexers.Match(filename)
	if lexer == nil {
		lexer = lexers.Fallback
	}
	lexer = chroma.Coalesce(lexer)
	return &Highlighter{lexer: lexer}
}

func (h *Highlighter) HighlightText(text string) [][]StyledRune {
	if h.lexer == nil {
		return plainTextLines(text)
	}

	iterator, err := h.lexer.Tokenise(nil, text)
	if err != nil {
		return plainTextLines(text)
	}

	var lines [][]StyledRune
	var currentLine []StyledRune

	for _, token := range iterator.Tokens() {
		style := tokenToStyle(token.Type)
		for _, r := range token.Value {
			if r == '\n' {
				lines = append(lines, currentLine)
				currentLine = []StyledRune{}
			} else {
				currentLine = append(currentLine, StyledRune{
					Rune:  r,
					Style: style,
				})
			}
		}
	}
	
	if len(currentLine) > 0 {
		lines = append(lines, currentLine)
	}
	
	return lines
}

func (h *Highlighter) HighlightLine(line string) []StyledRune {
	if h.lexer == nil {
		return plainText(line)
	}

	iterator, err := h.lexer.Tokenise(nil, line+"\n")
	if err != nil {
		return plainText(line)
	}

	var result []StyledRune
	for _, token := range iterator.Tokens() {
		style := tokenToStyle(token.Type)
		for _, r := range token.Value {
			if r != '\n' {
				result = append(result, StyledRune{
					Rune:  r,
					Style: style,
				})
			}
		}
	}
	return result
}

type StyledRune struct {
	Rune  rune
	Style tcell.Style
}

func plainText(line string) []StyledRune {
	result := make([]StyledRune, 0, len(line))
	for _, r := range line {
		result = append(result, StyledRune{
			Rune:  r,
			Style: tcell.StyleDefault,
		})
	}
	return result
}

func plainTextLines(text string) [][]StyledRune {
	lines := strings.Split(text, "\n")
	result := make([][]StyledRune, len(lines))
	for i, line := range lines {
		result[i] = plainText(line)
	}
	return result
}

func tokenToStyle(tokenType chroma.TokenType) tcell.Style {
	base := tcell.StyleDefault
	typeStr := tokenType.String()

	switch {
	case tokenType == chroma.Keyword:
		return base.Foreground(tcell.ColorBlue).Bold(true)
	case tokenType == chroma.KeywordNamespace || tokenType == chroma.KeywordType:
		return base.Foreground(tcell.ColorBlue)
	case tokenType == chroma.LiteralStringDouble || tokenType == chroma.LiteralStringSingle:
		return base.Foreground(tcell.ColorGreen)
	case strings.Contains(typeStr, "String"):
		return base.Foreground(tcell.ColorGreen)
	case tokenType == chroma.Comment || tokenType == chroma.CommentSingle || tokenType == chroma.CommentMultiline:
		return base.Foreground(tcell.ColorGray)
	case tokenType == chroma.LiteralNumberInteger || tokenType == chroma.LiteralNumberFloat:
		return base.Foreground(tcell.ColorYellow)
	case strings.Contains(typeStr, "Number"):
		return base.Foreground(tcell.ColorYellow)
	case tokenType == chroma.NameTag:
		return base.Foreground(tcell.NewRGBColor(100, 200, 255))
	case tokenType == chroma.NameAttribute:
		return base.Foreground(tcell.NewRGBColor(255, 200, 100))
	case tokenType == chroma.Operator || tokenType == chroma.Punctuation:
		return base.Foreground(tcell.ColorWhite)
	case tokenType == chroma.NameFunction:
		return base.Foreground(tcell.NewRGBColor(0, 255, 255))
	case strings.Contains(typeStr, "Builtin"):
		return base.Foreground(tcell.NewRGBColor(255, 0, 255))
	case tokenType == chroma.KeywordConstant:
		return base.Foreground(tcell.NewRGBColor(255, 100, 100))
	default:
		return base
	}
}
