package utils

import (
	"path/filepath"
	"strings"
)

func FileIcon(path string) string {
	if path == "" {
		return ""
	}

	if strings.HasSuffix(path, string(filepath.Separator)) {
		return ""
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
