package fswatch

import (
	"os"
	"path/filepath"
	"sync"
	"time"
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

// Watcher monitors a directory for changes
type Watcher struct {
	root     string
	Events   chan Event
	Errors   chan error
	mu       sync.Mutex
	snapshot map[string]time.Time
	stop     chan struct{}
}

// New creates a new file system watcher
func New(root string) *Watcher {
	return &Watcher{
		root:     root,
		Events:   make(chan Event, 100),
		Errors:   make(chan error, 10),
		snapshot: make(map[string]time.Time),
		stop:     make(chan struct{}),
	}
}

// Start begins monitoring the root directory
func (w *Watcher) Start() {
	// Initial snapshot
	w.scan()

	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				w.scan()
			case <-w.stop:
				return
			}
		}
	}()
}

// Stop halts the watcher
func (w *Watcher) Stop() {
	close(w.stop)
}

func (w *Watcher) scan() {
	w.mu.Lock()
	defer w.mu.Unlock()

	current := make(map[string]time.Time)

	err := filepath.Walk(w.root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		// Skip hidden directories (like .git) to save performance
		if info.IsDir() && filepath.Base(path) != "." && filepath.Base(path)[0] == '.' {
			return filepath.SkipDir
		}

		current[path] = info.ModTime()

		lastMod, exists := w.snapshot[path]
		if !exists {
			w.Events <- Event{Path: path, Type: EventCreated}
		} else if info.ModTime().After(lastMod) {
			w.Events <- Event{Path: path, Type: EventModified}
		}

		return nil
	})

	if err != nil {
		w.Errors <- err
		return
	}

	// Check for deletions
	for path := range w.snapshot {
		if _, exists := current[path]; !exists {
			w.Events <- Event{Path: path, Type: EventDeleted}
		}
	}

	w.snapshot = current
}
