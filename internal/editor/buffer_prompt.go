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
				p.bufferMgr.Delete()
				e.switchToBuffer()
				p.msgManager.SetTransient("Buffer saved and closed")
				p.mode = ModeNormal
			}
			p.renderCache.invalidate()

		case 'n', 'N':
			// Close without saving
			p.bufferMgr.Delete()
			e.switchToBuffer()
			p.msgManager.SetTransient("Buffer closed without saving")
			p.mode = ModeNormal
			p.renderCache.invalidate()
		}
	}
}
