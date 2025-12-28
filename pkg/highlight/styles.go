package highlight

import (
	"github.com/alecthomas/chroma/v2"
	"github.com/gdamore/tcell/v2"
)

// Keyword styles
func keywordStyle(t chroma.TokenType, base tcell.Style) tcell.Style {
	switch t {
	case chroma.KeywordConstant:
		return base.Foreground(tcell.NewRGBColor(255, 100, 100)).Bold(true)
	case chroma.KeywordType:
		return base.Foreground(tcell.NewRGBColor(100, 150, 255))
	case chroma.KeywordNamespace:
		return base.Foreground(tcell.NewRGBColor(150, 100, 255))
	default:
		return base.Foreground(tcell.ColorBlue).Bold(true)
	}
}

// String styles
func stringStyle(t chroma.TokenType, base tcell.Style) tcell.Style {
	switch t {
	case chroma.LiteralStringEscape:
		return base.Foreground(tcell.NewRGBColor(100, 255, 150)).Bold(true)
	case chroma.LiteralStringRegex:
		return base.Foreground(tcell.NewRGBColor(255, 200, 100))
	case chroma.LiteralStringSymbol:
		return base.Foreground(tcell.NewRGBColor(150, 255, 150))
	case chroma.LiteralStringDoc:
		return base.Foreground(tcell.NewRGBColor(100, 200, 100)).Italic(true)
	default:
		return base.Foreground(tcell.ColorGreen)
	}
}

// Number styles
func numberStyle(t chroma.TokenType, base tcell.Style) tcell.Style {
	switch t {
	case chroma.LiteralNumberHex:
		return base.Foreground(tcell.NewRGBColor(255, 200, 100))
	case chroma.LiteralNumberBin:
		return base.Foreground(tcell.NewRGBColor(200, 255, 100))
	case chroma.LiteralNumberOct:
		return base.Foreground(tcell.NewRGBColor(255, 255, 100))
	case chroma.LiteralNumberFloat:
		return base.Foreground(tcell.NewRGBColor(255, 220, 100))
	default:
		return base.Foreground(tcell.ColorYellow)
	}
}

// Comment styles
func commentStyle(base tcell.Style) tcell.Style {
	return base.Foreground(tcell.ColorGray).Italic(true)
}

// Name styles (functions, classes, variables)
func nameStyle(t chroma.TokenType, base tcell.Style) tcell.Style {
	switch t {
	case chroma.NameFunction, chroma.NameFunctionMagic:
		return base.Foreground(tcell.NewRGBColor(100, 200, 255))
	case chroma.NameClass:
		return base.Foreground(tcell.NewRGBColor(150, 255, 200)).Bold(true)
	case chroma.NameBuiltin, chroma.NameBuiltinPseudo:
		return base.Foreground(tcell.NewRGBColor(255, 100, 255))
	case chroma.NameDecorator:
		return base.Foreground(tcell.NewRGBColor(255, 150, 100))
	case chroma.NameException:
		return base.Foreground(tcell.NewRGBColor(255, 100, 100)).Bold(true)
	case chroma.NameConstant:
		return base.Foreground(tcell.NewRGBColor(255, 150, 150))
	case chroma.NameTag:
		return base.Foreground(tcell.NewRGBColor(100, 200, 255))
	case chroma.NameAttribute:
		return base.Foreground(tcell.NewRGBColor(200, 150, 255))
	case chroma.NameVariable, chroma.NameVariableInstance:
		return base.Foreground(tcell.NewRGBColor(200, 200, 255))
	case chroma.NameVariableClass, chroma.NameVariableGlobal:
		return base.Foreground(tcell.NewRGBColor(180, 180, 255)).Bold(true)
	case chroma.NameNamespace:
		return base.Foreground(tcell.NewRGBColor(150, 200, 255))
	default:
		return base.Foreground(tcell.NewRGBColor(200, 200, 200))
	}
}

// Operator styles
func operatorStyle(t chroma.TokenType, base tcell.Style) tcell.Style {
	switch t {
	case chroma.OperatorWord:
		return base.Foreground(tcell.NewRGBColor(150, 150, 255)).Bold(true)
	default:
		return base.Foreground(tcell.ColorWhite)
	}
}

// Literal styles
func literalStyle(t chroma.TokenType, base tcell.Style) tcell.Style {
	switch t {
	case chroma.LiteralDate:
		return base.Foreground(tcell.NewRGBColor(255, 200, 150))
	default:
		return base.Foreground(tcell.NewRGBColor(200, 255, 200))
	}
}
