package editor

import (
	"fmt"
	"strings"

	"github.com/Adelodunpeter25/vx/internal/replace"
	"github.com/gdamore/tcell/v2"
)

func (e *Editor) renderStatusLine() {
	p := e.active()
	y := e.height - 1
	style := tcell.StyleDefault.Reverse(true)

	// Clear status line
	for x := 0; x < e.width; x++ {
		e.term.SetCell(x, y, ' ', style)
	}

	// Handle special modes
	if p.mode == ModeCommand {
		e.renderCommandStatus(y, style)
		return
	}

	if p.mode == ModeSearch {
		e.renderSearchStatus(y, style)
		return
	}

	if p.mode == ModeReplace {
		e.renderReplaceStatus(y, style)
		return
	}

	// Normal status line
	e.renderNormalStatus(y, style)
}

func (e *Editor) renderCommandStatus(y int, style tcell.Style) {
	cmd := ":" + e.active().commandBuf
	e.term.DrawText(0, y, cmd, style)
}

func (e *Editor) renderSearchStatus(y int, style tcell.Style) {
	search := "/" + e.active().searchBuf
	e.term.DrawText(0, y, search, style)
}

func (e *Editor) renderReplaceStatus(y int, style tcell.Style) {
	state := e.active().replace.GetState()

	switch state {
	case replace.StateSearchInput:
		prompt := "Find: " + e.active().replace.GetSearchTerm()
		e.term.DrawText(0, y, prompt, style)
	case replace.StateReplaceInput:
		prompt := fmt.Sprintf("Find: %s | Replace: %s", e.active().replace.GetSearchTerm(), e.active().replace.GetReplaceTerm())
		e.term.DrawText(0, y, prompt, style)
	case replace.StateConfirm:
		prompt := fmt.Sprintf("Replace? [y/n/q] (%d/%d)", e.active().replace.GetCurrentIndex(), e.active().replace.GetMatchCount())
		e.term.DrawText(0, y, prompt, style)
	}
}

func (e *Editor) renderNormalStatus(y int, style tcell.Style) {
	// Show mode
	var mode string
	if e.active().preview.IsEnabled() {
		mode = "PREVIEW"
	} else {
		mode = e.active().mode.String()
	}
	e.term.DrawText(0, y, " "+mode+" ", style)
	modeWidth := len(mode) + 2

	// Show message if present
	message := e.active().msgManager.Get()
	if message != "" {
		e.renderStatusMessage(y, style, modeWidth, message)
		return
	}

	// Show file info
	e.renderFileInfo(y, style, modeWidth)
}

func (e *Editor) renderStatusMessage(y int, style tcell.Style, modeWidth int, message string) {
	// Check if message is a file info message (contains KB/MB and "lines")
	if strings.Contains(message, " lines") && (strings.Contains(message, "KB") || strings.Contains(message, "MB") || strings.Contains(message, "GB") || strings.Contains(message, " B,")) {
		e.renderFileInfoMessage(y, style, modeWidth, message)
	} else {
		e.term.DrawText(modeWidth+1, y, message, style)
	}
}

func (e *Editor) renderFileInfo(y int, style tcell.Style, modeWidth int) {
	p := e.active()
	filename := p.buffer.Filename()
	if filename == "" {
		filename = "[No Name]"
	}
	modified := ""
	if p.buffer.IsModified() {
		modified = " [+]"
	}
	info := filename + modified
	e.term.DrawText(modeWidth+1, y, info, style)

	// Show buffer count and cursor position
	e.renderRightInfo(y, style)
}

func (e *Editor) renderRightInfo(y int, style tcell.Style) {
	// Show buffer count if multiple buffers
	if e.active().bufferMgr.Count() > 1 {
		p := e.active()
		bufInfo := fmt.Sprintf(" Buffer %d/%d ", p.bufferMgr.CurrentIndex(), p.bufferMgr.Count())
		bufInfoX := e.width - len(bufInfo)

		// Don't show cursor position in preview mode
		if !p.preview.IsEnabled() {
			pos := fmt.Sprintf(" %d,%d ", p.cursorY+1, p.cursorX+1)
			bufInfoX -= len(pos)
			e.term.DrawText(e.width-len(pos), y, pos, style)
		}

		e.term.DrawText(bufInfoX, y, bufInfo, style)
	} else {
		// Don't show cursor position in preview mode
		p := e.active()
		if !p.preview.IsEnabled() {
			pos := fmt.Sprintf(" %d,%d ", p.cursorY+1, p.cursorX+1)
			e.term.DrawText(e.width-len(pos), y, pos, style)
		}
	}
}

func (e *Editor) renderFileInfoMessage(y int, style tcell.Style, modeWidth int, message string) {
	// Parse message: "filename" size, lines
	parts := strings.SplitN(message, "\"", 3)
	if len(parts) < 3 {
		e.term.DrawText(modeWidth+1, y, message, style)
		return
	}

	filename := parts[1]
	rest := strings.TrimSpace(parts[2])

	// Draw filename after mode
	e.term.DrawText(modeWidth+1, y, "\""+filename+"\"", style)

	// Draw size and lines on right
	e.term.DrawText(e.width-len(rest)-1, y, rest, style)
}
