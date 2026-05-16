package tests

import (
	"os"
	"testing"

	"github.com/Adelodunpeter25/vx/internal/utils"
)

func TestValidateUTF8_Valid(t *testing.T) {
	s := utils.ValidateUTF8("hello world")
	if s != "hello world" {
		t.Fatalf("valid UTF-8 should pass through, got '%s'", s)
	}
}

func TestValidateUTF8_Invalid(t *testing.T) {
	s := utils.ValidateUTF8("hello\xff\xfe world")
	// Should contain replacement characters
	if s == "hello\xff\xfe world" {
		t.Fatal("invalid UTF-8 should be cleaned")
	}
	if len(s) == 0 {
		t.Fatal("result should not be empty")
	}
}

func TestFormatFileSize(t *testing.T) {
	tests := []struct {
		bytes int64
		want  string
	}{
		{0, "0 B"},
		{500, "500 B"},
		{1024, "1.0 KB"},
		{1536, "1.5 KB"},
		{1048576, "1.0 MB"},
		{1073741824, "1.0 GB"},
	}
	for _, tt := range tests {
		got := utils.FormatFileSize(tt.bytes)
		if got != tt.want {
			t.Errorf("FormatFileSize(%d) = '%s', want '%s'", tt.bytes, got, tt.want)
		}
	}
}

func TestFormatLineCount(t *testing.T) {
	if utils.FormatLineCount(1) != "1 line" {
		t.Fatalf("expected '1 line', got '%s'", utils.FormatLineCount(1))
	}
	if utils.FormatLineCount(5) != "5 lines" {
		t.Fatalf("expected '5 lines', got '%s'", utils.FormatLineCount(5))
	}
}

func TestFormatFileInfo(t *testing.T) {
	info := utils.FormatFileInfo("test.go", 1024, 42)
	if info != `"test.go" 1.0 KB, 42 lines` {
		t.Fatalf("unexpected format: %s", info)
	}
}

func TestFileError_Error(t *testing.T) {
	err := utils.NewFileError("load", "test.txt", &testErr{"not found"})
	if err.Error() != "load test.txt: not found" {
		t.Fatalf("unexpected error format: %s", err.Error())
	}
}

type testErr struct{ msg string }

func (e *testErr) Error() string { return e.msg }

func TestFormatUserError(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"no such file or directory", "File not found"},
		{"permission denied", "Permission denied - check file permissions"},
		{"read-only file system", "Cannot save - file system is read-only"},
		{"no space left on device", "Cannot save - disk is full"},
		{"is a directory", "Cannot open - this is a directory, not a file"},
	}
	for _, tt := range tests {
		got := utils.FormatUserError(&testErr{tt.input})
		if got != tt.want {
			t.Errorf("FormatUserError(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestCountLines(t *testing.T) {
	dir := t.TempDir()
	path := dir + "/lines.txt"
	os.WriteFile(path, []byte("a\nb\nc\nd\ne\n"), 0644)

	count, err := utils.CountLines(path)
	if err != nil {
		t.Fatalf("CountLines error: %v", err)
	}
	if count != 5 {
		t.Fatalf("expected 5 lines, got %d", count)
	}
}
