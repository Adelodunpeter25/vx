package tests

import (
	"testing"

	splitpane "github.com/Adelodunpeter25/vx/internal/split-pane"
)

func TestLayoutSideBySide_SinglePane(t *testing.T) {
	rects, divX := splitpane.LayoutSideBySide(80, 24, 1, 0.5)
	if len(rects) != 1 {
		t.Fatalf("expected 1 rect, got %d", len(rects))
	}
	if rects[0].Width != 80 {
		t.Fatalf("expected width 80, got %d", rects[0].Width)
	}
	if rects[0].Height != 24 {
		t.Fatalf("expected height 24, got %d", rects[0].Height)
	}
	if divX != -1 {
		t.Fatalf("single pane should have dividerX=-1, got %d", divX)
	}
}

func TestLayoutSideBySide_TwoPanes(t *testing.T) {
	rects, divX := splitpane.LayoutSideBySide(100, 24, 2, 0.5)
	if len(rects) != 2 {
		t.Fatalf("expected 2 rects, got %d", len(rects))
	}
	// Total width should be 100 - 1 (divider) = 99
	if rects[0].Width+rects[1].Width != 99 {
		t.Fatalf("pane widths should sum to 99, got %d+%d", rects[0].Width, rects[1].Width)
	}
	if divX < 0 {
		t.Fatal("divider X should be non-negative for 2 panes")
	}
}

func TestLayoutSideBySide_RatioClamping(t *testing.T) {
	rects, _ := splitpane.LayoutSideBySide(100, 24, 2, 0.01) // too low, clamped to 0.1
	if rects[0].Width < 9 {
		t.Fatalf("ratio should be clamped to at least 0.1, left width is %d", rects[0].Width)
	}

	rects, _ = splitpane.LayoutSideBySide(100, 24, 2, 0.99) // too high, clamped to 0.9
	if rects[1].Width < 9 {
		t.Fatalf("ratio should be clamped to at most 0.9, right width is %d", rects[1].Width)
	}
}

func TestLayoutSideBySide_NarrowWidth(t *testing.T) {
	// Very narrow terminal — should still work without panicking
	rects, _ := splitpane.LayoutSideBySide(1, 24, 2, 0.5)
	if len(rects) == 0 {
		t.Fatal("should return at least one rect")
	}
}
