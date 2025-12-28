package editor

import (
	"github.com/Adelodunpeter25/vx/internal/terminal"
	"github.com/gdamore/tcell/v2"
)

type Editor struct {
	term   *terminal.Terminal
	width  int
	height int
	quit   bool
}

func New(term *terminal.Terminal) *Editor {
	width, height := term.Size()
	return &Editor{
		term:   term,
		width:  width,
		height: height,
	}
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
	if ev.Rune == 'q' {
		e.quit = true
	}
}

func (e *Editor) render() {
	e.term.Clear()
	
	msg := "vx - press 'q' to quit"
	x := (e.width - len(msg)) / 2
	y := e.height / 2
	
	e.term.DrawText(x, y, msg, tcell.StyleDefault)
	e.term.Show()
}
