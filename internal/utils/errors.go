package utils

import (
	"fmt"
	"strings"
)

// FileError wraps file operation errors with context
type FileError struct {
	Op       string
	Filename string
	Err      error
}

func (e *FileError) Error() string {
	return fmt.Sprintf("%s %s: %v", e.Op, e.Filename, e.Err)
}

// NewFileError creates a new file error
func NewFileError(op, filename string, err error) error {
	return &FileError{
		Op:       op,
		Filename: filename,
		Err:      err,
	}
}

// FormatUserError converts technical errors into user-friendly messages
func FormatUserError(err error) string {
	if err == nil {
		return ""
	}
	
	errMsg := err.Error()
	
	// File not found
	if strings.Contains(errMsg, "no such file or directory") {
		return "File not found"
	}
	
	// Permission denied
	if strings.Contains(errMsg, "permission denied") {
		return "Permission denied - check file permissions"
	}
	
	// File too large
	if strings.Contains(errMsg, "file too large") {
		return errMsg // Already user-friendly
	}
	
	// Too many lines
	if strings.Contains(errMsg, "too many lines") {
		return errMsg // Already user-friendly
	}
	
	// No write since last change
	if strings.Contains(errMsg, "no write since last change") {
		return errMsg // Already user-friendly
	}
	
	// No file name
	if strings.Contains(errMsg, "no file name") {
		return "No filename specified"
	}
	
	// Read-only file system
	if strings.Contains(errMsg, "read-only file system") {
		return "Cannot save - file system is read-only"
	}
	
	// Disk full
	if strings.Contains(errMsg, "no space left on device") {
		return "Cannot save - disk is full"
	}
	
	// Is a directory
	if strings.Contains(errMsg, "is a directory") {
		return "Cannot open - this is a directory, not a file"
	}
	
	// Default: return original error but make it friendlier
	return fmt.Sprintf("Error: %s", errMsg)
}

// FormatSaveError formats save-specific errors
func FormatSaveError(filename string, err error) string {
	if err == nil {
		return ""
	}
	
	baseMsg := FormatUserError(err)
	
	// Add context for save operations
	if strings.Contains(baseMsg, "Permission denied") {
		return fmt.Sprintf("Cannot save '%s' - permission denied", filename)
	}
	
	if strings.Contains(baseMsg, "read-only") {
		return fmt.Sprintf("Cannot save '%s' - file is read-only", filename)
	}
	
	return baseMsg
}

// FormatLoadError formats load-specific errors
func FormatLoadError(filename string, err error) string {
	if err == nil {
		return ""
	}
	
	baseMsg := FormatUserError(err)
	
	// Add context for load operations
	if baseMsg == "File not found" {
		return fmt.Sprintf("'%s' not found - creating new file", filename)
	}
	
	if strings.Contains(baseMsg, "directory") {
		return fmt.Sprintf("'%s' is a directory, not a file", filename)
	}
	
	return fmt.Sprintf("Cannot open '%s': %s", filename, baseMsg)
}
