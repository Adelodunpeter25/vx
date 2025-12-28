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
