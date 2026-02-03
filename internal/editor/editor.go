package editor

import (
	"github.com/Adelodunpeter25/vx/internal/buffer"
	filebrowser "github.com/Adelodunpeter25/vx/internal/file-browser"
	splitpane "github.com/Adelodunpeter25/vx/internal/split-pane"
	"github.com/Adelodunpeter25/vx/internal/syntax"
	"github.com/Adelodunpeter25/vx/internal/terminal"
	"github.com/Adelodunpeter25/vx/internal/utils"
	"github.com/gdamore/tcell/v2"
)

type Editor struct {
	term        *terminal.Terminal
	width       int
	height      int
	panes       []*Pane
	activePane  int
	splitRatio  float64
	dragSplit   bool
	dragBrowser bool
	fileBrowser *filebrowser.State
	quit        bool
}

func New(term *terminal.Terminal) *Editor {
	width, height := term.Size()
	buf := buffer.New()
	pane := NewPane(buf, "")
	return &Editor{
		term:        term,
		width:       width,
		height:      height,
		panes:       []*Pane{pane},
		activePane:  0,
		splitRatio:  0.5,
		fileBrowser: filebrowser.New(""),
	}
}

func NewWithFile(term *terminal.Terminal, filename string) (*Editor, error) {
	buf, err := buffer.Load(filename)
	if err != nil {
		// Check if it's a partial load (recoverable error)
		if buf != nil {
			// File loaded with warnings
			width, height := term.Size()
			pane := NewPane(buf, filename)
			pane.msgManager.SetError("Warning: " + utils.FormatLoadError(filename, err))
			ed := &Editor{
				term:        term,
				width:       width,
				height:      height,
				panes:       []*Pane{pane},
				activePane:  0,
				splitRatio:  0.5,
				fileBrowser: filebrowser.New(""),
			}
			return ed, nil
		}
		return nil, err
	}

	width, height := term.Size()
	pane := NewPane(buf, filename)
	ed := &Editor{
		term:        term,
		width:       width,
		height:      height,
		panes:       []*Pane{pane},
		activePane:  0,
		splitRatio:  0.5,
		fileBrowser: filebrowser.New(""),
	}

	// Show file info message on load
	ed.showFileInfo()

	// Show warning if syntax highlighting disabled due to size
	if ed.active().syntax.IsTooLarge() {
		ed.active().msgManager.SetPersistent("File too large for syntax highlighting")
	}

	return ed, nil
}

func (e *Editor) showFileInfo() {
	p := e.active()
	size, err := p.buffer.GetFileSize()
	if err != nil {
		return
	}

	filename := p.buffer.Filename()
	if filename == "" {
		filename = "[No Name]"
	}

	p.msgManager.SetPersistent(utils.FormatFileInfo(filename, size, p.buffer.LineCount()))
}

func (e *Editor) active() *Pane {
	if e.activePane < 0 || e.activePane >= len(e.panes) {
		return nil
	}
	return e.panes[e.activePane]
}

func (e *Editor) addPaneWithBuffer(buf *buffer.Buffer, filename string) {
	pane := NewPane(buf, filename)
	e.panes = append(e.panes, pane)
	e.activePane = len(e.panes) - 1
}

func (e *Editor) nextPane() {
	if len(e.panes) <= 1 {
		return
	}
	e.activePane = (e.activePane + 1) % len(e.panes)
}

func (e *Editor) previousPane() {
	if len(e.panes) <= 1 {
		return
	}
	e.activePane = (e.activePane - 1 + len(e.panes)) % len(e.panes)
}

func (e *Editor) deleteCurrentPane() {
	p := e.active()
	if p == nil {
		return
	}
	if len(e.panes) == 1 {
		p.msgManager.SetTransient("Cannot close last pane")
		return
	}
	if p.buffer.IsModified() {
		p.mode = ModeBufferPrompt
		p.msgManager.SetPersistent("Save changes? [y/n]")
		p.renderCache.invalidate()
		return
	}
	e.panes = append(e.panes[:e.activePane], e.panes[e.activePane+1:]...)
	if e.activePane >= len(e.panes) {
		e.activePane = len(e.panes) - 1
	}
	e.active().msgManager.SetTransient("Pane closed")
}

func (e *Editor) handleMouseEventForPane(ev *terminal.Event) {
	if len(e.panes) == 0 {
		return
	}
	if ev.MouseY >= e.height-1 {
		return
	}
	contentHeight := e.height - 1
	contentX := 0
	contentWidth := e.width
	if e.fileBrowser != nil && e.fileBrowser.Open {
		fbWidth := e.fileBrowser.Width
		if fbWidth < 10 {
			fbWidth = 10
		}
		if fbWidth > e.width-10 {
			fbWidth = e.width - 10
		}
		if e.dragBrowser {
			if e.handleBrowserResizeDrag(ev, fbWidth) {
				return
			}
		}
		if ev.MouseX < fbWidth && ev.MouseY < contentHeight {
			if ev.Button == tcell.Button1 || ev.Button == tcell.WheelUp || ev.Button == tcell.WheelDown {
				e.fileBrowser.Focused = true
			}
			action := e.fileBrowser.HandleMouse(ev, 0, 0, fbWidth, contentHeight)
			if action.PreviewPath != "" {
				e.previewFileInActivePane(action.PreviewPath)
			}
			if action.OpenPath != "" {
				e.openFileInActivePane(action.OpenPath)
			}
			return
		}
		if e.handleBrowserResizeDrag(ev, fbWidth) {
			return
		}
		contentX = fbWidth
		contentWidth = e.width - contentX
	}
	rects, dividerX := splitpane.LayoutSideBySide(contentWidth, contentHeight, len(e.panes), e.splitRatio)
	if len(rects) >= 2 && dividerX >= 0 {
		if e.handleSplitterDrag(ev, contentX+dividerX, contentX, contentWidth) {
			return
		}
	}
	for i, rect := range rects {
		rect.X += contentX
		if ev.MouseX >= rect.X && ev.MouseX < rect.X+rect.Width && ev.MouseY >= rect.Y && ev.MouseY < rect.Y+rect.Height {
			// Focus only on click or scroll, not on hover/move.
			if ev.Button == tcell.Button1 || ev.Button == tcell.WheelUp || ev.Button == tcell.WheelDown {
				e.activePane = i
				if e.fileBrowser != nil {
					e.fileBrowser.Focused = false
				}
			}
			local := *ev
			local.MouseX = ev.MouseX - rect.X
			local.MouseY = ev.MouseY - rect.Y
			e.handleMouseEvent(&local)
			return
		}
	}
}

func (e *Editor) handleSplitterDrag(ev *terminal.Event, dividerX int, contentX int, contentWidth int) bool {
	if ev.Button == tcell.Button1 && abs(ev.MouseX-dividerX) <= 1 {
		e.dragSplit = true
	}
	if ev.Button == tcell.ButtonNone && e.dragSplit {
		e.dragSplit = false
		return true
	}
	if e.dragSplit {
		width := contentWidth
		if width <= 0 {
			return true
		}
		minLeft := 10
		minRight := 10
		maxLeft := width - 1 - minRight
		if maxLeft < minLeft {
			maxLeft = minLeft
		}
		left := ev.MouseX - contentX
		if left < minLeft {
			left = minLeft
		}
		if left > maxLeft {
			left = maxLeft
		}
		available := width - 1
		if available <= 0 {
			return true
		}
		e.splitRatio = float64(left) / float64(available)
		e.active().renderCache.invalidate()
		return true
	}
	return false
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

func (e *Editor) Run() error {
	e.render()

	for !e.quit {
		e.handleEvent()
	}
	return nil
}

func (e *Editor) handleEvent() {
	e.ensurePaneCount()
	ev := e.term.ReadEvent()
	if ev == nil {
		return
	}

	switch ev.Type {
	case terminal.EventKey:
		e.handleKey(ev)
	case terminal.EventResize:
		e.handleResize()
	case terminal.EventMouse:
		e.handleMouseEventForPane(ev)
		e.active().renderCache.invalidate()
		e.render()
	}
}

func (e *Editor) handleResize() {
	e.width, e.height = e.term.Size()

	// Ensure cursor stays visible after resize
	e.clampCursor()
	e.adjustScroll()

	e.active().renderCache.invalidate()
	e.render()
}

func (e *Editor) handleKey(ev *terminal.Event) {
	if e.fileBrowser != nil && e.fileBrowser.Open && e.fileBrowser.Focused {
		action := e.fileBrowser.HandleKey(ev)
		if action.PreviewPath != "" {
			e.previewFileInActivePane(action.PreviewPath)
		}
		if action.OpenPath != "" {
			e.openFileInActivePane(action.OpenPath)
		}
		e.active().renderCache.invalidate()
		e.render()
		return
	}
	switch e.active().mode {
	case ModeNormal:
		e.handleNormalMode(ev)
	case ModeInsert:
		e.handleInsertMode(ev)
	case ModeCommand:
		e.handleCommandMode(ev)
	case ModeSearch:
		e.handleSearchMode(ev)
	case ModeReplace:
		// Convert terminal.Event to tcell.EventKey for replace mode
		tcellEv := tcell.NewEventKey(ev.Key, ev.Rune, tcell.ModNone)
		e.handleReplaceMode(tcellEv)
	case ModeBufferPrompt:
		tcellEv := tcell.NewEventKey(ev.Key, ev.Rune, tcell.ModNone)
		e.handleBufferPromptMode(tcellEv)
	}
	e.active().renderCache.invalidate()
	e.render()
}

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

func (e *Editor) ensurePaneCount() {
	if len(e.panes) == 0 {
		buf := buffer.New()
		pane := NewPane(buf, "")
		e.panes = []*Pane{pane}
		e.activePane = 0
	}
}
