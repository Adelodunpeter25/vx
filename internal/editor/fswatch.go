package editor

import (
	"time"

	"github.com/Adelodunpeter25/vx/internal/fswatch"
)

func (e *Editor) initWatcher(root string) {
	e.watcher = fswatch.New(root)
	e.watcher.Start()

	go func() {
		for {
			select {
			case <-e.watcher.Events:
				// Coalesce events (debounce)
				time.Sleep(100 * time.Millisecond)
				// Drain pending events
				for len(e.watcher.Events) > 0 {
					<-e.watcher.Events
				}

				if e.fileBrowser != nil {
					e.fileBrowser.Refresh()
					e.render()
				}
			case <-e.watcher.Errors:
				// Silent errors for now
			}
		}
	}()
}
