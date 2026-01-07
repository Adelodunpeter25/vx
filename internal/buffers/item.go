package buffers

import (
	"github.com/Adelodunpeter25/vx/internal/buffer"
	"github.com/Adelodunpeter25/vx/internal/syntax"
)

// BufferItem represents a single buffer with its associated state
type BufferItem struct {
	Buffer   *buffer.Buffer
	Syntax   *syntax.Engine
	CursorX  int
	CursorY  int
	OffsetX  int
	OffsetY  int
}

// NewBufferItem creates a new buffer item
func NewBufferItem(buf *buffer.Buffer, filename string) *BufferItem {
	return &BufferItem{
		Buffer:  buf,
		Syntax:  syntax.New(filename),
		CursorX: 0,
		CursorY: 0,
		OffsetX: 0,
		OffsetY: 0,
	}
}
