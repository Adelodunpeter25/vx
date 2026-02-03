package terminalpane

import "sync"

type Scrollback struct {
	mu     sync.RWMutex
	lines  [][]CellSnapshot
	maxLen int
	offset int
}

type CellSnapshot struct {
	Ch    rune
	Style uint64
}

func NewScrollback(max int) *Scrollback {
	if max < 1 {
		max = 1000
	}
	return &Scrollback{maxLen: max}
}

func (s *Scrollback) Append(line []CellSnapshot) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.lines = append(s.lines, line)
	if len(s.lines) > s.maxLen {
		s.lines = s.lines[len(s.lines)-s.maxLen:]
	}
}

func (s *Scrollback) SetOffset(offset int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if offset < 0 {
		offset = 0
	}
	if offset > len(s.lines) {
		offset = len(s.lines)
	}
	s.offset = offset
}

func (s *Scrollback) Offset() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.offset
}

func (s *Scrollback) GetLine(idx int) []CellSnapshot {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if idx < 0 || idx >= len(s.lines) {
		return nil
	}
	return s.lines[idx]
}

func (s *Scrollback) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.lines)
}
