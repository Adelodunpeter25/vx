package editor

import (
	"github.com/Adelodunpeter25/vx/internal/replace"
	"github.com/gdamore/tcell/v2"
)

func (e *Editor) handleReplaceMode(ev *tcell.EventKey) {
	state := e.replace.GetState()

	switch ev.Key() {
	case tcell.KeyEscape:
		e.replace.Cancel()
		e.mode = ModeNormal
		e.message = ""
		e.renderCache.invalidate()
		return

	case tcell.KeyEnter:
		if state == replace.StateSearchInput {
			// Perform search using existing search engine
			lines := make([]string, e.buffer.LineCount())
			for i := 0; i < e.buffer.LineCount(); i++ {
				lines[i] = e.buffer.Line(i)
			}
			matches := e.search.Search(lines, e.replace.GetSearchTerm())
			e.replace.ConfirmSearch(matches)
			if len(matches) == 0 {
				e.message = "No matches found"
				e.mode = ModeNormal
				e.replace.Cancel()
			}
			e.renderCache.invalidate()
		} else if state == replace.StateReplaceInput {
			// Start confirmation
			e.replace.ConfirmReplace()
			match := e.replace.GetCurrentMatch()
			if match != nil {
				e.cursorY = match.Line
				e.cursorX = match.Col
				e.adjustScroll()
			}
			e.renderCache.invalidate()
		}
		return

	case tcell.KeyBackspace, tcell.KeyBackspace2:
		if state == replace.StateSearchInput {
			e.replace.BackspaceSearch()
		} else if state == replace.StateReplaceInput {
			e.replace.BackspaceReplace()
		}
		e.renderCache.invalidate()
		return

	case tcell.KeyRune:
		r := ev.Rune()

		if state == replace.StateConfirm {
			// Handle y/n/q during confirmation
			switch r {
			case 'y':
				// Replace current match
				match := e.replace.GetCurrentMatch()
				if match != nil {
					// Delete old text and insert new text
					searchLen := len(e.replace.GetSearchTerm())
					
					// Delete characters one by one from the end
					for i := 0; i < searchLen; i++ {
						e.buffer.DeleteRune(match.Line, match.Col+searchLen-i)
					}
					
					// Insert replacement text character by character
					replaceTerm := e.replace.GetReplaceTerm()
					for i, r := range replaceTerm {
						e.buffer.InsertRune(match.Line, match.Col+i, r)
					}
				}
				// Move to next match
				if !e.replace.NextMatch() {
					e.message = "Replace complete"
					e.mode = ModeNormal
				} else {
					match = e.replace.GetCurrentMatch()
					if match != nil {
						e.cursorY = match.Line
						e.cursorX = match.Col
						e.adjustScroll()
					}
				}
				e.renderCache.invalidate()

			case 'n':
				// Skip to next match
				if !e.replace.NextMatch() {
					e.message = "Replace complete"
					e.mode = ModeNormal
				} else {
					match := e.replace.GetCurrentMatch()
					if match != nil {
						e.cursorY = match.Line
						e.cursorX = match.Col
						e.adjustScroll()
					}
				}
				e.renderCache.invalidate()

			case 'q':
				// Quit replace
				e.replace.Cancel()
				e.mode = ModeNormal
				e.message = "Replace cancelled"
				e.renderCache.invalidate()
			}
		} else if state == replace.StateSearchInput {
			e.replace.AppendToSearch(r)
			e.renderCache.invalidate()
		} else if state == replace.StateReplaceInput {
			e.replace.AppendToReplace(r)
			e.renderCache.invalidate()
		}
	}
}
