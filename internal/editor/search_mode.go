package editor

import (
	"fmt"

	"github.com/Adelodunpeter25/vx/internal/terminal"
	"github.com/gdamore/tcell/v2"
)

func (e *Editor) handleSearchMode(ev *terminal.Event) {
	p := e.active()
	if ev.Key == tcell.KeyEscape {
		p.mode = ModeNormal
		p.searchBuf = ""
		p.search.Clear()
		p.msgManager.Clear()
		return
	}

	if ev.Key == tcell.KeyEnter {
		// Just exit search mode, results already visible
		if p.search.HasMatches() {
			p.msgManager.SetPersistent(fmt.Sprintf("/%s [%d/%d]", p.searchBuf, p.search.CurrentIndex(), p.search.MatchCount()))
		}
		p.mode = ModeNormal
		return
	}

	if ev.Key == tcell.KeyBackspace || ev.Key == tcell.KeyBackspace2 {
		if len(p.searchBuf) > 0 {
			p.searchBuf = p.searchBuf[:len(p.searchBuf)-1]
			// Update search in real-time
			e.performIncrementalSearch()
		}
		return
	}

	if ev.Rune != 0 {
		p.searchBuf += string(ev.Rune)
		// Update search in real-time
		e.performIncrementalSearch()
	}
}

func (e *Editor) performIncrementalSearch() {
	p := e.active()
	if p.searchBuf == "" {
		p.search.Clear()
		p.msgManager.Clear()
		return
	}

	// Get all lines from buffer
	lines := make([]string, p.buffer.LineCount())
	for i := 0; i < p.buffer.LineCount(); i++ {
		lines[i] = p.buffer.Line(i)
	}

	// Perform search
	matches := p.search.Search(lines, p.searchBuf)

	if len(matches) == 0 {
		p.msgManager.SetPersistent(fmt.Sprintf("Pattern not found: %s", p.searchBuf))
		return
	}

	// Jump to first match
	match := p.search.Current()
	if match != nil {
		p.cursorY = match.Line
		p.cursorX = match.Col
		e.adjustScroll()
	}

	p.msgManager.SetPersistent(fmt.Sprintf("/%s [%d/%d]", p.searchBuf, p.search.CurrentIndex(), p.search.MatchCount()))
}

func (e *Editor) performSearch() {
	p := e.active()
	if p.searchBuf == "" {
		p.mode = ModeNormal
		p.search.Clear()
		p.msgManager.Clear()
		return
	}

	// Get all lines from buffer
	lines := make([]string, p.buffer.LineCount())
	for i := 0; i < p.buffer.LineCount(); i++ {
		lines[i] = p.buffer.Line(i)
	}

	// Perform search
	matches := p.search.Search(lines, p.searchBuf)

	if len(matches) == 0 {
		p.msgManager.SetPersistent(fmt.Sprintf("Pattern not found: %s", p.searchBuf))
		p.mode = ModeNormal
		return
	}

	// Jump to first match
	match := p.search.Current()
	if match != nil {
		p.cursorY = match.Line
		p.cursorX = match.Col
		e.adjustScroll()
	}

	p.msgManager.SetPersistent(fmt.Sprintf("/%s [%d/%d]", p.searchBuf, p.search.CurrentIndex(), p.search.MatchCount()))
	p.mode = ModeNormal
}
