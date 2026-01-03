package highlight

import (
	"strings"

	"github.com/alecthomas/chroma/v2"
	"github.com/gdamore/tcell/v2"
)

func tokenToStyle(tokenType chroma.TokenType) tcell.Style {
	base := tcell.StyleDefault
	typeStr := tokenType.String()

	// Comments (highest priority - check first)
	if isComment(tokenType, typeStr) {
		return commentStyle(base)
	}

	// Strings
	if isString(tokenType, typeStr) {
		return stringStyle(tokenType, base)
	}

	// Numbers
	if isNumber(tokenType, typeStr) {
		return numberStyle(tokenType, base)
	}

	// Keywords
	if isKeyword(tokenType, typeStr) {
		return keywordStyle(tokenType, base)
	}

	// Names (functions, classes, variables)
	if isName(tokenType, typeStr) {
		return nameStyle(tokenType, base)
	}

	// Operators and punctuation
	if isOperator(tokenType, typeStr) {
		return operatorStyle(tokenType, base)
	}

	// Generic tokens
	if isGeneric(tokenType, typeStr) {
		return genericStyle(tokenType, base)
	}

	// Literals
	if isLiteral(tokenType, typeStr) {
		return literalStyle(tokenType, base)
	}

	// Text and whitespace
	if isText(tokenType, typeStr) {
		return base
	}

	// Error tokens
	if isError(tokenType, typeStr) {
		return errorStyle(base)
	}

	return base
}

func isKeyword(t chroma.TokenType, s string) bool {
	return t == chroma.Keyword ||
		t == chroma.KeywordConstant ||
		t == chroma.KeywordDeclaration ||
		t == chroma.KeywordNamespace ||
		t == chroma.KeywordPseudo ||
		t == chroma.KeywordReserved ||
		t == chroma.KeywordType ||
		strings.Contains(s, "Keyword")
}

func isString(t chroma.TokenType, s string) bool {
	return t == chroma.String ||
		t == chroma.LiteralString ||
		t == chroma.LiteralStringDouble ||
		t == chroma.LiteralStringSingle ||
		t == chroma.LiteralStringBacktick ||
		t == chroma.LiteralStringChar ||
		t == chroma.LiteralStringDoc ||
		t == chroma.LiteralStringEscape ||
		t == chroma.LiteralStringHeredoc ||
		t == chroma.LiteralStringInterpol ||
		t == chroma.LiteralStringOther ||
		t == chroma.LiteralStringRegex ||
		t == chroma.LiteralStringSymbol ||
		strings.Contains(s, "String")
}

func isNumber(t chroma.TokenType, s string) bool {
	return t == chroma.Number ||
		t == chroma.LiteralNumber ||
		t == chroma.LiteralNumberBin ||
		t == chroma.LiteralNumberFloat ||
		t == chroma.LiteralNumberHex ||
		t == chroma.LiteralNumberInteger ||
		t == chroma.LiteralNumberIntegerLong ||
		t == chroma.LiteralNumberOct ||
		strings.Contains(s, "Number")
}

func isComment(t chroma.TokenType, s string) bool {
	return t == chroma.Comment ||
		t == chroma.CommentHashbang ||
		t == chroma.CommentMultiline ||
		t == chroma.CommentPreproc ||
		t == chroma.CommentPreprocFile ||
		t == chroma.CommentSingle ||
		t == chroma.CommentSpecial ||
		strings.Contains(s, "Comment")
}

func isName(t chroma.TokenType, s string) bool {
	return t == chroma.Name ||
		t == chroma.NameAttribute ||
		t == chroma.NameBuiltin ||
		t == chroma.NameBuiltinPseudo ||
		t == chroma.NameClass ||
		t == chroma.NameConstant ||
		t == chroma.NameDecorator ||
		t == chroma.NameEntity ||
		t == chroma.NameException ||
		t == chroma.NameFunction ||
		t == chroma.NameFunctionMagic ||
		t == chroma.NameLabel ||
		t == chroma.NameNamespace ||
		t == chroma.NameOther ||
		t == chroma.NameProperty ||
		t == chroma.NameTag ||
		t == chroma.NameVariable ||
		t == chroma.NameVariableClass ||
		t == chroma.NameVariableGlobal ||
		t == chroma.NameVariableInstance ||
		t == chroma.NameVariableMagic ||
		strings.Contains(s, "Name")
}

func isOperator(t chroma.TokenType, s string) bool {
	return t == chroma.Operator ||
		t == chroma.OperatorWord ||
		t == chroma.Punctuation ||
		strings.Contains(s, "Operator") ||
		strings.Contains(s, "Punctuation")
}

func isLiteral(t chroma.TokenType, s string) bool {
	return t == chroma.Literal ||
		t == chroma.LiteralDate ||
		strings.Contains(s, "Literal")
}

func isGeneric(t chroma.TokenType, s string) bool {
	return t == chroma.Generic ||
		t == chroma.GenericDeleted ||
		t == chroma.GenericEmph ||
		t == chroma.GenericError ||
		t == chroma.GenericHeading ||
		t == chroma.GenericInserted ||
		t == chroma.GenericOutput ||
		t == chroma.GenericPrompt ||
		t == chroma.GenericStrong ||
		t == chroma.GenericSubheading ||
		t == chroma.GenericTraceback ||
		strings.Contains(s, "Generic")
}

func isText(t chroma.TokenType, s string) bool {
	return t == chroma.Text ||
		t == chroma.TextWhitespace ||
		strings.Contains(s, "Text") ||
		strings.Contains(s, "Whitespace")
}

func isError(t chroma.TokenType, s string) bool {
	return t == chroma.Error ||
		strings.Contains(s, "Error")
}
