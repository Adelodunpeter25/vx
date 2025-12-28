package buffer

func (b *Buffer) InsertRune(line, col int, r rune) {
	if line < 0 || line >= len(b.lines) {
		return
	}
	
	lineStr := b.lines[line]
	if col < 0 || col > len(lineStr) {
		col = len(lineStr)
	}
	
	b.lines[line] = lineStr[:col] + string(r) + lineStr[col:]
	b.markModified()
}

func (b *Buffer) DeleteRune(line, col int) {
	if line < 0 || line >= len(b.lines) {
		return
	}
	
	lineStr := b.lines[line]
	if col <= 0 || col > len(lineStr) {
		return
	}
	
	b.lines[line] = lineStr[:col-1] + lineStr[col:]
	b.markModified()
}

func (b *Buffer) InsertLine(line int) {
	if line < 0 || line > len(b.lines) {
		return
	}
	
	b.lines = append(b.lines[:line], append([]string{""}, b.lines[line:]...)...)
	b.markModified()
}

func (b *Buffer) DeleteLine(line int) {
	if line < 0 || line >= len(b.lines) || len(b.lines) == 1 {
		return
	}
	
	b.lines = append(b.lines[:line], b.lines[line+1:]...)
	b.markModified()
}

func (b *Buffer) SplitLine(line, col int) {
	if line < 0 || line >= len(b.lines) {
		return
	}
	
	lineStr := b.lines[line]
	if col < 0 || col > len(lineStr) {
		col = len(lineStr)
	}
	
	b.lines[line] = lineStr[:col]
	b.lines = append(b.lines[:line+1], append([]string{lineStr[col:]}, b.lines[line+1:]...)...)
	b.markModified()
}

func (b *Buffer) JoinLine(line int) {
	if line < 0 || line >= len(b.lines)-1 {
		return
	}
	
	b.lines[line] = b.lines[line] + b.lines[line+1]
	b.lines = append(b.lines[:line+1], b.lines[line+2:]...)
	b.markModified()
}
