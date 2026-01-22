package editor

import (
	"fmt"

	"github.com/Adelodunpeter25/vx/internal/terminal"
	"github.com/gdamore/tcell/v2"
)

func (e *Editor) handleSearchMode(ev *terminal.Event) {
	if ev.Key == tcell.KeyEscape {
		e.mode = ModeNormal
		e.searchBuf = ""
		e.search.Clear()
		e.msgManager.Clear()
		return
	}
	
	if ev.Key == tcell.KeyEnter {
		// Just exit search mode, results already visible
		if e.search.HasMatches() {
			e.msgManager.SetPersistent(fmt.Sprintf("/%s [%d/%d]", e.searchBuf, e.search.CurrentIndex(), e.search.MatchCount()))
		}
		e.mode = ModeNormal
		return
	}
	
	if ev.Key == tcell.KeyBackspace || ev.Key == tcell.KeyBackspace2 {
		if len(e.searchBuf) > 0 {
			e.searchBuf = e.searchBuf[:len(e.searchBuf)-1]
			// Update search in real-time
			e.performIncrementalSearch()
		}
		return
	}
	
	if ev.Rune != 0 {
		e.searchBuf += string(ev.Rune)
		// Update search in real-time
		e.performIncrementalSearch()
	}
}

func (e *Editor) performIncrementalSearch() {
	if e.searchBuf == "" {
		e.search.Clear()
		e.msgManager.Clear()
		return
	}
	
	// Get all lines from buffer
	lines := make([]string, e.buffer.LineCount())
	for i := 0; i < e.buffer.LineCount(); i++ {
		lines[i] = e.buffer.Line(i)
	}
	
	// Perform search
	matches := e.search.Search(lines, e.searchBuf)
	
	if len(matches) == 0 {
		e.msgManager.SetPersistent(fmt.Sprintf("Pattern not found: %s", e.searchBuf))
		return
	}
	
	// Jump to first match
	match := e.search.Current()
	if match != nil {
		e.cursorY = match.Line
		e.cursorX = match.Col
		e.adjustScroll()
	}
	
	e.msgManager.SetPersistent(fmt.Sprintf("/%s [%d/%d]", e.searchBuf, e.search.CurrentIndex(), e.search.MatchCount()))
}

func (e *Editor) performSearch() {
	if e.searchBuf == "" {
		e.mode = ModeNormal
		e.search.Clear()
		e.msgManager.Clear()
		return
	}
	
	// Get all lines from buffer
	lines := make([]string, e.buffer.LineCount())
	for i := 0; i < e.buffer.LineCount(); i++ {
		lines[i] = e.buffer.Line(i)
	}
	
	// Perform search
	matches := e.search.Search(lines, e.searchBuf)
	
	if len(matches) == 0 {
		e.msgManager.SetPersistent(fmt.Sprintf("Pattern not found: %s", e.searchBuf))
		e.mode = ModeNormal
		return
	}
	
	// Jump to first match
	match := e.search.Current()
	if match != nil {
		e.cursorY = match.Line
		e.cursorX = match.Col
		e.adjustScroll()
	}
	
	e.msgManager.SetPersistent(fmt.Sprintf("/%s [%d/%d]", e.searchBuf, e.search.CurrentIndex(), e.search.MatchCount()))
	e.mode = ModeNormal
}
