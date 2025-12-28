package buffer

type Buffer struct {
	lines    []string
	filename string
	modified bool
}

func New() *Buffer {
	return &Buffer{
		lines: []string{""},
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
