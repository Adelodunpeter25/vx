package osc

import (
	"path/filepath"
	"strings"
)

// Title returns the terminal title for the current editor context.
// If fileName is empty, the title falls back to "vx - {dir}".
// If fileName is set, the title becomes "{file} - {dir} - vx".
func Title(fileName, dirPath string) string {
	dir := displayName(dirPath)
	if fileName == "" {
		return "vx - " + dir
	}
	file := displayName(fileName)
	return file + " - " + dir + " - vx"
}

func displayName(path string) string {
	if path == "" {
		return "vx"
	}
	cleaned := filepath.Clean(path)
	base := filepath.Base(cleaned)
	if base == "." || base == string(filepath.Separator) || base == "" {
		base = cleaned
	}
	if base == "" {
		return "vx"
	}
	return strings.TrimSpace(base)
}
