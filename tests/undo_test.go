package tests

import (
	"testing"

	"github.com/Adelodunpeter25/vx/internal/undo"
)

func TestStack_PushAndUndo(t *testing.T) {
	s := undo.NewStack()
	if s.CanUndo() {
		t.Fatal("new stack should not be undoable")
	}

	s.Push(undo.Action{Type: undo.ActionInsertRune, Line: 0, Col: 0, Text: "a"})
	if !s.CanUndo() {
		t.Fatal("stack should be undoable after push")
	}

	action := s.Undo()
	if action == nil {
		t.Fatal("undo returned nil")
	}
	if action.Line != 0 || action.Text != "a" {
		t.Fatalf("unexpected action: %+v", action)
	}
	if s.CanUndo() {
		t.Fatal("should not be undoable after single undo")
	}
}

func TestStack_UndoRedo(t *testing.T) {
	s := undo.NewStack()
	s.Push(undo.Action{Type: undo.ActionInsertRune, Line: 1, Col: 5, Text: "x"})
	s.Push(undo.Action{Type: undo.ActionDeleteRune, Line: 2, Col: 3, OldText: "y"})

	// Undo second action
	action := s.Undo()
	if action == nil || action.Type != undo.ActionDeleteRune {
		t.Fatalf("expected delete action, got: %+v", action)
	}
	if !s.CanRedo() {
		t.Fatal("should be redoable after undo")
	}

	// Undo first action
	action = s.Undo()
	if action == nil || action.Type != undo.ActionInsertRune {
		t.Fatalf("expected insert action, got: %+v", action)
	}
	if s.CanUndo() {
		t.Fatal("should not be undoable after undoing all")
	}

	// Redo first action
	action = s.Redo()
	if action == nil || action.Type != undo.ActionInsertRune {
		t.Fatalf("expected insert action on redo, got: %+v", action)
	}
}

func TestStack_PushTruncatesRedoHistory(t *testing.T) {
	s := undo.NewStack()
	s.Push(undo.Action{Type: undo.ActionInsertRune, Line: 0, Col: 0, Text: "a"})
	s.Push(undo.Action{Type: undo.ActionInsertRune, Line: 1, Col: 0, Text: "b"})

	s.Undo() // undo "b"
	if !s.CanRedo() {
		t.Fatal("should be redoable")
	}

	// Push new action — should truncate redo of "b"
	s.Push(undo.Action{Type: undo.ActionInsertLine, Line: 1, Col: 0, Text: "c"})
	if s.CanRedo() {
		t.Fatal("push should have truncated redo history")
	}

	// Undo should go to "c", not "b"
	action := s.Undo()
	if action == nil || action.Text != "c" {
		t.Fatalf("expected 'c', got: %+v", action)
	}
}

func TestStack_Clear(t *testing.T) {
	s := undo.NewStack()
	s.Push(undo.Action{Type: undo.ActionInsertRune, Line: 0, Col: 0, Text: "a"})
	s.Push(undo.Action{Type: undo.ActionDeleteRune, Line: 1, Col: 0})
	s.Clear()

	if s.CanUndo() || s.CanRedo() {
		t.Fatal("cleared stack should have no undo/redo")
	}
	if s.Undo() != nil {
		t.Fatal("undo on cleared stack should return nil")
	}
	if s.Redo() != nil {
		t.Fatal("redo on cleared stack should return nil")
	}
}

func TestStack_UndoOnEmpty(t *testing.T) {
	s := undo.NewStack()
	if s.Undo() != nil {
		t.Fatal("undo on empty stack should return nil")
	}
	if s.Redo() != nil {
		t.Fatal("redo on empty stack should return nil")
	}
}
