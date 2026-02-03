package editor

import (
	"github.com/Adelodunpeter25/vx/internal/buffer"
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

func (e *Editor) Run() error {
	e.render()

	for !e.quit {
		e.handleEvent()
	}
	return nil
}

func (e *Editor) handleEvent() {
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
		e.handleMouseEvent(ev)
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
