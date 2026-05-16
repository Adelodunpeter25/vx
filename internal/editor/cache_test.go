package editor

import (
	"testing"
)

func TestRenderCache_Invalidate(t *testing.T) {
	rc := newRenderCache()
	if !rc.needsRedraw {
		t.Fatal("new cache should need redraw")
	}

	rc.needsRedraw = false
	if rc.needsRedraw {
		t.Fatal("should not need redraw after clearing")
	}

	rc.invalidate()
	if !rc.needsRedraw {
		t.Fatal("should need redraw after invalidate")
	}
}

func TestRenderCache_LineChanged(t *testing.T) {
	rc := newRenderCache()
	// Line not in cache — should report changed
	if !rc.lineChanged(0, "hello") {
		t.Fatal("unseen line should report as changed")
	}

	rc.updateLine(0, "hello")
	// Same content — should NOT report changed
	if rc.lineChanged(0, "hello") {
		t.Fatal("same content should not report as changed")
	}
	// Different content — should report changed
	if !rc.lineChanged(0, "world") {
		t.Fatal("different content should report as changed")
	}
}

func TestRenderCache_InvalidateLine(t *testing.T) {
	rc := newRenderCache()
	rc.updateLine(0, "a")
	rc.updateLine(1, "b")

	rc.invalidateLine(0)
	if !rc.lineChanged(0, "a") {
		t.Fatal("invalidated line should report as changed")
	}
	if rc.lineChanged(1, "b") {
		t.Fatal("non-invalidated line should still be cached")
	}
}
