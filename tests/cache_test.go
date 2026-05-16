package tests

import (
	"testing"

	"github.com/Adelodunpeter25/vx/internal/editor"
)

func TestRenderCache_Invalidate(t *testing.T) {
	rc := editor.NewRenderCache()
	if !rc.NeedsRedraw() {
		t.Fatal("new cache should need redraw")
	}

	rc.SetNeedsRedraw(false)
	if rc.NeedsRedraw() {
		t.Fatal("should not need redraw after clearing")
	}

	rc.MarkNeedsRedraw()
	if !rc.NeedsRedraw() {
		t.Fatal("should need redraw after invalidate")
	}
}

func TestRenderCache_LineChanged(t *testing.T) {
	rc := editor.NewRenderCache()
	if !rc.LineChanged(0, "hello") {
		t.Fatal("unseen line should report as changed")
	}

	rc.UpdateLine(0, "hello")
	if rc.LineChanged(0, "hello") {
		t.Fatal("same content should not report as changed")
	}
	if !rc.LineChanged(0, "world") {
		t.Fatal("different content should report as changed")
	}
}

func TestRenderCache_InvalidateLine(t *testing.T) {
	rc := editor.NewRenderCache()
	rc.UpdateLine(0, "a")
	rc.UpdateLine(1, "b")

	rc.InvalidateLine(0)
	if !rc.LineChanged(0, "a") {
		t.Fatal("invalidated line should report as changed")
	}
	if rc.LineChanged(1, "b") {
		t.Fatal("non-invalidated line should still be cached")
	}
}
