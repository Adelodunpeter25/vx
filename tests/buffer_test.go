package tests

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Adelodunpeter25/vx/internal/buffer"
)

func TestBuffer_New(t *testing.T) {
	b := buffer.New()
	if b.LineCount() != 1 {
		t.Fatalf("new buffer should have 1 line, got %d", b.LineCount())
	}
	if b.Line(0) != "" {
		t.Fatal("new buffer's first line should be empty")
	}
	if b.IsModified() {
		t.Fatal("new buffer should not be modified")
	}
	if b.Filename() != "" {
		t.Fatal("new buffer should have no filename")
	}
}

func TestBuffer_Load(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.txt")
	os.WriteFile(path, []byte("line1\nline2\nline3\n"), 0644)

	b, err := buffer.Load(path)
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	if b.LineCount() != 3 {
		t.Fatalf("expected 3 lines, got %d", b.LineCount())
	}
	if b.Line(0) != "line1" {
		t.Fatalf("line 0: expected 'line1', got '%s'", b.Line(0))
	}
	if b.Line(2) != "line3" {
		t.Fatalf("line 2: expected 'line3', got '%s'", b.Line(2))
	}
	if b.Filename() != path {
		t.Fatalf("filename mismatch: %s", b.Filename())
	}
	if b.IsModified() {
		t.Fatal("loaded buffer should not be modified")
	}
}

func TestBuffer_LoadEmpty(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "empty.txt")
	os.WriteFile(path, []byte(""), 0644)

	b, err := buffer.Load(path)
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	if b.LineCount() != 1 {
		t.Fatalf("empty file should have 1 line, got %d", b.LineCount())
	}
}

func TestBuffer_LoadNonExistent(t *testing.T) {
	b, err := buffer.Load("/nonexistent/path/file.txt")
	if err != nil {
		t.Fatalf("non-existent file should not error, got: %v", err)
	}
	if b == nil {
		t.Fatal("should return empty buffer for non-existent file")
	}
	if b.LineCount() != 1 {
		t.Fatalf("non-existent file buffer should have 1 line, got %d", b.LineCount())
	}
}

func TestBuffer_Save(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "save_test.txt")

	b := buffer.New()
	b.SetFilename(path)
	b.InsertRune(0, 0, 'h')
	b.InsertRune(0, 1, 'i')

	if err := b.Save(); err != nil {
		t.Fatalf("save error: %v", err)
	}
	if b.IsModified() {
		t.Fatal("saved buffer should not be modified")
	}

	// Read back
	data, _ := os.ReadFile(path)
	if string(data) != "hi" {
		t.Fatalf("expected 'hi', got '%s'", string(data))
	}
}

func TestBuffer_Reload(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "reload_test.txt")
	os.WriteFile(path, []byte("original\n"), 0644)

	b, _ := buffer.Load(path)
	if b.Line(0) != "original" {
		t.Fatalf("expected 'original', got '%s'", b.Line(0))
	}

	// Modify file externally
	os.WriteFile(path, []byte("modified\n"), 0644)

	if err := b.Reload(); err != nil {
		t.Fatalf("reload error: %v", err)
	}
	if b.Line(0) != "modified" {
		t.Fatalf("after reload expected 'modified', got '%s'", b.Line(0))
	}
	if b.IsModified() {
		t.Fatal("reloaded buffer should not be modified")
	}
}

func TestBuffer_LineBounds(t *testing.T) {
	b := buffer.New()
	if b.Line(-1) != "" {
		t.Fatal("Line(-1) should return empty")
	}
	if b.Line(1) != "" {
		t.Fatal("Line(1) on single-line buffer should return empty")
	}
}

func TestBuffer_ModVersion(t *testing.T) {
	b := buffer.New()
	v0 := b.ModVersion()
	b.InsertRune(0, 0, 'x')
	if b.ModVersion() != v0+1 {
		t.Fatal("modVersion should increment on edit")
	}
}
