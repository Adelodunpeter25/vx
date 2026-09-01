package buffer

import (
	"fmt"
	"unicode/utf8"
	"time"

	"github.com/Adelodunpeter25/vx/internal/undo"
	"github.com/Adelodunpeter25/vx/internal/utils"
)

type Buffer struct {
	lines      []string
	filename   string
	modified   bool
	modVersion int // Increments on each modification
	undoStack  *undo.Stack
	lazy       *utils.LazyFileReader
	totalLines int
	modTime    time.Time
	fileSize   int64
}

func New() *Buffer {
	return &Buffer{
		lines:      []string{""},
		modVersion: 0,
		undoStack:  undo.NewStack(),
		totalLines: 1,
	}
}

func (b *Buffer) LineCount() int {
	if b.lazy != nil {
		if b.totalLines > 0 {
			return b.totalLines
		}
		return 0
	}
	return len(b.lines)
}

func (b *Buffer) Line(n int) string {
	if b.lazy != nil {
		if n < 0 || (b.totalLines > 0 && n >= b.totalLines) {
			return ""
		}
		b.ensureLineLoaded(n)
	}
	if n < 0 || n >= len(b.lines) {
		return ""
	}
	return b.lines[n]
}

func (b *Buffer) LineRuneCount(n int) int {
	line := b.Line(n)
	if line == "" && (n < 0 || n >= b.LineCount()) {
		return 0
	}
	return utf8.RuneCountInString(line)
}

func (b *Buffer) IsModified() bool {
	return b.modified
}

func (b *Buffer) Filename() string {
	return b.filename
}

func (b *Buffer) IsLazy() bool {
	return b.lazy != nil
}

func (b *Buffer) SetFilename(filename string) {
	b.filename = filename
}

func (b *Buffer) ModVersion() int {
	return b.modVersion
}

func (b *Buffer) ModTime() time.Time {
	return b.modTime
}

func (b *Buffer) setModTime(t time.Time) {
	b.modTime = t
}

func (b *Buffer) markModified() {
	b.modified = true
	b.modVersion++
}

func (b *Buffer) UndoStack() *undo.Stack {
	return b.undoStack
}

func (b *Buffer) ensureLineLoaded(n int) {
	if b.lazy == nil {
		return
	}
	for n >= len(b.lines) && !b.lazy.IsFullyLoaded() {
		chunk, err := b.lazy.LoadChunk()
		if err != nil {
			break
		}
		if len(chunk) == 0 {
			break
		}
		b.lines = append(b.lines, chunk...)
	}
	if b.lazy.IsFullyLoaded() {
		_ = b.lazy.Close()
		b.lazy = nil
	}
}

// Reload re-reads the file from disk into this buffer.
// The undo stack is reset. Returns an error if the reload fails.
func (b *Buffer) Reload() error {
	if b.filename == "" {
		return fmt.Errorf("no filename set")
	}
	newBuf, err := Load(b.filename)
	if err != nil {
		return err
	}
	b.lines = newBuf.lines
	b.lazy = newBuf.lazy
	b.totalLines = newBuf.totalLines
	b.modified = false
	b.modVersion++
	b.undoStack = undo.NewStack()
	b.modTime = newBuf.modTime
	b.fileSize = newBuf.fileSize
	return nil
}

func (b *Buffer) ensureAllLoaded() {
	if b.lazy == nil {
		return
	}
	for !b.lazy.IsFullyLoaded() {
		_, err := b.lazy.LoadChunk()
		if err != nil {
			break
		}
	}
	if b.lazy.IsFullyLoaded() {
		b.lines = append(b.lines, b.lazy.GetLines()[len(b.lines):]...)
		_ = b.lazy.Close()
		b.lazy = nil
	}
}
