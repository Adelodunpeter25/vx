package editor

import (
	"github.com/Adelodunpeter25/vx/internal/command"
	"github.com/gdamore/tcell/v2"
)

func (e *Editor) handleCommandMode(ev *terminal.Event) {
	if ev.Key == tcell.KeyEscape {
		e.mode = ModeNormal
		e.commandBuf = ""
		e.message = ""
		return
	}
	
	if ev.Key == tcell.KeyEnter {
		result := command.Execute(e.commandBuf, e.buffer)
		if result.Error != nil {
			e.message = result.Error.Error()
		} else if result.Message != "" {
			e.message = result.Message
		}
		if result.Quit {
			e.quit = true
		}
		e.mode = ModeNormal
		e.commandBuf = ""
		return
	}
	
	if ev.Key == tcell.KeyBackspace || ev.Key == tcell.KeyBackspace2 {
		if len(e.commandBuf) > 0 {
			e.commandBuf = e.commandBuf[:len(e.commandBuf)-1]
		}
		return
	}
	
	if ev.Rune != 0 {
		e.commandBuf += string(ev.Rune)
	}
}
