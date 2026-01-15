package wrap

// WrapLine wraps a single line into segments that fit within maxWidth
func WrapLine(text string, lineNum, maxWidth int) []Line {
	if maxWidth <= 0 {
		return []Line{{Text: text, StartCol: 0, LineNum: lineNum, IsWrapped: false}}
	}
	
	runes := []rune(text)
	if len(runes) <= maxWidth {
		return []Line{{Text: text, StartCol: 0, LineNum: lineNum, IsWrapped: false}}
	}
	
	var segments []Line
	startCol := 0
	
	for startCol < len(runes) {
		endCol := startCol + maxWidth
		if endCol > len(runes) {
			endCol = len(runes)
		}
		
		segments = append(segments, Line{
			Text:      string(runes[startCol:endCol]),
			StartCol:  startCol,
			LineNum:   lineNum,
			IsWrapped: startCol > 0,
		})
		
		startCol = endCol
	}
	
	return segments
}
