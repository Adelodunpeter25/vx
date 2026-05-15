package palette

import (
	"os"
	"path/filepath"
	"strings"
)

func SearchFiles(root string, pattern string) []Item {
	var items []Item
	pattern = strings.ToLower(pattern)

	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			if strings.HasPrefix(info.Name(), ".") && info.Name() != "." {
				return filepath.SkipDir
			}
			return nil
		}

		rel, err := filepath.Rel(root, path)
		if err != nil {
			rel = path
		}

		if pattern == "" || strings.Contains(strings.ToLower(rel), pattern) {
			items = append(items, Item{
				Label: rel,
				Data:  path,
				Icon:  fileIcon(path),
			})
		}
		
		// Limit results for performance
		if len(items) > 100 {
			return filepath.SkipAll
		}

		return nil
	})

	return items
}

func fileIcon(path string) string {
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
	default:
		return ""
	}
}
