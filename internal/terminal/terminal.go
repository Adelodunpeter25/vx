package terminal

import (
	"github.com/gdamore/tcell/v2"
)

type Terminal struct {
	screen tcell.Screen
}

func New() (*Terminal, error) {
	screen, err := tcell.NewScreen()
	if err != nil {
		return nil, err
	}
	if err := screen.Init(); err != nil {
		return nil, err
	}
	screen.EnableMouse()
	screen.Clear()
	return &Terminal{screen: screen}, nil
}

func (t *Terminal) Close() {
	t.screen.Fini()
}

func (t *Terminal) Size() (width, height int) {
	return t.screen.Size()
}

func (t *Terminal) Clear() {
	t.screen.Clear()
}

func (t *Terminal) Show() {
	t.screen.Show()
}

func (t *Terminal) PollEvent() tcell.Event {
	return t.screen.PollEvent()
}
