package editor

import (
	"github.com/Adelodunpeter25/vx/internal/utils"
	"github.com/gdamore/tcell/v2"
)

func (e *Editor) handleBufferPromptMode(ev *tcell.EventKey) {
	switch ev.Key() {
	case tcell.KeyEscape:
		e.mode = ModeNormal
		e.msgManager.Clear()
		e.renderCache.invalidate()
		return

	case tcell.KeyRune:
		r := ev.Rune()
		switch r {
		case 'y', 'Y':
			// Save and close buffer
			if err := e.buffer.Save(); err != nil {
				e.msgManager.SetError(utils.FormatSaveError(e.buffer.Filename(), err))
				e.mode = ModeNormal
			} else {
				e.bufferMgr.Delete()
				e.switchToBuffer()
				e.msgManager.SetTransient("Buffer saved and closed")
				e.mode = ModeNormal
			}
			e.renderCache.invalidate()

		case 'n', 'N':
			// Close without saving
			e.bufferMgr.Delete()
			e.switchToBuffer()
			e.msgManager.SetTransient("Buffer closed without saving")
			e.mode = ModeNormal
			e.renderCache.invalidate()
		}
	}
}
