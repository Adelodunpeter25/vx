package tests

import (
	"os"
	"path/filepath"
	"testing"

	filebrowser "github.com/Adelodunpeter25/vx/internal/file-browser"
)

func TestBrowser_RefreshPreservesExpansion(t *testing.T) {
	dir := t.TempDir()
	// Create directory structure:
	// dir/
	//   sub1/
	//     file1.txt
	//   sub2/
	//     file2.txt
	os.MkdirAll(filepath.Join(dir, "sub1"), 0755)
	os.MkdirAll(filepath.Join(dir, "sub2"), 0755)
	os.WriteFile(filepath.Join(dir, "sub1", "file1.txt"), []byte("hi"), 0644)
	os.WriteFile(filepath.Join(dir, "sub2", "file2.txt"), []byte("hi"), 0644)

	s := filebrowser.New(dir)
	s.Open = true

	// Expand sub1 and sub2 by simulating visible nodes
	nodes := s.Visible()
	// Find and expand sub1 and sub2
	for _, n := range nodes {
		if n.Name == "sub1" || n.Name == "sub2" {
			n.Expanded = true
		}
	}

	// Get visible nodes after expansion (triggers loadChildren)
	_ = s.Visible()

	// Verify sub1 and sub2 are expanded
	var sub1Expanded, sub2Expanded bool
	for _, n := range nodes {
		if n.Name == "sub1" {
			sub1Expanded = n.Expanded
		}
		if n.Name == "sub2" {
			sub2Expanded = n.Expanded
		}
	}
	if !sub1Expanded || !sub2Expanded {
		t.Fatal("sub1 and sub2 should be expanded before refresh")
	}

	// Refresh
	s.Refresh()

	// Get visible nodes after refresh
	nodes = s.Visible()

	// Verify sub1 and sub2 are still expanded
	sub1Expanded = false
	sub2Expanded = false
	for _, n := range nodes {
		if n.Name == "sub1" {
			sub1Expanded = n.Expanded
		}
		if n.Name == "sub2" {
			sub2Expanded = n.Expanded
		}
	}
	if !sub1Expanded {
		t.Fatal("sub1 should still be expanded after refresh")
	}
	if !sub2Expanded {
		t.Fatal("sub2 should still be expanded after refresh")
	}
}

func TestBrowser_RefreshShowsNewFiles(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "existing.txt"), []byte("hi"), 0644)

	s := filebrowser.New(dir)
	s.Open = true

	nodes := s.Visible()
	found := false
	for _, n := range nodes {
		if n.Name == "existing.txt" {
			found = true
		}
	}
	if !found {
		t.Fatal("existing.txt should be visible")
	}

	// Create new file externally
	os.WriteFile(filepath.Join(dir, "new_file.txt"), []byte("new"), 0644)

	// Refresh and check
	s.Refresh()
	nodes = s.Visible()
	found = false
	for _, n := range nodes {
		if n.Name == "new_file.txt" {
			found = true
		}
	}
	if !found {
		t.Fatal("new_file.txt should be visible after refresh")
	}
}

func TestBrowser_RefreshRemovesDeletedFiles(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "to_delete.txt"), []byte("hi"), 0644)

	s := filebrowser.New(dir)
	s.Open = true

	nodes := s.Visible()
	found := false
	for _, n := range nodes {
		if n.Name == "to_delete.txt" {
			found = true
		}
	}
	if !found {
		t.Fatal("to_delete.txt should be visible")
	}

	// Delete file externally
	os.Remove(filepath.Join(dir, "to_delete.txt"))

	// Refresh and check
	s.Refresh()
	nodes = s.Visible()
	found = false
	for _, n := range nodes {
		if n.Name == "to_delete.txt" {
			found = true
		}
	}
	if found {
		t.Fatal("to_delete.txt should not be visible after refresh")
	}
}

func TestBrowser_InvalidationCascadesToCollapsedDirs(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "a", "b"), 0755)
	os.WriteFile(filepath.Join(dir, "a", "b", "deep.txt"), []byte("deep"), 0644)

	s := filebrowser.New(dir)
	s.Open = true

	// Expand "a" to see "b"
	nodes := s.Visible()
	for _, n := range nodes {
		if n.Name == "a" {
			n.Expanded = true
		}
	}
	nodes = s.Visible()

	// Expand "b" to see "deep.txt"
	for _, n := range nodes {
		if n.Name == "b" {
			n.Expanded = true
		}
	}
	nodes = s.Visible()

	foundDeep := false
	for _, n := range nodes {
		if n.Name == "deep.txt" {
			foundDeep = true
		}
	}
	if !foundDeep {
		t.Fatal("deep.txt should be visible")
	}

	// Now collapse "a" (set Expanded=false for a and b)
	for _, n := range nodes {
		if n.Name == "a" {
			n.Expanded = false
		}
	}

	// Create new file under a/b
	os.WriteFile(filepath.Join(dir, "a", "b", "new_deep.txt"), []byte("new"), 0644)

	// Refresh — must invalidate collapsed dirs too
	s.Refresh()

	// Re-expand "a" and "b"
	nodes = s.Visible()
	for _, n := range nodes {
		if n.Name == "a" {
			n.Expanded = true
		}
	}
	nodes = s.Visible()
	for _, n := range nodes {
		if n.Name == "b" {
			n.Expanded = true
		}
	}
	nodes = s.Visible()

	foundNew := false
	for _, n := range nodes {
		if n.Name == "new_deep.txt" {
			foundNew = true
		}
	}
	if !foundNew {
		t.Fatal("new_deep.txt should be visible after refresh and re-expand")
	}
}
