package visual

// Position represents a cursor position in the buffer
type Position struct {
	Line int
	Col  int
}

// Selection represents a text selection
type Selection struct {
	start  *Position
	end    *Position
	active bool
}

// New creates a new selection manager
func New() *Selection {
	return &Selection{
		active: false,
	}
}

// Start begins a selection at the given position
func (s *Selection) Start(line, col int) {
	s.start = &Position{Line: line, Col: col}
	s.end = &Position{Line: line, Col: col}
	s.active = true
}

// Update updates the end position of the selection
func (s *Selection) Update(line, col int) {
	if s.active && s.start != nil {
		s.end = &Position{Line: line, Col: col}
	}
}

// Clear clears the selection
func (s *Selection) Clear() {
	s.start = nil
	s.end = nil
	s.active = false
}

// IsActive returns whether a selection is active
func (s *Selection) IsActive() bool {
	return s.active && s.start != nil && s.end != nil
}

// GetRange returns the normalized start and end positions
func (s *Selection) GetRange() (startLine, startCol, endLine, endCol int, ok bool) {
	if !s.IsActive() {
		return 0, 0, 0, 0, false
	}

	// Normalize so start is always before end
	if s.start.Line < s.end.Line || (s.start.Line == s.end.Line && s.start.Col <= s.end.Col) {
		return s.start.Line, s.start.Col, s.end.Line, s.end.Col, true
	}
	return s.end.Line, s.end.Col, s.start.Line, s.start.Col, true
}

// Contains checks if a position is within the selection
func (s *Selection) Contains(line, col int) bool {
	startLine, startCol, endLine, endCol, ok := s.GetRange()
	if !ok {
		return false
	}

	if line < startLine || line > endLine {
		return false
	}
	if line == startLine && col < startCol {
		return false
	}
	if line == endLine && col > endCol {
		return false
	}
	return true
}
