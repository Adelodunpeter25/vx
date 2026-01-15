package wrap

// ScreenPosition converts buffer position to screen position
func ScreenPosition(line string, bufferCol, maxWidth int) (screenRow, screenCol int) {
	if maxWidth <= 0 || bufferCol <= 0 {
		return 0, bufferCol
	}
	
	screenRow = bufferCol / maxWidth
	screenCol = bufferCol % maxWidth
	return
}

// BufferPosition converts screen position back to buffer position
func BufferPosition(screenRow, screenCol, maxWidth int) int {
	if maxWidth <= 0 {
		return screenCol
	}
	return screenRow*maxWidth + screenCol
}
