package editor

import (
	"github.com/Adelodunpeter25/vx/internal/utils"
	"github.com/gdamore/tcell/v2"
)

func (e *Editor) handleBufferPromptMode(ev *tcell.EventKey) {
	p := e.active()
	switch ev.Key() {
	case tcell.KeyEscape:
		p.mode = ModeNormal
		p.msgManager.Clear()
		p.renderCache.invalidate()
		return

	case tcell.KeyRune:
		r := ev.Rune()
		switch r {
		case 'y', 'Y':
			// Save and close buffer
			if err := p.buffer.Save(); err != nil {
				p.msgManager.SetError(utils.FormatSaveError(p.buffer.Filename(), err))
				p.mode = ModeNormal
			} else {
				e.panes = append(e.panes[:e.activePane], e.panes[e.activePane+1:]...)
				if e.activePane >= len(e.panes) {
					e.activePane = len(e.panes) - 1
				}
				if e.active() != nil {
					e.active().msgManager.SetTransient("Pane saved and closed")
				}
				p.mode = ModeNormal
			}
			p.renderCache.invalidate()

		case 'n', 'N':
			// Close without saving
			e.panes = append(e.panes[:e.activePane], e.panes[e.activePane+1:]...)
			if e.activePane >= len(e.panes) {
				e.activePane = len(e.panes) - 1
			}
			if e.active() != nil {
				e.active().msgManager.SetTransient("Pane closed without saving")
			}
			p.mode = ModeNormal
			p.renderCache.invalidate()
		}
	}
}
