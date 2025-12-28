package editor

import (
	"github.com/Adelodunpeter25/vx/internal/buffer"
	"github.com/Adelodunpeter25/vx/internal/syntax"
	"github.com/Adelodunpeter25/vx/internal/terminal"
)

type Editor struct {
	term        *terminal.Terminal
	buffer      *buffer.Buffer
	syntax      *syntax.Engine
	renderCache *RenderCache
	width       int
	height      int
	cursorX     int
	cursorY     int
	offsetY     int
	mode        Mode
	quit        bool
	commandBuf  string
	message     string
}

func New(term *terminal.Terminal) *Editor {
	width, height := term.Size()
	return &Editor{
		term:        term,
		buffer:      buffer.New(),
		syntax:      syntax.New(""),
		renderCache: newRenderCache(),
		width:       width,
		height:      height,
		mode:        ModeNormal,
	}
}

func NewWithFile(term *terminal.Terminal, filename string) (*Editor, error) {
	buf, err := buffer.Load(filename)
	if err != nil {
		// Check if it's a partial load (recoverable error)
		if buf != nil {
			// File loaded with warnings
			width, height := term.Size()
			ed := &Editor{
				term:        term,
				buffer:      buf,
				syntax:      syntax.New(filename),
				renderCache: newRenderCache(),
				width:       width,
				height:      height,
				mode:        ModeNormal,
				message:     "Warning: " + err.Error(),
			}
			return ed, nil
		}
		return nil, err
	}
	
	width, height := term.Size()
	ed := &Editor{
		term:        term,
		buffer:      buf,
		syntax:      syntax.New(filename),
		renderCache: newRenderCache(),
		width:       width,
		height:      height,
		mode:        ModeNormal,
	}
	
	// Show warning if syntax highlighting disabled due to size
	if ed.syntax.IsTooLarge() {
		ed.message = "File too large for syntax highlighting"
	}
	
	return ed, nil
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
	}
}

func (e *Editor) handleResize() {
	e.width, e.height = e.term.Size()
	
	// Ensure cursor stays visible after resize
	e.clampCursor()
	e.adjustScroll()
	
	e.renderCache.invalidate()
	e.render()
}

func (e *Editor) handleKey(ev *terminal.Event) {
	switch e.mode {
	case ModeNormal:
		e.handleNormalMode(ev)
	case ModeInsert:
		e.handleInsertMode(ev)
	case ModeCommand:
		e.handleCommandMode(ev)
	}
	e.renderCache.invalidate()
	e.render()
}
