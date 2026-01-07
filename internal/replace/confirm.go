package replace

import "github.com/Adelodunpeter25/vx/internal/search"

// GetCurrentMatch returns the current match being confirmed
func (e *Engine) GetCurrentMatch() *search.Match {
	if e.currentIdx >= 0 && e.currentIdx < len(e.matches) {
		return &e.matches[e.currentIdx]
	}
	return nil
}

// GetMatchCount returns total number of matches
func (e *Engine) GetMatchCount() int {
	return len(e.matches)
}

// GetCurrentIndex returns current match index (1-based)
func (e *Engine) GetCurrentIndex() int {
	return e.currentIdx + 1
}

// NextMatch moves to next match
func (e *Engine) NextMatch() bool {
	e.currentIdx++
	if e.currentIdx >= len(e.matches) {
		e.state = StateInactive
		return false
	}
	return true
}
