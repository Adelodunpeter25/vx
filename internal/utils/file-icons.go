package utils

import (
	"path/filepath"
	"strings"
)

func FileIcon(path string) string {
	return FileIconFor(path, false)
}

func FileIconFor(path string, isDir bool) string {
	if path == "" {
		return ""
	}

	if isDir {
		return "󰉋"
	}

	switch strings.ToLower(filepath.Ext(path)) {
	case ".go":
		return ""
	case ".md":
		return "󰍔"
	case ".json":
		return ""
	case ".js", ".mjs", ".cjs":
		return ""
	case ".ts", ".tsx":
		return ""
	case ".py":
		return ""
	case ".rs":
		return ""
	case ".toml":
		return ""
	case ".yaml", ".yml":
		return ""
	case ".txt":
		return "󰈙"
	default:
		return ""
	}
}

func DirIcon(expanded bool) string {
	if expanded {
		return ""
	}
	return "󰉋"
}
