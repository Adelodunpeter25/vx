package visual

// Selection represents a text selection
type Selection struct {
	StartLine int
	StartCol  int
	EndLine   int
	EndCol    int
	Active    bool
}

// Manager handles visual selection state
type Manager struct {
	selection Selection
}

// New creates a new visual selection manager
func New() *Manager {
	return &Manager{
		selection: Selection{Active: false},
	}
}

// StartSelection begins a new selection at the given position
func (m *Manager) StartSelection(line, col int) {
	m.selection = Selection{
		StartLine: line,
		StartCol:  col,
		EndLine:   line,
		EndCol:    col,
		Active:    true,
	}
}

// UpdateSelection updates the end position of the selection
func (m *Manager) UpdateSelection(line, col int) {
	if m.selection.Active {
		m.selection.EndLine = line
		m.selection.EndCol = col
	}
}

// ClearSelection deactivates the selection
func (m *Manager) ClearSelection() {
	m.selection.Active = false
}

// IsActive returns true if selection is active
func (m *Manager) IsActive() bool {
	return m.selection.Active
}

// GetSelection returns the current selection
func (m *Manager) GetSelection() Selection {
	return m.selection
}

// HasSelection returns true if there's an active selection with content
func (m *Manager) HasSelection() bool {
	if !m.selection.Active {
		return false
	}
	
	// Check if selection has any content
	return m.selection.StartLine != m.selection.EndLine || 
		   m.selection.StartCol != m.selection.EndCol
}
