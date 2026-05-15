package palette

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/Adelodunpeter25/vx/internal/utils"
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
				Icon:  utils.FileIcon(path),
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
