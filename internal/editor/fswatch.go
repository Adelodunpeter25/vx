package editor

import (
	"path/filepath"

	"github.com/Adelodunpeter25/vx/internal/fswatch"
	"github.com/Adelodunpeter25/vx/internal/syntax"
)

func (e *Editor) initWatcher(root string) {
	e.watcher = fswatch.New(root)
	if err := e.watcher.Start(); err != nil {
		return // watcher failure is non-fatal
	}

	// This goroutine only forwards events into the tcell event loop.
	// It never touches editor state directly — no data races.
	go func() {
		for {
			select {
			case ev, ok := <-e.watcher.Events:
				if !ok {
					return
				}
				e.term.PostFileChange(ev.Path)
			case _, ok := <-e.watcher.Errors:
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

// clampCursorForPane clamps cursor position to valid bounds within the buffer.
func (e *Editor) clampCursorForPane(p *Pane) {
	if p.cursorY >= p.buffer.LineCount() {
		p.cursorY = p.buffer.LineCount() - 1
	}
	if p.cursorY < 0 {
		p.cursorY = 0
	}
	line := p.buffer.Line(p.cursorY)
	lineLen := lineRuneCount(line)
	if p.cursorX > lineLen {
		p.cursorX = lineLen
	}
}

// adjustScrollForPaneKeepScroll preserves scroll position after a reload,
// only adjusting if the cursor went out of the visible area.
func (e *Editor) adjustScrollForPaneKeepScroll(p *Pane) {
	saved := e.activePane
	for i, pane := range e.panes {
		if pane == p {
			e.activePane = i
			break
		}
	}
	savedVisualOffset := p.visualOffsetY
	savedOffsetY := p.offsetY
	e.adjustScroll()
	// If the cursor is still within the original viewport, restore the scroll
	if p.visualOffsetY == savedVisualOffset {
		p.visualOffsetY = savedVisualOffset
		p.offsetY = savedOffsetY
	}
	e.activePane = saved
}
