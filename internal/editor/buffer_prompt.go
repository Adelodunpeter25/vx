package editor

import "github.com/gdamore/tcell/v2"

func (e *Editor) handleBufferPromptMode(ev *tcell.EventKey) {
	switch ev.Key() {
	case tcell.KeyEscape:
		e.mode = ModeNormal
		e.message = ""
		e.renderCache.invalidate()
		return

	case tcell.KeyRune:
		r := ev.Rune()
		switch r {
		case 'y', 'Y':
			// Save and close buffer
			if err := e.buffer.Save(); err != nil {
				e.message = "Error saving: " + err.Error()
				e.mode = ModeNormal
			} else {
				e.bufferMgr.Delete()
				e.switchToBuffer()
				e.message = "Buffer saved and closed"
				e.mode = ModeNormal
			}
			e.renderCache.invalidate()

		case 'n', 'N':
			// Close without saving
			e.bufferMgr.Delete()
			e.switchToBuffer()
			e.message = "Buffer closed without saving"
			e.mode = ModeNormal
			e.renderCache.invalidate()
		}
	}
}
