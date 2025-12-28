package markdown

import (
	"regexp"
	"strings"
)

// Element represents a parsed markdown element
type Element struct {
	Type    ElementType
	Content string
	Level   int // For headers
}

type ElementType int

const (
	TypeText ElementType = iota
	TypeHeader
	TypeBold
	TypeItalic
	TypeCode
	TypeCodeBlock
	TypeList
	TypeLink
	TypeBlockquote
)

// Parse converts markdown text to structured elements
func Parse(text string) []Element {
	lines := strings.Split(text, "\n")
	elements := make([]Element, 0)
	inCodeBlock := false
	
	for _, line := range lines {
		// Code blocks
		if strings.HasPrefix(line, "```") {
			inCodeBlock = !inCodeBlock
			if !inCodeBlock {
				continue
			}
			elements = append(elements, Element{Type: TypeCodeBlock, Content: ""})
			continue
		}
		
		if inCodeBlock {
			elements = append(elements, Element{Type: TypeCodeBlock, Content: line})
			continue
		}
		
		// Headers
		if strings.HasPrefix(line, "#") {
			level := 0
			for _, r := range line {
				if r == '#' {
					level++
				} else {
					break
				}
			}
			content := strings.TrimSpace(line[level:])
			elements = append(elements, Element{Type: TypeHeader, Content: content, Level: level})
			continue
		}
		
		// Lists
		if strings.HasPrefix(strings.TrimSpace(line), "- ") || strings.HasPrefix(strings.TrimSpace(line), "* ") {
			content := strings.TrimSpace(line[2:])
			elements = append(elements, Element{Type: TypeList, Content: content})
			continue
		}
		
		// Blockquotes
		if strings.HasPrefix(strings.TrimSpace(line), "> ") {
			content := strings.TrimSpace(line[2:])
			elements = append(elements, Element{Type: TypeBlockquote, Content: content})
			continue
		}
		
		// Regular text with inline formatting
		if line != "" {
			elements = append(elements, Element{Type: TypeText, Content: line})
		}
	}
	
	return elements
}

// ParseInline handles inline markdown formatting (bold, italic, code, links)
func ParseInline(text string) string {
	// Bold **text**
	boldRe := regexp.MustCompile(`\*\*(.+?)\*\*`)
	text = boldRe.ReplaceAllString(text, "$1")
	
	// Italic *text*
	italicRe := regexp.MustCompile(`\*(.+?)\*`)
	text = italicRe.ReplaceAllString(text, "$1")
	
	// Inline code `code`
	codeRe := regexp.MustCompile("`(.+?)`")
	text = codeRe.ReplaceAllString(text, "$1")
	
	// Links [text](url)
	linkRe := regexp.MustCompile(`\[(.+?)\]\((.+?)\)`)
	text = linkRe.ReplaceAllString(text, "$1")
	
	return text
}
