package editor

import (
	"os"

	"github.com/Adelodunpeter25/vx/internal/buffer"
	filebrowser "github.com/Adelodunpeter25/vx/internal/file-browser"
	"github.com/Adelodunpeter25/vx/internal/syntax"
	"github.com/Adelodunpeter25/vx/internal/terminal"
	"github.com/gdamore/tcell/v2"
)

func (e *Editor) openFileInActivePane(path string) {
	p := e.active()
	if p.buffer.IsModified() {
		p.msgManager.SetError("No write since last change (use :e! to override)")
		return
	}
	newBuf, err := buffer.Load(path)
	if err != nil {
		p.msgManager.SetError("Error: " + err.Error())
		return
	}
	p.buffer = newBuf
	p.syntax = syntax.New(newBuf.Filename())
	p.cursorX = 0
	p.cursorY = 0
	p.offsetY = 0
	p.renderCache.invalidate()
	e.showFileInfo()
}

func (e *Editor) previewFileInActivePane(path string) {
	p := e.active()
	if p.buffer.IsModified() {
		p.msgManager.SetError("No write since last change (use :e! to override)")
		return
	}
	newBuf, err := buffer.Load(path)
	if err != nil {
		p.msgManager.SetError("Error: " + err.Error())
		return
	}
	p.buffer = newBuf
	p.syntax = syntax.New(newBuf.Filename())
	p.cursorX = 0
	p.cursorY = 0
	p.offsetY = 0
	p.renderCache.invalidate()
	e.showFileInfo()
}

func (e *Editor) toggleFileBrowser() {
	if e.fileBrowser == nil {
		e.fileBrowser = filebrowser.New("")
	}
	e.fileBrowser.Open = !e.fileBrowser.Open
	if e.fileBrowser.Open {
		e.fileBrowser.Focused = true
	} else {
		e.fileBrowser.Focused = false
	}
}

func (e *Editor) handleBrowserResizeDrag(ev *terminal.Event, dividerX int) bool {
	if ev.Button == tcell.Button1 && abs(ev.MouseX-dividerX) <= 1 {
		e.dragBrowser = true
	}
	if ev.Button == tcell.ButtonNone && e.dragBrowser {
		e.dragBrowser = false
		return true
	}
	if e.dragBrowser {
		minWidth := 10
		maxWidth := e.width - 10
		if maxWidth < minWidth {
			maxWidth = minWidth
		}
		newWidth := ev.MouseX
		if newWidth < minWidth {
			newWidth = minWidth
		}
		if newWidth > maxWidth {
			newWidth = maxWidth
		}
		if e.fileBrowser != nil {
			e.fileBrowser.Width = newWidth
		}
		e.active().renderCache.invalidate()
		return true
	}
	return false
}

func (e *Editor) startCdPrompt(initial string) {
	if e.cdPrompt == nil {
		e.cdPrompt = filebrowser.NewCdPrompt(initial)
	} else {
		if initial == "" {
			e.cdPrompt = filebrowser.NewCdPrompt("")
		} else {
			e.cdPrompt = filebrowser.NewCdPrompt(initial)
		}
	}
	e.active().mode = ModeCdPrompt
}

func (e *Editor) handleCdPrompt(ev *terminal.Event) {
	if e.cdPrompt == nil {
		e.cdPrompt = filebrowser.NewCdPrompt("")
	}
	action := e.cdPrompt.HandleKey(ev)
	if action.Cancel {
		e.active().mode = ModeNormal
		return
	}
	if action.Apply {
		path := filebrowser.ExpandHome(action.Path)
		if path == "" {
			path = "."
		}
		if err := os.Chdir(path); err != nil {
			e.active().msgManager.SetError("Error: " + err.Error())
		} else {
			if e.fileBrowser != nil {
				e.fileBrowser.SetRoot(path)
			}
			e.active().msgManager.SetTransient("Changed directory to " + abbreviateHome(path))
		}
		e.active().mode = ModeNormal
	}
}
