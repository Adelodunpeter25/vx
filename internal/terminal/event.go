package terminal

import "github.com/gdamore/tcell/v2"

type EventType int

const (
	EventKey EventType = iota
	EventResize
	EventQuit
)

type Event struct {
	Type EventType
	Key  tcell.Key
	Rune rune
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
	}
	return nil
}
