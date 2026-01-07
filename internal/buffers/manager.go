package buffers

import (
	"fmt"

	"github.com/Adelodunpeter25/vx/internal/buffer"
)

// Manager manages multiple buffers
type Manager struct {
	buffers []*BufferItem
	current int
}

// New creates a new buffer manager with initial buffer
func New(buf *buffer.Buffer, filename string) *Manager {
	return &Manager{
		buffers: []*BufferItem{NewBufferItem(buf, filename)},
		current: 0,
	}
}

// Current returns the current buffer item
func (m *Manager) Current() *BufferItem {
	if m.current >= 0 && m.current < len(m.buffers) {
		return m.buffers[m.current]
	}
	return nil
}

// Count returns the number of buffers
func (m *Manager) Count() int {
	return len(m.buffers)
}

// CurrentIndex returns the current buffer index (1-based)
func (m *Manager) CurrentIndex() int {
	return m.current + 1
}

// Add adds a new buffer and switches to it
func (m *Manager) Add(buf *buffer.Buffer, filename string) {
	item := NewBufferItem(buf, filename)
	m.buffers = append(m.buffers, item)
	m.current = len(m.buffers) - 1
}

// Next switches to the next buffer
func (m *Manager) Next() {
	if len(m.buffers) > 1 {
		m.current = (m.current + 1) % len(m.buffers)
	}
}

// Previous switches to the previous buffer
func (m *Manager) Previous() {
	if len(m.buffers) > 1 {
		m.current = (m.current - 1 + len(m.buffers)) % len(m.buffers)
	}
}

// Delete removes the current buffer
func (m *Manager) Delete() error {
	if len(m.buffers) == 1 {
		return fmt.Errorf("cannot close last buffer")
	}

	// Remove current buffer
	m.buffers = append(m.buffers[:m.current], m.buffers[m.current+1:]...)

	// Adjust current index
	if m.current >= len(m.buffers) {
		m.current = len(m.buffers) - 1
	}

	return nil
}
