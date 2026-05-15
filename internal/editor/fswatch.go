package editor

import (
	"path/filepath"

	"github.com/Adelodunpeter25/vx/internal/buffer"
	"github.com/Adelodunpeter25/vx/internal/fswatch"
)

func (e *Editor) initWatcher(root string) {
	e.watcher = fswatch.New(root)
	e.watcher.Start()

	go func() {
		for {
			select {
			case ev := <-e.watcher.Events:
				// Check if the modified file is open in any pane
				absEvPath, _ := filepath.Abs(ev.Path)

				for _, p := range e.panes {
					if p.buffer.Filename() == "" {
						continue
					}
					absBufPath, _ := filepath.Abs(p.buffer.Filename())

					if absEvPath == absBufPath && ev.Type == fswatch.EventModified {
						if !p.buffer.IsModified() {
							// Reload buffer content while preserving cursor
							newBuf, err := buffer.Load(p.buffer.Filename())
							if err == nil {
								p.buffer = newBuf
								e.clampCursorFor(p)
								p.msgManager.SetTransient("File reloaded (external change)")
							}
						}
					}
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

func (e *Editor) clampCursorFor(p *Pane) {
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
