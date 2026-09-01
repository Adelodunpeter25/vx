package editor

import (
	"os"
	"path/filepath"

	"github.com/Adelodunpeter25/vx/internal/fswatch"
	"github.com/Adelodunpeter25/vx/internal/syntax"
)

func (e *Editor) initWatcher(root string) {
	w := fswatch.New(root)
	e.watcher = w
	// Start watcher async so recursive directory walk (filepath.Walk + inotify_add_watch)
	// does not block the first frame. File change events will arrive late, which is fine.
	go func() {
		if err := w.Start(); err != nil {
			return // watcher failure is non-fatal
		}
		for {
			select {
			case ev, ok := <-w.Events:
				if !ok {
					return
				}
				e.term.PostFileChange(ev.Path)
			case _, ok := <-w.Errors:
				if !ok {
					return
				}
			}
		}
	}()
}

// stopWatcher cleanly shuts down the file watcher and its goroutine.
func (e *Editor) stopWatcher() {
	if e.watcher != nil {
		e.watcher.Stop()
	}
}

// handleFileChange processes a file change event in the main goroutine.
func (e *Editor) handleFileChange(path string) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return
	}

	reloaded := false
	for _, p := range e.panes {
		if p.buffer.Filename() == "" {
			continue
		}
		absBufPath, err := filepath.Abs(p.buffer.Filename())
		if err != nil {
			continue
		}
		if absPath != absBufPath {
			continue
		}

		stat, err := os.Stat(p.buffer.Filename())
		if err == nil && !stat.ModTime().After(p.buffer.ModTime()) {
			continue
		}

		if p.buffer.IsModified() {
			p.msgManager.SetPersistent("File changed on disk (buffer has unsaved changes)")
			continue
		}

		// Reload buffer content
		if err := p.buffer.Reload(); err != nil {
			p.msgManager.SetError("Failed to reload file: " + err.Error())
			continue
		}

		// Refresh syntax engine
		p.syntax = syntax.New(p.buffer.Filename())

		// Clamp cursor and preserve scroll position
		e.clampCursorForPane(p)
		e.adjustScrollForPaneKeepScroll(p)

		// Invalidate render cache
		p.renderCache.invalidate()

		p.msgManager.SetTransient("File reloaded (external change)")
		reloaded = true
	}

	// Refresh file browser so new/deleted files appear
	if e.fileBrowser != nil {
		e.fileBrowser.Refresh()
	}

	// Only re-render if something actually changed or the browser is open
	if reloaded || e.fileBrowser != nil {
		e.render()
	}
}

// adjustScrollForPaneKeepScroll preserves scroll position after a reload,
// only adjusting if the cursor went out of the visible area.
func (e *Editor) adjustScrollForPaneKeepScroll(p *Pane) {
	if p == nil {
		return
	}
	savedVisualOffset := p.visualOffsetY
	savedOffsetY := p.offsetY
	e.adjustScrollForPane(p)
	// If viewport didn't need to move for cursor, restore original scroll
	if p.visualOffsetY == savedVisualOffset {
		p.visualOffsetY = savedVisualOffset
		p.offsetY = savedOffsetY
	}
}
