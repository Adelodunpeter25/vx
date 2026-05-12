package markdown

import "strings"

type Element struct {
	Type     ElementType
	Content  string
	Level    int
	Segments []Segment
	Number   int
	Checked  *bool
	URL      string
	Lang     string
}

type Segment struct {
	Text         string
	Bold         bool
	Italic       bool
	Code         bool
	Strikethrough bool
	Link         bool
	LinkURL      string
}

type ElementType int

const (
	TypeText ElementType = iota
	TypeHeader
	TypeCodeBlock
	TypeList
	TypeBlockquote
	TypeHorizontalRule
	TypeOrderedList
	TypeTaskList
	TypeParagraph
	TypeEmptyLine
)

func Parse(text string) []Element {
	lines := strings.Split(text, "\n")
	var elements []Element
	var paragraphLines []string
	inCodeBlock := false

	flushParagraph := func() {
		if len(paragraphLines) > 0 {
			content := strings.Join(paragraphLines, "\n")
			segments := parseInlineFormatting(content)
			elements = append(elements, Element{
				Type:     TypeParagraph,
				Content:  content,
				Segments: segments,
			})
			paragraphLines = nil
		}
	}

	for i := 0; i < len(lines); i++ {
		line := lines[i]

		if strings.HasPrefix(line, "```") {
			if inCodeBlock {
				inCodeBlock = false
				continue
			}
			flushParagraph()
			inCodeBlock = true
			continue
		}

		if inCodeBlock {
			elements = append(elements, Element{Type: TypeCodeBlock, Content: line})
			continue
		}

		trimmed := strings.TrimSpace(line)

		if trimmed == "" {
			flushParagraph()
			elements = append(elements, Element{Type: TypeEmptyLine})
			continue
		}

		if isHorizontalRule(trimmed) {
			flushParagraph()
			elements = append(elements, Element{Type: TypeHorizontalRule})
			continue
		}

		if strings.HasPrefix(line, "#") && (len(line) == 1 || line[1] == ' ' || line[1] == '\t') {
			flushParagraph()
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
				Type:     TypeHeader,
				Content:  content,
				Level:    level,
				Segments: segments,
			})
			continue
		}

		if strings.HasPrefix(trimmed, "> ") {
			flushParagraph()
			content := strings.TrimSpace(trimmed[2:])
			segments := parseInlineFormatting(content)
			elements = append(elements, Element{
				Type:     TypeBlockquote,
				Content:  content,
				Segments: segments,
			})
			continue
		}

		if matched, num, content := matchOrderedList(trimmed); matched {
			flushParagraph()
			segments := parseInlineFormatting(content)
			elements = append(elements, Element{
				Type:     TypeOrderedList,
				Content:  content,
				Number:   num,
				Segments: segments,
			})
			continue
		}

		if matched, checked, content := matchTaskList(trimmed); matched {
			flushParagraph()
			segments := parseInlineFormatting(content)
			elements = append(elements, Element{
				Type:     TypeTaskList,
				Content:  content,
				Checked:  checked,
				Segments: segments,
			})
			continue
		}

		if strings.HasPrefix(trimmed, "- ") || strings.HasPrefix(trimmed, "* ") {
			flushParagraph()
			content := strings.TrimSpace(trimmed[2:])
			segments := parseInlineFormatting(content)
			elements = append(elements, Element{
				Type:     TypeList,
				Content:  content,
				Segments: segments,
			})
			continue
		}

		paragraphLines = append(paragraphLines, line)
	}

	if inCodeBlock {
		elements = append(elements, Element{Type: TypeCodeBlock, Content: ""})
	}
	flushParagraph()

	return elements
}

func isHorizontalRule(s string) bool {
	if len(s) < 3 {
		return false
	}
	first := s[0]
	if first != '-' && first != '*' && first != '_' {
		return false
	}
	for _, r := range s {
		if byte(r) != first {
			return false
		}
	}
	return true
}

func matchOrderedList(s string) (bool, int, string) {
	for i := 0; i < len(s); i++ {
		if s[i] >= '0' && s[i] <= '9' {
			continue
		}
		if s[i] == '.' && i > 0 && i+1 < len(s) && s[i+1] == ' ' {
			num := 0
			for j := 0; j < i; j++ {
				num = num*10 + int(s[j]-'0')
			}
			return true, num, strings.TrimSpace(s[i+2:])
		}
		break
	}
	return false, 0, ""
}

func matchTaskList(s string) (bool, *bool, string) {
	if !strings.HasPrefix(s, "- ") && !strings.HasPrefix(s, "* ") {
		return false, nil, ""
	}
	rest := s[2:]
	if strings.HasPrefix(rest, "[ ] ") {
		checked := false
		return true, &checked, strings.TrimSpace(rest[4:])
	}
	if strings.HasPrefix(rest, "[x] ") || strings.HasPrefix(rest, "[X] ") {
		checked := true
		return true, &checked, strings.TrimSpace(rest[4:])
	}
	return false, nil, ""
}

func parseInlineFormatting(text string) []Segment {
	var segments []Segment
	runes := []rune(text)
	i := 0

	for i < len(runes) {
		if runes[i] == '\\' && i+1 < len(runes) {
			segments = append(segments, Segment{Text: string(runes[i+1])})
			i += 2
			continue
		}

		if i+1 < len(runes) && runes[i] == '`' && runes[i+1] == '`' {
			end := i + 2
			for end < len(runes) {
				if end+1 < len(runes) && runes[end] == '`' && runes[end+1] == '`' {
					segments = append(segments, Segment{
						Text: string(runes[i+2 : end]),
						Code: true,
					})
					i = end + 2
					goto next
				}
				end++
			}
			segments = append(segments, Segment{Text: "``"})
			i += 2
			continue
		}

		if runes[i] == '`' {
			end := i + 1
			for end < len(runes) && runes[end] != '`' {
				end++
			}
			if end < len(runes) {
				segments = append(segments, Segment{
					Text: string(runes[i+1 : end]),
					Code: true,
				})
				i = end + 1
				continue
			}
			segments = append(segments, Segment{Text: "`"})
			i++
			continue
		}

		if i+1 < len(runes) && runes[i] == '[' {
			linkText, linkURL, consumed := tryParseLink(runes, i)
			if linkText != "" {
				segments = append(segments, Segment{
					Text:    linkText,
					Link:    true,
					LinkURL: linkURL,
				})
				i = consumed
				continue
			}
		}

		if i+2 < len(runes) && runes[i] == '~' && runes[i+1] == '~' {
			end := i + 2
			for end+1 < len(runes) {
				if runes[end] == '~' && runes[end+1] == '~' {
					inner := parseInlineFormatting(string(runes[i+2 : end]))
					for _, s := range inner {
						s.Strikethrough = true
						segments = append(segments, s)
					}
					i = end + 2
					goto next
				}
				end++
			}
		}

		if i+2 < len(runes) && runes[i] == '*' && runes[i+1] == '*' && runes[i+2] == '*' {
			end := i + 3
			for end+2 < len(runes) {
				if runes[end] == '*' && runes[end+1] == '*' && runes[end+2] == '*' {
					inner := parseInlineFormatting(string(runes[i+3 : end]))
					for _, s := range inner {
						s.Bold = true
						s.Italic = true
						segments = append(segments, s)
					}
					i = end + 3
					goto next
				}
				end++
			}
		}

		if i+2 < len(runes) && runes[i] == '*' && runes[i+1] == '*' {
			end := i + 2
			for end+1 < len(runes) {
				if runes[end] == '*' && runes[end+1] == '*' {
					inner := parseInlineFormatting(string(runes[i+2 : end]))
					for _, s := range inner {
						s.Bold = true
						segments = append(segments, s)
					}
					i = end + 2
					goto next
				}
				end++
			}
		}

		if runes[i] == '*' {
			end := i + 1
			for end < len(runes) && runes[end] != '*' {
				end++
			}
			if end < len(runes) && end > i+1 {
				inner := parseInlineFormatting(string(runes[i+1 : end]))
				for _, s := range inner {
					s.Italic = true
					segments = append(segments, s)
				}
				i = end + 1
				continue
			}
		}

		segments = append(segments, Segment{Text: string(runes[i])})
		i++

	next:
	}

	return mergeTextSegments(segments)
}

func tryParseLink(runes []rune, start int) (text, url string, consumed int) {
	if start+2 >= len(runes) || runes[start] != '[' {
		return "", "", start
	}
	closeBracket := start + 1
	for closeBracket < len(runes) && runes[closeBracket] != ']' {
		closeBracket++
	}
	if closeBracket >= len(runes) || closeBracket+1 >= len(runes) || runes[closeBracket+1] != '(' {
		return "", "", start
	}
	closeParen := closeBracket + 2
	for closeParen < len(runes) && runes[closeParen] != ')' {
		closeParen++
	}
	if closeParen >= len(runes) {
		return "", "", start
	}
	text = string(runes[start+1 : closeBracket])
	url = string(runes[closeBracket+2 : closeParen])
	return text, url, closeParen + 1
}

func mergeTextSegments(segments []Segment) []Segment {
	if len(segments) == 0 {
		return segments
	}
	merged := []Segment{segments[0]}
	for _, s := range segments[1:] {
		last := &merged[len(merged)-1]
		if !last.Bold && !last.Italic && !last.Code && !last.Strikethrough && !last.Link &&
			!s.Bold && !s.Italic && !s.Code && !s.Strikethrough && !s.Link {
			last.Text += s.Text
			continue
		}
		merged = append(merged, s)
	}
	return merged
}
