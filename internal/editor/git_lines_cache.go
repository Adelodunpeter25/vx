package editor

import (
	"context"
	"path/filepath"
	"time"

	"github.com/Adelodunpeter25/vx/internal/git"
)

func (e *Editor) requestGitLines(p *Pane) {
	if e == nil || p == nil || p.buffer == nil {
		return
	}
	filename := p.buffer.Filename()
	if filename == "" {
		return
	}
	key := filepath.Clean(filename)
	version := p.buffer.ModVersion()

	e.gitLinesMu.Lock()
	if e.gitLines == nil {
		e.gitLines = make(map[string]*gitLineCacheEntry)
	}
	entry, ok := e.gitLines[key]
	if ok {
		if entry.ready && entry.file == key && entry.version == version {
			e.gitLinesMu.Unlock()
			return
		}
		if entry.loading && entry.version == version {
			e.gitLinesMu.Unlock()
			return
		}
	}
	if !ok {
		entry = &gitLineCacheEntry{}
		e.gitLines[key] = entry
	}
	entry.file = key
	entry.version = version
	entry.loading = true
	entry.ready = false
	e.gitLinesMu.Unlock()

	go func(file string, ver int) {
		ctx, cancel := context.WithTimeout(context.Background(), 1500*time.Millisecond)
		defer cancel()

		lines, err := git.DiffLineChanges(ctx, file)

		e.gitLinesMu.Lock()
		defer e.gitLinesMu.Unlock()

		current, ok := e.gitLines[file]
		if !ok || current.version != ver {
			return
		}
		current.loading = false
		if err != nil {
			current.lines = nil
			current.ready = true
			return
		}
		current.lines = lines
		current.ready = true
	}(key, version)
}

func (e *Editor) gitLinesFor(p *Pane) map[int]git.LineChange {
	if e == nil || p == nil || p.buffer == nil {
		return nil
	}
	filename := p.buffer.Filename()
	if filename == "" {
		return nil
	}
	key := filepath.Clean(filename)
	version := p.buffer.ModVersion()

	e.gitLinesMu.RLock()
	entry, ok := e.gitLines[key]
	e.gitLinesMu.RUnlock()
	if !ok || entry == nil || entry.version != version || !entry.ready {
		e.requestGitLines(p)
		return nil
	}
	return entry.lines
}
