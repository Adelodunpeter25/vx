package editor

import (
	"github.com/Adelodunpeter25/vx/internal/replace"
	"github.com/gdamore/tcell/v2"
)

func (e *Editor) handleReplaceMode(ev *tcell.EventKey) {
	p := e.active()
	state := p.replace.GetState()

	switch ev.Key() {
	case tcell.KeyEscape:
		p.replace.Cancel()
		p.mode = ModeNormal
		p.msgManager.Clear()
		p.renderCache.invalidate()
		return

	case tcell.KeyEnter:
		if state == replace.StateSearchInput {
			// Perform search using existing search engine
			lines := make([]string, p.buffer.LineCount())
			for i := 0; i < p.buffer.LineCount(); i++ {
				lines[i] = p.buffer.Line(i)
			}
			matches := p.search.Search(lines, p.replace.GetSearchTerm())
			p.replace.ConfirmSearch(matches)
			if len(matches) == 0 {
				p.msgManager.SetTransient("No matches found")
				p.mode = ModeNormal
				p.replace.Cancel()
			}
			p.renderCache.invalidate()
		} else if state == replace.StateReplaceInput {
			// Start confirmation
			p.replace.ConfirmReplace()
			match := p.replace.GetCurrentMatch()
			if match != nil {
				p.cursorY = match.Line
				p.cursorX = match.Col
				e.adjustScroll()
			}
			p.renderCache.invalidate()
		}
		return

	case tcell.KeyBackspace, tcell.KeyBackspace2:
		if state == replace.StateSearchInput {
			p.replace.BackspaceSearch()
		} else if state == replace.StateReplaceInput {
			p.replace.BackspaceReplace()
		}
		p.renderCache.invalidate()
		return

	case tcell.KeyRune:
		r := ev.Rune()

		if state == replace.StateConfirm {
			// Handle y/n/q during confirmation
			switch r {
			case 'y':
				// Replace current match
				match := p.replace.GetCurrentMatch()
				if match != nil {
					// Delete old text and insert new text
					searchLen := len(p.replace.GetSearchTerm())

					// Delete characters one by one from the end
					for i := 0; i < searchLen; i++ {
						p.buffer.DeleteRune(match.Line, match.Col+searchLen-i)
					}

					// Insert replacement text character by character
					replaceTerm := p.replace.GetReplaceTerm()
					for i, r := range replaceTerm {
						p.buffer.InsertRune(match.Line, match.Col+i, r)
					}
				}
				// Move to next match
				if !p.replace.NextMatch() {
					p.msgManager.SetTransient("Replace complete")
					p.mode = ModeNormal
					p.search.Clear() // Clear search highlights
				} else {
					match = p.replace.GetCurrentMatch()
					if match != nil {
						p.cursorY = match.Line
						p.cursorX = match.Col
						e.adjustScroll()
					}
				}
				p.renderCache.invalidate()

			case 'n':
				// Skip to next match
				if !p.replace.NextMatch() {
					p.msgManager.SetTransient("Replace complete")
					p.mode = ModeNormal
					p.search.Clear() // Clear search highlights
				} else {
					match := p.replace.GetCurrentMatch()
					if match != nil {
						p.cursorY = match.Line
						p.cursorX = match.Col
						e.adjustScroll()
					}
				}
				p.renderCache.invalidate()

			case 'q':
				// Quit replace
				p.replace.Cancel()
				p.mode = ModeNormal
				p.search.Clear() // Clear search highlights
				p.msgManager.SetTransient("Replace cancelled")
				p.renderCache.invalidate()
			}
		} else if state == replace.StateSearchInput {
			p.replace.AppendToSearch(r)
			p.renderCache.invalidate()
		} else if state == replace.StateReplaceInput {
			p.replace.AppendToReplace(r)
			p.renderCache.invalidate()
		}
	}
}
