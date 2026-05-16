package tests

import (
	"testing"
	"time"

	"github.com/Adelodunpeter25/vx/internal/editor"
)

func TestMessageManager_SetAndGet(t *testing.T) {
	mm := editor.NewMessageManager()
	mm.SetPersistent("hello")
	if mm.Get() != "hello" {
		t.Fatalf("expected 'hello', got '%s'", mm.Get())
	}
}

func TestMessageManager_TransientExpiry(t *testing.T) {
	mm := editor.NewMessageManager()
	mm.SetTransient("flash")
	if mm.Get() != "flash" {
		t.Fatal("transient message should be visible immediately")
	}
	mm.Clear()
	if mm.Get() != "" {
		t.Fatal("message should be empty after clear")
	}
}

func TestMessageManager_Clear(t *testing.T) {
	mm := editor.NewMessageManager()
	mm.SetPersistent("persistent")
	mm.Clear()
	if mm.Get() != "" {
		t.Fatal("message should be empty after clear")
	}
}

func TestMessageManager_ClearIfTransient(t *testing.T) {
	mm := editor.NewMessageManager()
	mm.SetPersistent("persist")
	mm.ClearIfTransient()
	if mm.Get() != "persist" {
		t.Fatal("persistent message should survive ClearIfTransient")
	}

	mm.SetTransient("transient")
	time.Sleep(600 * time.Millisecond)
	mm.ClearIfTransient()
	if mm.Get() != "" {
		t.Fatal("transient message should be cleared after 500ms guard")
	}
}

func TestMessageManager_EmptyGet(t *testing.T) {
	mm := editor.NewMessageManager()
	if mm.Get() != "" {
		t.Fatal("empty manager should return empty string")
	}
}
