package utils

import "fmt"

// FileError represents a file-related error with context
type FileError struct {
	Filename string
	Op       string
	Err      error
}

func (e *FileError) Error() string {
	return fmt.Sprintf("%s: %s: %v", e.Op, e.Filename, e.Err)
}

// NewFileError creates a new file error
func NewFileError(op, filename string, err error) *FileError {
	return &FileError{
		Filename: filename,
		Op:       op,
		Err:      err,
	}
}

// IsRecoverable checks if an error is recoverable
func IsRecoverable(err error) bool {
	if err == nil {
		return true
	}
	
	// Check for common recoverable errors
	switch err.Error() {
	case "file too large":
		return false
	case "too many lines":
		return false
	default:
		return true
	}
}
