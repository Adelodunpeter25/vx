package tests

import (
	"testing"

	"github.com/Adelodunpeter25/vx/internal/wrap"
)

func TestWrapLine_Short(t *testing.T) {
	lines := wrap.WrapLine("hello", 0, 80)
	if len(lines) != 1 {
		t.Fatalf("short line should have 1 segment, got %d", len(lines))
	}
	if lines[0].Text != "hello" {
		t.Fatalf("expected 'hello', got '%s'", lines[0].Text)
	}
	if lines[0].IsWrapped {
		t.Fatal("first segment should not be wrapped")
	}
}

func TestWrapLine_Long(t *testing.T) {
	lines := wrap.WrapLine("abcdefghijklmnop", 0, 5)
	if len(lines) != 4 {
		t.Fatalf("expected 4 segments, got %d", len(lines))
	}
	if lines[0].Text != "abcde" {
		t.Fatalf("seg 0: expected 'abcde', got '%s'", lines[0].Text)
	}
	if !lines[1].IsWrapped {
		t.Fatal("seg 1 should be wrapped")
	}
	if lines[1].StartCol != 5 {
		t.Fatalf("seg 1 StartCol: expected 5, got %d", lines[1].StartCol)
	}
}

func TestWrapLine_Empty(t *testing.T) {
	lines := wrap.WrapLine("", 0, 10)
	if len(lines) != 1 {
		t.Fatalf("empty line should have 1 segment, got %d", len(lines))
	}
}

func TestWrapLine_ZeroWidth(t *testing.T) {
	lines := wrap.WrapLine("hello", 0, 0)
	if len(lines) != 1 {
		t.Fatalf("zero width should return 1 segment, got %d", len(lines))
	}
}

func TestVisualLineCount(t *testing.T) {
	if wrap.VisualLineCount("hello", 80) != 1 {
		t.Fatal("short line should be 1 visual line")
	}
	if wrap.VisualLineCount("abcdefghijklmnop", 5) != 4 {
		t.Fatal("16 chars with width 5 should be 4 visual lines")
	}
	if wrap.VisualLineCount("", 80) != 1 {
		t.Fatal("empty line should be 1 visual line")
	}
	if wrap.VisualLineCount("hello", 0) != 1 {
		t.Fatal("zero width should return 1")
	}
}

func TestTotalVisualLines(t *testing.T) {
	lines := []string{"hello", "abcdefghijklmnop", "hi"}
	total := wrap.TotalVisualLines(lines, 0, 2, 5)
	// "hello" = 1, "abcdefghijklmnop" = 4 (16/5=3.2->4), "hi" = 1
	expected := 1 + 4 + 1
	if total != expected {
		t.Fatalf("expected %d visual lines, got %d", expected, total)
	}
}
