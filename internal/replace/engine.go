package replace

import "github.com/Adelodunpeter25/vx/internal/search"

// Engine manages find and replace operations
type Engine struct {
	state       State
	searchTerm  string
	replaceTerm string
	matches     []search.Match
	currentIdx  int
}

func New() *Engine {
	return &Engine{
		state: StateInactive,
	}
}

// Start begins a new replace operation
func (e *Engine) Start() {
	e.state = StateSearchInput
	e.searchTerm = ""
	e.replaceTerm = ""
	e.matches = nil
	e.currentIdx = 0
}

// IsActive returns true if replace mode is active
func (e *Engine) IsActive() bool {
	return e.state != StateInactive
}

// GetState returns current state
func (e *Engine) GetState() State {
	return e.state
}

// Cancel cancels the replace operation
func (e *Engine) Cancel() {
	e.state = StateInactive
	e.searchTerm = ""
	e.replaceTerm = ""
	e.matches = nil
	e.currentIdx = 0
}
