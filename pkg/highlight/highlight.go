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

func tokenToStyle(tokenType chroma.TokenType) tcell.Style {
	base := tcell.StyleDefault

	switch {
	case tokenType == chroma.Keyword:
		return base.Foreground(tcell.ColorBlue).Bold(true)
	case tokenType == chroma.KeywordNamespace || tokenType == chroma.KeywordType:
		return base.Foreground(tcell.ColorBlue)
	case tokenType == chroma.String || tokenType == chroma.LiteralString:
		return base.Foreground(tcell.ColorGreen)
	case tokenType == chroma.LiteralStringDouble || tokenType == chroma.LiteralStringSingle:
		return base.Foreground(tcell.ColorGreen)
	case tokenType == chroma.Comment || tokenType == chroma.CommentSingle || tokenType == chroma.CommentMultiline:
		return base.Foreground(tcell.ColorGray)
	case tokenType == chroma.Number || tokenType == chroma.LiteralNumber:
		return base.Foreground(tcell.ColorYellow)
	case tokenType == chroma.LiteralNumberInteger || tokenType == chroma.LiteralNumberFloat:
		return base.Foreground(tcell.ColorYellow)
	case tokenType == chroma.Operator || tokenType == chroma.Punctuation:
		return base.Foreground(tcell.ColorWhite)
	case tokenType == chroma.Name || tokenType == chroma.NameFunction:
		return base.Foreground(tcell.NewRGBColor(0, 255, 255))
	case tokenType == chroma.NameTag || tokenType == chroma.NameAttribute:
		return base.Foreground(tcell.NewRGBColor(0, 200, 255))
	case strings.Contains(tokenType.String(), "Builtin"):
		return base.Foreground(tcell.NewRGBColor(255, 0, 255))
	case tokenType == chroma.LiteralStringSymbol:
		return base.Foreground(tcell.NewRGBColor(255, 165, 0))
	default:
		return base
	}
}
