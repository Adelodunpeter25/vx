package buffer

import "github.com/Adelodunpeter25/vx/internal/undo"

func (b *Buffer) InsertRune(line, col int, r rune) {
	if line < 0 || line >= len(b.lines) {
		return
	}
	
	lineStr := b.lines[line]
	if col < 0 || col > len(lineStr) {
		col = len(lineStr)
	}
	
	// Record undo action
	b.undoStack.Push(undo.Action{
		Type: undo.ActionInsertRune,
		Line: line,
		Col:  col,
		Text: string(r),
	})
	
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
	
	// Record undo action
	b.undoStack.Push(undo.Action{
		Type:    undo.ActionDeleteRune,
		Line:    line,
		Col:     col,
		OldText: string(lineStr[col-1]),
	})
	
	b.lines[line] = lineStr[:col-1] + lineStr[col:]
	b.markModified()
}

func (b *Buffer) InsertLine(line int) {
	if line < 0 || line > len(b.lines) {
		return
	}
	
	// Record undo action
	b.undoStack.Push(undo.Action{
		Type: undo.ActionInsertLine,
		Line: line,
	})
	
	b.lines = append(b.lines[:line], append([]string{""}, b.lines[line:]...)...)
	b.markModified()
}

func (b *Buffer) DeleteLine(line int) {
	if line < 0 || line >= len(b.lines) || len(b.lines) == 1 {
		return
	}
	
	// Record undo action
	b.undoStack.Push(undo.Action{
		Type:    undo.ActionDeleteLine,
		Line:    line,
		OldText: b.lines[line],
	})
	
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
	
	// Record undo action
	b.undoStack.Push(undo.Action{
		Type: undo.ActionSplitLine,
		Line: line,
		Col:  col,
	})
	
	b.lines[line] = lineStr[:col]
	b.lines = append(b.lines[:line+1], append([]string{lineStr[col:]}, b.lines[line+1:]...)...)
	b.markModified()
}

func (b *Buffer) JoinLine(line int) {
	if line < 0 || line >= len(b.lines)-1 {
		return
	}
	
	// Record undo action
	b.undoStack.Push(undo.Action{
		Type:    undo.ActionJoinLine,
		Line:    line,
		OldText: b.lines[line+1],
	})
	
	b.lines[line] = b.lines[line] + b.lines[line+1]
	b.lines = append(b.lines[:line+1], b.lines[line+2:]...)
	b.markModified()
}

// Undo operations (without recording to undo stack)
func (b *Buffer) undoInsertRune(line, col int) {
	if line < 0 || line >= len(b.lines) {
		return
	}
	lineStr := b.lines[line]
	if col < 0 || col >= len(lineStr) {
		return
	}
	b.lines[line] = lineStr[:col] + lineStr[col+1:]
}

func (b *Buffer) undoDeleteRune(line, col int, r string) {
	if line < 0 || line >= len(b.lines) {
		return
	}
	lineStr := b.lines[line]
	if col < 0 || col > len(lineStr) {
		col = len(lineStr)
	}
	b.lines[line] = lineStr[:col-1] + r + lineStr[col-1:]
}

func (b *Buffer) undoInsertLine(line int) {
	if line < 0 || line >= len(b.lines) || len(b.lines) == 1 {
		return
	}
	b.lines = append(b.lines[:line], b.lines[line+1:]...)
}

func (b *Buffer) undoDeleteLine(line int, text string) {
	if line < 0 || line > len(b.lines) {
		return
	}
	b.lines = append(b.lines[:line], append([]string{text}, b.lines[line:]...)...)
}

func (b *Buffer) undoSplitLine(line, col int) {
	if line < 0 || line >= len(b.lines)-1 {
		return
	}
	b.lines[line] = b.lines[line] + b.lines[line+1]
	b.lines = append(b.lines[:line+1], b.lines[line+2:]...)
}

func (b *Buffer) undoJoinLine(line int, text string) {
	if line < 0 || line >= len(b.lines) {
		return
	}
	lineStr := b.lines[line]
	splitPos := len(lineStr) - len(text)
	if splitPos < 0 {
		splitPos = 0
	}
	b.lines[line] = lineStr[:splitPos]
	b.lines = append(b.lines[:line+1], append([]string{text}, b.lines[line+1:]...)...)
}

// Undo performs an undo operation
func (b *Buffer) Undo() bool {
	action := b.undoStack.Undo()
	if action == nil {
		return false
	}
	
	switch action.Type {
	case undo.ActionInsertRune:
		b.undoInsertRune(action.Line, action.Col)
	case undo.ActionDeleteRune:
		b.undoDeleteRune(action.Line, action.Col, action.OldText)
	case undo.ActionInsertLine:
		b.undoInsertLine(action.Line)
	case undo.ActionDeleteLine:
		b.undoDeleteLine(action.Line, action.OldText)
	case undo.ActionSplitLine:
		b.undoSplitLine(action.Line, action.Col)
	case undo.ActionJoinLine:
		b.undoJoinLine(action.Line, action.OldText)
	}
	
	b.markModified()
	return true
}

// Redo performs a redo operation
func (b *Buffer) Redo() bool {
	action := b.undoStack.Redo()
	if action == nil {
		return false
	}
	
	switch action.Type {
	case undo.ActionInsertRune:
		lineStr := b.lines[action.Line]
		b.lines[action.Line] = lineStr[:action.Col] + action.Text + lineStr[action.Col:]
	case undo.ActionDeleteRune:
		lineStr := b.lines[action.Line]
		b.lines[action.Line] = lineStr[:action.Col-1] + lineStr[action.Col:]
	case undo.ActionInsertLine:
		b.lines = append(b.lines[:action.Line], append([]string{""}, b.lines[action.Line:]...)...)
	case undo.ActionDeleteLine:
		b.lines = append(b.lines[:action.Line], b.lines[action.Line+1:]...)
	case undo.ActionSplitLine:
		lineStr := b.lines[action.Line]
		b.lines[action.Line] = lineStr[:action.Col]
		b.lines = append(b.lines[:action.Line+1], append([]string{lineStr[action.Col:]}, b.lines[action.Line+1:]...)...)
	case undo.ActionJoinLine:
		b.lines[action.Line] = b.lines[action.Line] + b.lines[action.Line+1]
		b.lines = append(b.lines[:action.Line+1], b.lines[action.Line+2:]...)
	}
	
	b.markModified()
	return true
}
