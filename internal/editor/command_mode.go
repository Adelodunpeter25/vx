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
	p := e.active()
	if ev.Key == tcell.KeyEscape {
		p.mode = ModeNormal
		p.commandBuf = ""
		p.msgManager.Clear()
		return
	}

	if ev.Key == tcell.KeyEnter {
		// Show "Saving..." for write commands
		if strings.HasPrefix(p.commandBuf, "w") {
			p.msgManager.SetPersistent("Saving...")
			e.render()
		}

		result := command.Execute(p.commandBuf, p.buffer)

		// Handle buffer operations
		if result.AddBuffer && result.NewBuffer != nil {
			e.addPaneWithBuffer(result.NewBuffer, result.NewBuffer.Filename())
		} else if result.DeleteBuffer {
			e.deleteCurrentBuffer()
			// If we're in prompt mode, don't reset to normal
			if p.mode == ModeBufferPrompt {
				p.commandBuf = ""
				return
			}
		} else if result.SwitchFile && result.NewBuffer != nil {
			// Handle file switching (replace current buffer)
			p.buffer = result.NewBuffer
			p.syntax = syntax.New(result.NewBuffer.Filename())
			p.cursorX = 0
			p.cursorY = 0
			p.offsetY = 0
			p.renderCache.invalidate()
		}

		if result.Error != nil {
			p.msgManager.SetError(utils.FormatUserError(result.Error))
		} else if result.Message != "" {
			p.msgManager.SetPersistent(result.Message)
		}
		if result.Quit {
			e.quit = true
		}
		p.mode = ModeNormal
		p.commandBuf = ""
		return
	}

	if ev.Key == tcell.KeyBackspace || ev.Key == tcell.KeyBackspace2 {
		if len(p.commandBuf) > 0 {
			p.commandBuf = p.commandBuf[:len(p.commandBuf)-1]
		}
		return
	}

	if ev.Rune != 0 {
		p.commandBuf += string(ev.Rune)
	}
}
