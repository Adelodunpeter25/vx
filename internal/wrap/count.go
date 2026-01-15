package wrap

// VisualLineCount returns how many screen rows a line takes when wrapped
func VisualLineCount(text string, maxWidth int) int {
	if maxWidth <= 0 {
		return 1
	}
	
	length := len([]rune(text))
	if length == 0 {
		return 1
	}
	
	count := length / maxWidth
	if length%maxWidth != 0 {
		count++
	}
	return count
}

// TotalVisualLines calculates total screen rows for a range of buffer lines
func TotalVisualLines(lines []string, startLine, endLine, maxWidth int) int {
	total := 0
	for i := startLine; i <= endLine && i < len(lines); i++ {
		total += VisualLineCount(lines[i], maxWidth)
	}
	return total
}
