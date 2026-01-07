package buffers

// SaveState saves the current cursor and scroll position to the buffer item
func (m *Manager) SaveState(cursorX, cursorY, offsetX, offsetY int) {
	current := m.Current()
	if current != nil {
		current.CursorX = cursorX
		current.CursorY = cursorY
		current.OffsetX = offsetX
		current.OffsetY = offsetY
	}
}

// RestoreState returns the saved cursor and scroll position
func (m *Manager) RestoreState() (cursorX, cursorY, offsetX, offsetY int) {
	current := m.Current()
	if current != nil {
		return current.CursorX, current.CursorY, current.OffsetX, current.OffsetY
	}
	return 0, 0, 0, 0
}
