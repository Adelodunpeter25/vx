package fswatch

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

// EventType defines the type of file system change
type EventType int

const (
	EventCreated EventType = iota
	EventModified
	EventDeleted
)

// Event represents a single file system change
type Event struct {
	Path string
	Type EventType
}

// Watcher monitors a directory for changes using fsnotify
type Watcher struct {
	root    string
	Events  chan Event
	Errors  chan error
	watcher *fsnotify.Watcher
	cancel  context.CancelFunc
}

// New creates a new file system watcher
func New(root string) *Watcher {
	return &Watcher{
		root:   root,
		Events: make(chan Event, 100),
		Errors: make(chan error, 10),
	}
}

// Start begins monitoring the root directory
func (w *Watcher) Start() error {
	fsw, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	w.watcher = fsw

	// Recursively add all directories under root
	if err := w.addDirsRecursive(w.root); err != nil {
		_ = fsw.Close()
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	w.cancel = cancel

	go w.loop(ctx)
	return nil
}

// Stop halts the watcher and cleans up goroutines
func (w *Watcher) Stop() {
	if w.cancel != nil {
		w.cancel()
	}
	if w.watcher != nil {
		_ = w.watcher.Close()
	}
}

func (w *Watcher) addDirsRecursive(root string) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // skip inaccessible paths
		}
		if !info.IsDir() {
			return nil
		}
		name := info.Name()
		if name != "." && name != root {
			if shouldSkip(name) {
				return filepath.SkipDir
			}
		}
		return w.watcher.Add(path)
	})
}

func shouldSkip(name string) bool {
	return strings.HasPrefix(name, ".") || name == "node_modules" || name == "vendor"
}

func (w *Watcher) loop(ctx context.Context) {
	// Debounce timer: coalesce rapid events within 50ms
	debounce := time.NewTimer(0)
	if !debounce.Stop() {
		<-debounce.C
	}
	defer debounce.Stop()

	pending := make(map[string]EventType)

	for {
		select {
		case <-ctx.Done():
			return

		case ev, ok := <-w.watcher.Events:
			if !ok {
				return
			}
			name := filepath.Base(ev.Name)
			if shouldSkip(name) {
				continue
			}

			// Translate fsnotify events to our event types
			switch {
			case ev.Op&(fsnotify.Create) != 0:
				// If a new directory was created, start watching it
				if info, err := os.Stat(ev.Name); err == nil && info.IsDir() {
					_ = w.addDirsRecursive(ev.Name)
				}
				pending[ev.Name] = EventCreated

			case ev.Op&(fsnotify.Write|fsnotify.Chmod) != 0:
				// Only emit modified for files, not directories
				if info, err := os.Stat(ev.Name); err == nil && !info.IsDir() {
					// Don't overwrite a Create with a Modify
					if _, exists := pending[ev.Name]; !exists {
						pending[ev.Name] = EventModified
					}
				}

			case ev.Op&(fsnotify.Remove|fsnotify.Rename) != 0:
				// Remove from watcher (no-op if already gone)
				_ = w.watcher.Remove(ev.Name)
				// Don't overwrite a Create with a Delete
				if _, exists := pending[ev.Name]; !exists {
					pending[ev.Name] = EventDeleted
				}
			}

			// Reset debounce timer
			debounce.Reset(50 * time.Millisecond)

		case err, ok := <-w.watcher.Errors:
			if !ok {
				return
			}
			select {
			case w.Errors <- err:
			default:
			}

		case <-debounce.C:
			// Flush pending events
			for path, typ := range pending {
				select {
				case w.Events <- Event{Path: path, Type: typ}:
				default:
					// Channel full — drop with warning
					select {
					case w.Errors <- fmt.Errorf("fswatch: event channel full, dropping event for %s", path):
					default:
					}
				}
			}
			pending = make(map[string]EventType)
		}
	}
}
