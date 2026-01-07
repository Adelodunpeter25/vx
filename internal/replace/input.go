package replace

import "github.com/Adelodunpeter25/vx/internal/search"

// GetSearchTerm returns the search term
func (e *Engine) GetSearchTerm() string {
	return e.searchTerm
}

// GetReplaceTerm returns the replacement term
func (e *Engine) GetReplaceTerm() string {
	return e.replaceTerm
}

// AppendToSearch adds a character to search term
func (e *Engine) AppendToSearch(r rune) {
	e.searchTerm += string(r)
}

// BackspaceSearch removes last character from search term
func (e *Engine) BackspaceSearch() {
	if len(e.searchTerm) > 0 {
		e.searchTerm = e.searchTerm[:len(e.searchTerm)-1]
	}
}

// AppendToReplace adds a character to replace term
func (e *Engine) AppendToReplace(r rune) {
	e.replaceTerm += string(r)
}

// BackspaceReplace removes last character from replace term
func (e *Engine) BackspaceReplace() {
	if len(e.replaceTerm) > 0 {
		e.replaceTerm = e.replaceTerm[:len(e.replaceTerm)-1]
	}
}

// ConfirmSearch moves to replace input state
func (e *Engine) ConfirmSearch(matches []search.Match) {
	e.matches = matches
	if len(matches) > 0 {
		e.state = StateReplaceInput
	} else {
		e.state = StateInactive
	}
}

// ConfirmReplace moves to confirmation state
func (e *Engine) ConfirmReplace() {
	if len(e.matches) > 0 {
		e.state = StateConfirm
		e.currentIdx = 0
	} else {
		e.state = StateInactive
	}
}
