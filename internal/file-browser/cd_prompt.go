package filebrowser

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/Adelodunpeter25/vx/internal/terminal"
	"github.com/gdamore/tcell/v2"
)

type CdAction struct {
	Apply  bool
	Cancel bool
	Path   string
}

type CdPrompt struct {
	Value       string
	suggestions []string
}

func NewCdPrompt(initial string) *CdPrompt {
	if initial == "" {
		if cwd, err := os.Getwd(); err == nil {
			initial = cwd
		}
	}
	return &CdPrompt{Value: initial}
}

func (c *CdPrompt) HandleKey(ev *terminal.Event) CdAction {
	if ev == nil {
		return CdAction{}
	}
	switch ev.Key {
	case tcell.KeyEscape:
		return CdAction{Cancel: true}
	case tcell.KeyEnter:
		return CdAction{Apply: true, Path: c.Value}
	case tcell.KeyBackspace, tcell.KeyBackspace2:
		c.Value = deletePathSegment(c.Value)
		c.suggestions = nil
		return CdAction{}
	case tcell.KeyTab:
		c.complete()
		return CdAction{}
	}
	if ev.Rune != 0 {
		c.Value += string(ev.Rune)
		c.suggestions = nil
	}
	return CdAction{}
}

func (c *CdPrompt) Render(width int) string {
	if width <= 0 {
		return ""
	}
	prefix := " cd " + c.Value
	if len(c.suggestions) > 0 {
		prefix += "  [" + strings.Join(c.suggestions, " ") + "]"
	}
	if len(prefix) > width {
		return prefix[:width]
	}
	return padRight(prefix, width)
}

func (c *CdPrompt) complete() {
	base, partial, displayBase := splitPath(c.Value)
	entries, err := os.ReadDir(base)
	if err != nil {
		c.suggestions = nil
		return
	}
	var matches []string
	for _, ent := range entries {
		if !ent.IsDir() {
			continue
		}
		name := ent.Name()
		if strings.HasPrefix(name, partial) {
			matches = append(matches, name)
		}
	}
	if len(matches) == 0 {
		c.suggestions = nil
		return
	}
	if len(matches) == 1 {
		c.suggestions = nil
		c.Value = filepath.Join(displayBase, matches[0])
		return
	}
	c.suggestions = matches
}

func splitPath(input string) (base string, partial string, displayBase string) {
	if input == "" {
		return ".", "", "."
	}
	expanded := ExpandHome(input)
	if strings.HasSuffix(expanded, string(os.PathSeparator)) {
		return expanded, "", input
	}
	base = filepath.Dir(expanded)
	partial = filepath.Base(expanded)
	displayBase = filepath.Dir(input)
	if displayBase == "." {
		displayBase = base
	}
	return base, partial, displayBase
}

func ExpandHome(path string) string {
	if path == "~" || strings.HasPrefix(path, "~/") {
		if home, err := os.UserHomeDir(); err == nil {
			return filepath.Join(home, strings.TrimPrefix(path, "~"))
		}
	}
	return path
}

func deletePathSegment(path string) string {
	if path == "" {
		return ""
	}
	trimmed := strings.TrimRight(path, string(os.PathSeparator))
	if trimmed == "" {
		return string(os.PathSeparator)
	}
	dir := filepath.Dir(trimmed)
	if dir == "." {
		return ""
	}
	if dir == string(os.PathSeparator) {
		return dir
	}
	return dir + string(os.PathSeparator)
}
