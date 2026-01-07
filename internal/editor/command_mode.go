package editor

import (
	"strings"

	"github.com/Adelodunpeter25/vx/internal/command"
	"github.com/Adelodunpeter25/vx/internal/syntax"
	"github.com/Adelodunpeter25/vx/internal/terminal"
	"github.com/Adelodunpeter25/vx/internal/utils"
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
		// Show "Saving..." for write commands
		if strings.HasPrefix(e.commandBuf, "w") {
			e.message = "Saving..."
			e.render()
		}
		
		result := command.Execute(e.commandBuf, e.buffer)
		
		// Handle buffer operations
		if result.AddBuffer && result.NewBuffer != nil {
			e.addBuffer(result.NewBuffer, result.NewBuffer.Filename())
		} else if result.DeleteBuffer {
			e.deleteCurrentBuffer()
			// If we're in prompt mode, don't reset to normal
			if e.mode == ModeBufferPrompt {
				e.commandBuf = ""
				return
			}
		} else if result.SwitchFile && result.NewBuffer != nil {
			// Handle file switching (replace current buffer)
			e.buffer = result.NewBuffer
			e.syntax = syntax.New(result.NewBuffer.Filename())
			e.cursorX = 0
			e.cursorY = 0
			e.offsetY = 0
			e.renderCache.invalidate()
		}
		
		if result.Error != nil {
			e.message = utils.FormatUserError(result.Error)
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
