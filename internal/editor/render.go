package editor

import (
	filebrowser "github.com/Adelodunpeter25/vx/internal/file-browser"
	splitpane "github.com/Adelodunpeter25/vx/internal/split-pane"
	"github.com/gdamore/tcell/v2"
)

func (e *Editor) render() {
	e.term.Clear()

	contentHeight := e.height - 1
	cdPromptActive := e.active() != nil && e.active().mode == ModeCdPrompt
	cdPromptRows := 0
	if cdPromptActive {
		cdPromptRows = 4
		if contentHeight > cdPromptRows {
			contentHeight -= cdPromptRows
		} else {
			contentHeight = 1
		}
	}
	contentX := 0
	contentWidth := e.width
	if e.fileBrowser != nil && e.fileBrowser.Open {
		if e.fileBrowser.Width < 10 {
			e.fileBrowser.Width = 10
		}
		if e.fileBrowser.Width > e.width-10 {
			e.fileBrowser.Width = e.width - 10
		}
		e.fileBrowser.Render(e.term, 0, 0, e.fileBrowser.Width, contentHeight)
		contentX = e.fileBrowser.Width + 1
		contentWidth = e.width - contentX
	}
	rects, dividerX := splitpane.LayoutSideBySide(contentWidth, contentHeight, len(e.panes), e.splitRatio)
	for i, rect := range rects {
		if i < len(e.panes) {
			isActive := i == e.activePane
			rect.X += contentX
			e.renderPane(e.panes[i], rect, isActive)
		}
	}

	// Draw divider for 2-pane layout
	if dividerX >= 0 {
		style := tcell.StyleDefault.Foreground(tcell.ColorGray)
		if e.dragSplit {
			style = style.Bold(true)
		}
		for y := 0; y < contentHeight; y++ {
			e.term.SetCell(contentX+dividerX, y, '│', style)
		}
	}
	if e.fileBrowser != nil && e.fileBrowser.Open {
		dividerStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite)
		if e.dragBrowser {
			dividerStyle = dividerStyle.Bold(true)
		}
		dividerX := e.fileBrowser.Width
		if dividerX >= 0 && dividerX < e.width {
			for y := 0; y < contentHeight; y++ {
				e.term.SetCell(dividerX, y, '│', dividerStyle)
			}
		}
	}

	e.renderStatusLine()
	if cdPromptActive {
		promptY := e.height - 1 - cdPromptRows
		if promptY < 0 {
			promptY = 0
		}
		filebrowser.RenderCdPrompt(e.term, e.cdPrompt, e.width, promptY, promptY+1, cdPromptRows-1)
	}
	if e.palette != nil && e.palette.Active {
		e.palette.Render(e.term, e.width, e.height)
	}
	e.term.Show()
}

