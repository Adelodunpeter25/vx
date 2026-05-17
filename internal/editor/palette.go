package editor

import (
	"github.com/Adelodunpeter25/vx/internal/palette"
	"github.com/Adelodunpeter25/vx/internal/terminal"
)

func (e *Editor) openPalette() {
	root := e.launchDir
	if e.fileBrowser != nil && e.fileBrowser.RootPath != "" {
		root = e.fileBrowser.RootPath
	}
	if root == "" {
		root = "."
	}

	e.palette = palette.New("")
	e.palette.Active = true
	e.active().mode = ModePalette

	// Initial file search
	items := palette.SearchFiles(root, "")
	e.palette.SetItems(items)

	e.palette.OnChange = func(input string) {
		items := palette.SearchFiles(root, input)
		e.palette.SetItems(items)
	}

	e.palette.OnSelect = func(item palette.Item) {
		path := item.Data.(string)
		e.openFileInActivePane(path)
		e.closePalette()
	}

	e.palette.OnCancel = func() {
		e.closePalette()
	}
}

func (e *Editor) closePalette() {
	if e.palette != nil {
		e.palette.Active = false
	}
	e.active().mode = ModeNormal
}

func (e *Editor) handlePaletteEvent(ev *terminal.Event) {
	if e.palette == nil || !e.palette.Active {
		e.active().mode = ModeNormal
		return
	}
	e.palette.HandleKey(ev)
}
