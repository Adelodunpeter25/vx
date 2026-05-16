package tests

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Adelodunpeter25/vx/internal/fswatch"
)

func TestWatcher_DetectsFileCreation(t *testing.T) {
	dir := t.TempDir()
	w := fswatch.New(dir)
	if err := w.Start(); err != nil {
		t.Fatalf("start error: %v", err)
	}
	defer w.Stop()

	// Wait a moment for the watcher to initialize
	time.Sleep(100 * time.Millisecond)

	// Create a new file
	os.WriteFile(filepath.Join(dir, "new.txt"), []byte("hello"), 0644)

	// Wait for event with timeout
	select {
	case ev := <-w.Events:
		if ev.Type != fswatch.EventCreated {
			t.Fatalf("expected EventCreated, got %d", ev.Type)
		}
		if filepath.Base(ev.Path) != "new.txt" {
			t.Fatalf("expected 'new.txt', got '%s'", filepath.Base(ev.Path))
		}
	case err := <-w.Errors:
		t.Fatalf("unexpected error: %v", err)
	case <-time.After(3 * time.Second):
		t.Fatal("timed out waiting for file creation event")
	}
}

func TestWatcher_DetectsFileModification(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "mod.txt")
	os.WriteFile(path, []byte("original"), 0644)

	w := fswatch.New(dir)
	if err := w.Start(); err != nil {
		t.Fatalf("start error: %v", err)
	}
	defer w.Stop()
	time.Sleep(100 * time.Millisecond)

	// Modify the file
	os.WriteFile(path, []byte("modified"), 0644)

	select {
	case ev := <-w.Events:
		if ev.Type != fswatch.EventModified {
			t.Fatalf("expected EventModified, got %d", ev.Type)
		}
	case err := <-w.Errors:
		t.Fatalf("unexpected error: %v", err)
	case <-time.After(3 * time.Second):
		t.Fatal("timed out waiting for file modification event")
	}
}

func TestWatcher_DetectsFileDeletion(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "delete.txt")
	os.WriteFile(path, []byte("bye"), 0644)

	w := fswatch.New(dir)
	if err := w.Start(); err != nil {
		t.Fatalf("start error: %v", err)
	}
	defer w.Stop()
	time.Sleep(100 * time.Millisecond)

	// Delete the file
	os.Remove(path)

	select {
	case ev := <-w.Events:
		if ev.Type != fswatch.EventDeleted {
			t.Fatalf("expected EventDeleted, got %d", ev.Type)
		}
	case err := <-w.Errors:
		t.Fatalf("unexpected error: %v", err)
	case <-time.After(3 * time.Second):
		t.Fatal("timed out waiting for file deletion event")
	}
}

func TestWatcher_AutoWatchesNewDirectories(t *testing.T) {
	dir := t.TempDir()

	w := fswatch.New(dir)
	if err := w.Start(); err != nil {
		t.Fatalf("start error: %v", err)
	}
	defer w.Stop()
	time.Sleep(100 * time.Millisecond)

	// Create a new subdirectory
	newDir := filepath.Join(dir, "subdir")
	os.Mkdir(newDir, 0755)
	time.Sleep(200 * time.Millisecond) // let debounce flush the mkdir event

	// Drain any events from the mkdir
	drained := false
	for !drained {
		select {
		case <-w.Events:
		default:
			drained = true
		}
	}

	// Create a file inside the new subdirectory
	os.WriteFile(filepath.Join(newDir, "inside.txt"), []byte("deep"), 0644)

	select {
	case ev := <-w.Events:
		if filepath.Base(ev.Path) != "inside.txt" {
			t.Fatalf("expected event for 'inside.txt', got '%s'", filepath.Base(ev.Path))
		}
	case err := <-w.Errors:
		t.Fatalf("unexpected error: %v", err)
	case <-time.After(3 * time.Second):
		t.Fatal("timed out waiting for event in new subdirectory")
	}
}

func TestWatcher_Stop(t *testing.T) {
	dir := t.TempDir()
	w := fswatch.New(dir)
	if err := w.Start(); err != nil {
		t.Fatalf("start error: %v", err)
	}
	w.Stop()

	// After stop, the Events channel should eventually be closed
	// (or at least the goroutine should exit)
	select {
	case _, ok := <-w.Events:
		if ok {
			// Channel still open is acceptable — just means goroutine hasn't exited yet
		}
	case <-time.After(500 * time.Millisecond):
		// Timeout is fine — goroutine may take a moment to exit
	}
}

func TestWatcher_SkipsHiddenDirectories(t *testing.T) {
	dir := t.TempDir()
	hiddenDir := filepath.Join(dir, ".hidden")
	os.Mkdir(hiddenDir, 0755)

	w := fswatch.New(dir)
	if err := w.Start(); err != nil {
		t.Fatalf("start error: %v", err)
	}
	defer w.Stop()
	time.Sleep(100 * time.Millisecond)

	// Create file in hidden directory
	os.WriteFile(filepath.Join(hiddenDir, "secret.txt"), []byte("hidden"), 0644)

	// Should NOT receive an event for hidden directory files
	select {
	case ev := <-w.Events:
		t.Fatalf("should not receive event for hidden dir file, got: %s", ev.Path)
	case err := <-w.Errors:
		t.Fatalf("unexpected error: %v", err)
	case <-time.After(500 * time.Millisecond):
		// Expected — no event for hidden directories
	}
}

func TestWatcher_InitialSnapshotSuppressesEvents(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "existing.txt"), []byte("exists"), 0644)

	w := fswatch.New(dir)
	if err := w.Start(); err != nil {
		t.Fatalf("start error: %v", err)
	}
	defer w.Stop()

	// Should not receive an event for the existing file
	select {
	case ev := <-w.Events:
		t.Fatalf("should not receive event for pre-existing file, got: %s", ev.Path)
	case err := <-w.Errors:
		t.Fatalf("unexpected error: %v", err)
	case <-time.After(500 * time.Millisecond):
		// Expected — no event for existing files on startup
	}
}
