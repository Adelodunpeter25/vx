package palette

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/Adelodunpeter25/vx/internal/utils"
)

var (
	fileCache      []string
	cacheRoot      string
	cacheTimestamp int64
)

// ListAllFiles returns all files in root, using a cache if available
func ListAllFiles(root string) []string {
	// Simple cache invalidation: if root changes, clear cache
	if root != cacheRoot {
		fileCache = nil
		cacheRoot = root
	}

	if fileCache != nil {
		return fileCache
	}

	var files []string
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			if utils.ShouldSkipDir(info.Name()) {
				return filepath.SkipDir
			}
			return nil
		}

		rel, err := filepath.Rel(root, path)
		if err != nil {
			rel = path
		}
		files = append(files, rel)
		return nil
	})

	fileCache = files
	return files
}

func SearchFiles(root string, pattern string) []Item {
	var items []Item
	pattern = strings.ToLower(pattern)
	files := ListAllFiles(root)

	hasWildcard := strings.Contains(pattern, "*") || strings.Contains(pattern, "?") || strings.Contains(pattern, "[")

	for _, rel := range files {
		path := filepath.Join(root, rel)
		relLower := strings.ToLower(rel)
		
		matched := false
		if pattern == "" {
			matched = true
		} else if hasWildcard {
			matched, _ = filepath.Match(pattern, relLower)
		} else {
			matched = strings.Contains(relLower, pattern)
		}

		if matched {
			items = append(items, Item{
				Label: rel,
				Data:  path,
				Icon:  utils.FileIcon(path),
			})
		}

		if len(items) > 100 {
			break
		}
	}

	return items
}
