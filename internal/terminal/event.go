package terminal

import (
	"time"

	"github.com/gdamore/tcell/v2"
)

type EventType int

const (
	EventKey EventType = iota
	EventResize
	EventMouse
	EventQuit
	EventFileChange
)

type Event struct {
	Type       EventType
	Key        tcell.Key
	Rune       rune
	Button     tcell.ButtonMask
	MouseX     int
	MouseY     int
	ChangePath string // populated for EventFileChange
}

// FileChangeEvent is posted into the tcell event loop by the file watcher goroutine.
type FileChangeEvent struct {
	path string
	when time.Time
}

func NewFileChangeEvent(path string) *FileChangeEvent {
	return &FileChangeEvent{path: path, when: time.Now()}
}

func (e *FileChangeEvent) When() time.Time { return e.when }
func (e *FileChangeEvent) Path() string    { return e.path }

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
	case *FileChangeEvent:
		return &Event{
			Type:       EventFileChange,
			ChangePath: ev.Path(),
		}
	}
	return nil
}
