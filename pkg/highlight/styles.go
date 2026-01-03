package highlight

import (
	"github.com/alecthomas/chroma/v2"
	"github.com/gdamore/tcell/v2"
)

// Keyword styles - VS Code blue
func keywordStyle(t chroma.TokenType, base tcell.Style) tcell.Style {
	switch t {
	case chroma.KeywordConstant:
		return base.Foreground(tcell.NewRGBColor(86, 156, 214)).Bold(true) // Light blue
	case chroma.KeywordType:
		return base.Foreground(tcell.NewRGBColor(78, 201, 176)) // Teal
	case chroma.KeywordNamespace:
		return base.Foreground(tcell.NewRGBColor(197, 134, 192)) // Purple
	default:
		return base.Foreground(tcell.NewRGBColor(86, 156, 214)) // Light blue
	}
}

// String styles - VS Code orange/red
func stringStyle(t chroma.TokenType, base tcell.Style) tcell.Style {
	switch t {
	case chroma.LiteralStringEscape:
		return base.Foreground(tcell.NewRGBColor(214, 157, 133)) // Light orange
	case chroma.LiteralStringRegex:
		return base.Foreground(tcell.NewRGBColor(209, 105, 105)) // Red
	case chroma.LiteralStringSymbol:
		return base.Foreground(tcell.NewRGBColor(206, 145, 120)) // Orange
	case chroma.LiteralStringDoc:
		return base.Foreground(tcell.NewRGBColor(106, 153, 85)) // Green
	case chroma.LiteralStringBacktick:
		return base.Foreground(tcell.NewRGBColor(206, 145, 120)) // Orange
	case chroma.LiteralStringChar:
		return base.Foreground(tcell.NewRGBColor(206, 145, 120)) // Orange
	case chroma.LiteralStringInterpol:
		return base.Foreground(tcell.NewRGBColor(156, 220, 254)) // Sky blue (interpolation)
	case chroma.LiteralStringHeredoc:
		return base.Foreground(tcell.NewRGBColor(206, 145, 120)) // Orange
	default:
		return base.Foreground(tcell.NewRGBColor(206, 145, 120)) // Orange
	}
}

// Number styles - VS Code light green
func numberStyle(t chroma.TokenType, base tcell.Style) tcell.Style {
	return base.Foreground(tcell.NewRGBColor(181, 206, 168)) // Light green
}

// Comment styles - VS Code gray/green
func commentStyle(base tcell.Style) tcell.Style {
	return base.Foreground(tcell.NewRGBColor(106, 153, 85)).Italic(true) // Green italic
}

// Name styles (functions, classes, variables)
func nameStyle(t chroma.TokenType, base tcell.Style) tcell.Style {
	switch t {
	case chroma.NameFunction, chroma.NameFunctionMagic:
		return base.Foreground(tcell.NewRGBColor(220, 220, 170)) // Yellow
	case chroma.NameClass:
		return base.Foreground(tcell.NewRGBColor(78, 201, 176)) // Teal
	case chroma.NameBuiltin, chroma.NameBuiltinPseudo:
		return base.Foreground(tcell.NewRGBColor(86, 156, 214)) // Light blue
	case chroma.NameDecorator:
		return base.Foreground(tcell.NewRGBColor(220, 220, 170)) // Yellow
	case chroma.NameException:
		return base.Foreground(tcell.NewRGBColor(78, 201, 176)) // Teal
	case chroma.NameConstant:
		return base.Foreground(tcell.NewRGBColor(79, 193, 255)) // Bright blue
	case chroma.NameTag:
		return base.Foreground(tcell.NewRGBColor(86, 156, 214)) // Light blue
	case chroma.NameAttribute:
		return base.Foreground(tcell.NewRGBColor(156, 220, 254)) // Sky blue
	case chroma.NameVariable, chroma.NameVariableInstance, chroma.NameVariableMagic:
		return base.Foreground(tcell.NewRGBColor(156, 220, 254)) // Sky blue
	case chroma.NameVariableClass, chroma.NameVariableGlobal:
		return base.Foreground(tcell.NewRGBColor(79, 193, 255)) // Bright blue
	case chroma.NameNamespace:
		return base.Foreground(tcell.NewRGBColor(197, 134, 192)) // Purple
	case chroma.NameLabel:
		return base.Foreground(tcell.NewRGBColor(197, 134, 192)) // Purple
	case chroma.NameEntity:
		return base.Foreground(tcell.NewRGBColor(220, 220, 170)) // Yellow
	case chroma.NameProperty:
		return base.Foreground(tcell.NewRGBColor(156, 220, 254)) // Sky blue
	default:
		return base.Foreground(tcell.NewRGBColor(212, 212, 212)) // Light gray
	}
}

// Operator styles - VS Code white
func operatorStyle(t chroma.TokenType, base tcell.Style) tcell.Style {
	switch t {
	case chroma.OperatorWord:
		return base.Foreground(tcell.NewRGBColor(86, 156, 214)) // Light blue
	default:
		return base.Foreground(tcell.NewRGBColor(212, 212, 212)) // Light gray
	}
}

// Literal styles
func literalStyle(t chroma.TokenType, base tcell.Style) tcell.Style {
	switch t {
	case chroma.LiteralDate:
		return base.Foreground(tcell.NewRGBColor(206, 145, 120)) // Orange
	default:
		return base.Foreground(tcell.NewRGBColor(181, 206, 168)) // Light green
	}
}

// Generic styles (markdown, diffs, etc)
func genericStyle(t chroma.TokenType, base tcell.Style) tcell.Style {
	switch t {
	case chroma.GenericHeading, chroma.GenericSubheading:
		return base.Foreground(tcell.NewRGBColor(86, 156, 214)).Bold(true) // Blue
	case chroma.GenericEmph:
		return base.Italic(true)
	case chroma.GenericStrong:
		return base.Bold(true)
	case chroma.GenericDeleted:
		return base.Foreground(tcell.NewRGBColor(244, 71, 71)) // Red
	case chroma.GenericInserted:
		return base.Foreground(tcell.NewRGBColor(106, 153, 85)) // Green
	case chroma.GenericError, chroma.GenericTraceback:
		return base.Foreground(tcell.NewRGBColor(244, 71, 71)).Bold(true) // Red
	case chroma.GenericPrompt:
		return base.Foreground(tcell.NewRGBColor(106, 153, 85)) // Green
	case chroma.GenericOutput:
		return base.Foreground(tcell.NewRGBColor(212, 212, 212)) // Light gray
	default:
		return base
	}
}

// Error style
func errorStyle(base tcell.Style) tcell.Style {
	return base.Foreground(tcell.NewRGBColor(244, 71, 71)).Underline(true) // Red underline
}
