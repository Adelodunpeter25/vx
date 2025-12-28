package search

import "strings"

// Match represents a search match location
type Match struct {
	Line int
	Col  int
	Len  int
}

// Engine handles search operations
type Engine struct {
	query   string
	matches []Match
	current int
}

func New() *Engine {
	return &Engine{
		matches: []Match{},
		current: -1,
	}
}

// Search finds all matches in buffer
func (e *Engine) Search(lines []string, query string) []Match {
	if query == "" {
		e.matches = []Match{}
		e.current = -1
		return e.matches
	}

	e.query = query
	e.matches = []Match{}
	e.current = -1

	// Case-insensitive search
	lowerQuery := strings.ToLower(query)

	for lineNum, line := range lines {
		lowerLine := strings.ToLower(line)
		col := 0
		for {
			idx := strings.Index(lowerLine[col:], lowerQuery)
			if idx == -1 {
				break
			}
			e.matches = append(e.matches, Match{
				Line: lineNum,
				Col:  col + idx,
				Len:  len(query),
			})
			col += idx + 1
		}
	}

	if len(e.matches) > 0 {
		e.current = 0
	}

	return e.matches
}

// Next moves to next match
func (e *Engine) Next() *Match {
	if len(e.matches) == 0 {
		return nil
	}

	e.current = (e.current + 1) % len(e.matches)
	return &e.matches[e.current]
}

// Previous moves to previous match
func (e *Engine) Previous() *Match {
	if len(e.matches) == 0 {
		return nil
	}

	e.current--
	if e.current < 0 {
		e.current = len(e.matches) - 1
	}
	return &e.matches[e.current]
}

// Current returns current match
func (e *Engine) Current() *Match {
	if e.current >= 0 && e.current < len(e.matches) {
		return &e.matches[e.current]
	}
	return nil
}

// Clear clears search results
func (e *Engine) Clear() {
	e.query = ""
	e.matches = []Match{}
	e.current = -1
}

// HasMatches returns true if there are matches
func (e *Engine) HasMatches() bool {
	return len(e.matches) > 0
}

// MatchCount returns number of matches
func (e *Engine) MatchCount() int {
	return len(e.matches)
}

// CurrentIndex returns current match index (1-based)
func (e *Engine) CurrentIndex() int {
	if e.current >= 0 {
		return e.current + 1
	}
	return 0
}

// Query returns current search query
func (e *Engine) Query() string {
	return e.query
}

// GetMatches returns all matches
func (e *Engine) GetMatches() []Match {
	return e.matches
}
