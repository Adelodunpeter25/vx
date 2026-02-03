package editor

import (
	"github.com/Adelodunpeter25/vx/internal/buffer"
	splitpane "github.com/Adelodunpeter25/vx/internal/split-pane"
	"github.com/Adelodunpeter25/vx/internal/terminal"
	"github.com/Adelodunpeter25/vx/internal/utils"
	"github.com/gdamore/tcell/v2"
)

type Editor struct {
	term       *terminal.Terminal
	width      int
	height     int
	panes      []*Pane
	activePane int
	splitRatio float64
	dragSplit  bool
	quit       bool
}

func New(term *terminal.Terminal) *Editor {
	width, height := term.Size()
	buf := buffer.New()
	pane := NewPane(buf, "")
	return &Editor{
		term:       term,
		width:      width,
		height:     height,
		panes:      []*Pane{pane},
		activePane: 0,
		splitRatio: 0.5,
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
				term:       term,
				width:      width,
				height:     height,
				panes:      []*Pane{pane},
				activePane: 0,
				splitRatio: 0.5,
			}
			return ed, nil
		}
		return nil, err
	}

	width, height := term.Size()
	pane := NewPane(buf, filename)
	ed := &Editor{
		term:       term,
		width:      width,
		height:     height,
		panes:      []*Pane{pane},
		activePane: 0,
		splitRatio: 0.5,
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
	rects, dividerX := splitpane.LayoutSideBySide(e.width, contentHeight, len(e.panes), e.splitRatio)
	if len(rects) >= 2 && dividerX >= 0 {
		if e.handleSplitterDrag(ev, dividerX) {
			return
		}
	}
	for i, rect := range rects {
		if ev.MouseX >= rect.X && ev.MouseX < rect.X+rect.Width && ev.MouseY >= rect.Y && ev.MouseY < rect.Y+rect.Height {
			// Focus only on click or scroll, not on hover/move.
			if ev.Button == tcell.Button1 || ev.Button == tcell.WheelUp || ev.Button == tcell.WheelDown {
				e.activePane = i
			}
			local := *ev
			local.MouseX = ev.MouseX - rect.X
			local.MouseY = ev.MouseY - rect.Y
			e.handleMouseEvent(&local)
			return
		}
	}
}

func (e *Editor) handleSplitterDrag(ev *terminal.Event, dividerX int) bool {
	if ev.Button == tcell.Button1 && abs(ev.MouseX-dividerX) <= 1 {
		e.dragSplit = true
	}
	if ev.Button == tcell.ButtonNone && e.dragSplit {
		e.dragSplit = false
		return true
	}
	if e.dragSplit {
		width := e.width
		if width <= 0 {
			return true
		}
		minLeft := 10
		minRight := 10
		maxLeft := width - 1 - minRight
		if maxLeft < minLeft {
			maxLeft = minLeft
		}
		left := ev.MouseX
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

func (e *Editor) ensurePaneCount() {
	if len(e.panes) == 0 {
		buf := buffer.New()
		pane := NewPane(buf, "")
		e.panes = []*Pane{pane}
		e.activePane = 0
	}
}
