package editor

import (
	"github.com/Adelodunpeter25/vx/internal/buffer"
	"github.com/Adelodunpeter25/vx/internal/terminal"
)

type Editor struct {
	term       *terminal.Terminal
	buffer     *buffer.Buffer
	width      int
	height     int
	cursorX    int
	cursorY    int
	offsetY    int
	mode       Mode
	quit       bool
	commandBuf string
	message    string
}

func New(term *terminal.Terminal) *Editor {
	width, height := term.Size()
	return &Editor{
		term:    term,
		buffer:  buffer.New(),
		width:   width,
		height:  height,
		mode:    ModeNormal,
	}
}

func NewWithFile(term *terminal.Terminal, filename string) (*Editor, error) {
	buf, err := buffer.Load(filename)
	if err != nil {
		return nil, err
	}
	
	width, height := term.Size()
	return &Editor{
		term:    term,
		buffer:  buf,
		width:   width,
		height:  height,
		mode:    ModeNormal,
	}, nil
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
		e.width, e.height = e.term.Size()
		e.render()
	}
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
	e.render()
}
