package buffer

import "github.com/Adelodunpeter25/vx/internal/undo"

type Buffer struct {
	lines      []string
	filename   string
	modified   bool
	modVersion int // Increments on each modification
	undoStack  *undo.Stack
}

func New() *Buffer {
	return &Buffer{
		lines:      []string{""},
		modVersion: 0,
		undoStack:  undo.NewStack(),
	}
}

func (b *Buffer) LineCount() int {
	return len(b.lines)
}

func (b *Buffer) Line(n int) string {
	if n < 0 || n >= len(b.lines) {
		return ""
	}
	return b.lines[n]
}

func (b *Buffer) IsModified() bool {
	return b.modified
}

func (b *Buffer) Filename() string {
	return b.filename
}

func (b *Buffer) SetFilename(filename string) {
	b.filename = filename
}

func (b *Buffer) ModVersion() int {
	return b.modVersion
}

func (b *Buffer) markModified() {
	b.modified = true
	b.modVersion++
}

func (b *Buffer) UndoStack() *undo.Stack {
	return b.undoStack
}
