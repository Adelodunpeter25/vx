package tests

import (
	"testing"

	"github.com/Adelodunpeter25/vx/internal/search"
)

func TestEngine_SearchBasic(t *testing.T) {
	e := search.New()
	lines := []string{"hello world", "foo bar", "hello again"}

	matches := e.Search(lines, "hello")
	if len(matches) != 2 {
		t.Fatalf("expected 2 matches, got %d", len(matches))
	}
	if matches[0].Line != 0 || matches[0].Col != 0 {
		t.Fatalf("match 0: expected (0,0), got (%d,%d)", matches[0].Line, matches[0].Col)
	}
	if matches[1].Line != 2 || matches[1].Col != 0 {
		t.Fatalf("match 1: expected (2,0), got (%d,%d)", matches[1].Line, matches[1].Col)
	}
}

func TestEngine_SearchCaseInsensitive(t *testing.T) {
	e := search.New()
	lines := []string{"Hello WORLD", "hello world", "HELLO"}

	matches := e.Search(lines, "hello")
	if len(matches) != 3 {
		t.Fatalf("expected 3 case-insensitive matches, got %d", len(matches))
	}
}

func TestEngine_SearchNoMatches(t *testing.T) {
	e := search.New()
	lines := []string{"abc", "def", "ghi"}

	matches := e.Search(lines, "xyz")
	if len(matches) != 0 {
		t.Fatalf("expected 0 matches, got %d", len(matches))
	}
	if e.HasMatches() {
		t.Fatal("HasMatches should be false")
	}
}

func TestEngine_SearchEmptyQuery(t *testing.T) {
	e := search.New()
	lines := []string{"abc", "def"}

	matches := e.Search(lines, "")
	if len(matches) != 0 {
		t.Fatalf("empty query should return no matches, got %d", len(matches))
	}
}

func TestEngine_NextPrevious(t *testing.T) {
	e := search.New()
	lines := []string{"aaa", "bbb", "aaa", "ccc"}

	matches := e.Search(lines, "aaa")
	if len(matches) != 2 {
		t.Fatalf("expected 2 matches, got %d", len(matches))
	}

	// Current starts at 0
	cur := e.Current()
	if cur == nil || cur.Line != 0 {
		t.Fatalf("expected current at line 0, got: %+v", cur)
	}

	// Next should move to match 1
	next := e.Next()
	if next == nil || next.Line != 2 {
		t.Fatalf("expected next at line 2, got: %+v", next)
	}

	// Next wraps around
	next = e.Next()
	if next == nil || next.Line != 0 {
		t.Fatalf("expected wrap to line 0, got: %+v", next)
	}

	// Previous wraps around backward
	prev := e.Previous()
	if prev == nil || prev.Line != 2 {
		t.Fatalf("expected prev at line 2, got: %+v", prev)
	}
}

func TestEngine_Clear(t *testing.T) {
	e := search.New()
	lines := []string{"hello", "world"}

	e.Search(lines, "hello")
	if !e.HasMatches() {
		t.Fatal("should have matches")
	}

	e.Clear()
	if e.HasMatches() {
		t.Fatal("should have no matches after clear")
	}
	if e.Query() != "" {
		t.Fatal("query should be empty after clear")
	}
}

func TestEngine_MatchCountAndIndex(t *testing.T) {
	e := search.New()
	lines := []string{"a", "b", "a", "a"}

	e.Search(lines, "a")
	if e.MatchCount() != 3 {
		t.Fatalf("expected 3 matches, got %d", e.MatchCount())
	}
	if e.CurrentIndex() != 1 {
		t.Fatalf("expected current index 1, got %d", e.CurrentIndex())
	}

	e.Next()
	if e.CurrentIndex() != 2 {
		t.Fatalf("expected current index 2, got %d", e.CurrentIndex())
	}
}
