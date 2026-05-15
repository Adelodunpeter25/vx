package editor

import (
	"github.com/Adelodunpeter25/vx/internal/palette"
	"github.com/Adelodunpeter25/vx/internal/terminal"
)

func (e *Editor) openPalette() {
	e.palette = palette.New("")
	e.palette.Active = true
	e.active().mode = ModePalette

	// Initial file search
	items := palette.SearchFiles(".", "")
	e.palette.SetItems(items)

	e.palette.OnChange = func(input string) {
		items := palette.SearchFiles(".", input)
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
