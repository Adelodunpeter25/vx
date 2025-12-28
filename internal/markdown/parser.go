package markdown

import (
	"regexp"
	"strings"
)

// Element represents a parsed markdown element
type Element struct {
	Type     ElementType
	Content  string
	Level    int // For headers
	Segments []Segment // For inline formatting
}

type Segment struct {
	Text  string
	Bold  bool
	Italic bool
	Code  bool
}

type ElementType int

const (
	TypeText ElementType = iota
	TypeHeader
	TypeCodeBlock
	TypeList
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
			segments := parseInlineFormatting(content)
			elements = append(elements, Element{
				Type: TypeHeader, 
				Content: content, 
				Level: level,
				Segments: segments,
			})
			continue
		}
		
		// Lists
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "- ") || strings.HasPrefix(trimmed, "* ") {
			content := strings.TrimSpace(trimmed[2:])
			segments := parseInlineFormatting(content)
			elements = append(elements, Element{
				Type: TypeList, 
				Content: content,
				Segments: segments,
			})
			continue
		}
		
		// Blockquotes
		if strings.HasPrefix(trimmed, "> ") {
			content := strings.TrimSpace(trimmed[2:])
			segments := parseInlineFormatting(content)
			elements = append(elements, Element{
				Type: TypeBlockquote, 
				Content: content,
				Segments: segments,
			})
			continue
		}
		
		// Regular text with inline formatting
		if line != "" {
			segments := parseInlineFormatting(line)
			elements = append(elements, Element{
				Type: TypeText, 
				Content: line,
				Segments: segments,
			})
		}
	}
	
	return elements
}

// parseInlineFormatting extracts bold, italic, and code segments
func parseInlineFormatting(text string) []Segment {
	segments := make([]Segment, 0)
	
	// Pattern to match **bold**, *italic*, and `code`
	pattern := regexp.MustCompile(`(\*\*[^*]+\*\*|\*[^*]+\*|` + "`" + `[^` + "`" + `]+` + "`" + `|[^*` + "`" + `]+)`)
	matches := pattern.FindAllString(text, -1)
	
	for _, match := range matches {
		seg := Segment{Text: match}
		
		// Check for bold **text**
		if strings.HasPrefix(match, "**") && strings.HasSuffix(match, "**") && len(match) > 4 {
			seg.Text = match[2:len(match)-2]
			seg.Bold = true
		} else if strings.HasPrefix(match, "*") && strings.HasSuffix(match, "*") && len(match) > 2 {
			// Check for italic *text*
			seg.Text = match[1:len(match)-1]
			seg.Italic = true
		} else if strings.HasPrefix(match, "`") && strings.HasSuffix(match, "`") && len(match) > 2 {
			// Check for code `text`
			seg.Text = match[1:len(match)-1]
			seg.Code = true
		}
		
		segments = append(segments, seg)
	}
	
	return segments
}
