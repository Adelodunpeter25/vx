package utils

import "fmt"

// FormatFileSize formats bytes into human-readable size
func FormatFileSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	
	units := []string{"KB", "MB", "GB", "TB"}
	return fmt.Sprintf("%.1f %s", float64(bytes)/float64(div), units[exp])
}

// FormatLineCount formats line count with proper pluralization
func FormatLineCount(count int) string {
	if count == 1 {
		return "1 line"
	}
	return fmt.Sprintf("%d lines", count)
}

// FormatFileInfo creates a file info message
func FormatFileInfo(filename string, size int64, lines int) string {
	return fmt.Sprintf("\"%s\" %s, %s", filename, FormatFileSize(size), FormatLineCount(lines))
}
