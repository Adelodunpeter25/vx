package terminal

import "github.com/gdamore/tcell/v2"

type EventType int

const (
	EventKey EventType = iota
	EventResize
	EventMouse
	EventQuit
)

type Event struct {
	Type   EventType
	Key    tcell.Key
	Rune   rune
	Button tcell.ButtonMask
	MouseX int
	MouseY int
}

func (t *Terminal) ReadEvent() *Event {
	ev := t.screen.PollEvent()
	switch ev := ev.(type) {
	case *tcell.EventKey:
		return &Event{
			Type: EventKey,
			Key:  ev.Key(),
			Rune: ev.Rune(),
		}
	case *tcell.EventResize:
		t.screen.Sync()
		return &Event{Type: EventResize}
	case *tcell.EventMouse:
		x, y := ev.Position()
		return &Event{
			Type:   EventMouse,
			Button: ev.Buttons(),
			MouseX: x,
			MouseY: y,
		}
	}
	return nil
}
