package editor

import (
	"github.com/Adelodunpeter25/vx/internal/fswatch"
)

func (e *Editor) initWatcher(root string) {
	e.watcher = fswatch.New(root)
	e.watcher.Start()

	go func() {
		for {
			select {
			case <-e.watcher.Events:
				if e.fileBrowser != nil {
					// We need to refresh the root children if anything changed
					// For simplicity, we'll just trigger a re-render and let fileBrowser
					// handle expansion if it needs to. 
					// We'll add a Refresh method to the fileBrowser state later if needed.
					// For now, just invalidate the root's loaded state.
					e.fileBrowser.Refresh()
				}
			case <-e.watcher.Errors:
				// Silent errors for now
			}
		}
	}()
}
