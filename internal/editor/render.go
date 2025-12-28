package editor

import "github.com/gdamore/tcell/v2"

func (e *Editor) render() {
	e.term.Clear()
	
	contentHeight := e.height - 1
	for i := 0; i < contentHeight; i++ {
		lineNum := e.offsetY + i
		if lineNum >= e.buffer.LineCount() {
			e.term.DrawText(0, i, "~", tcell.StyleDefault.Foreground(tcell.ColorBlue))
		} else {
			line := e.buffer.Line(lineNum)
			e.term.DrawText(0, i, line, tcell.StyleDefault)
		}
	}
	
	e.renderStatusLine()
	
	screenY := e.cursorY - e.offsetY
	e.term.SetCell(e.cursorX, screenY, ' ', tcell.StyleDefault.Reverse(true))
	
	e.term.Show()
}

func (e *Editor) renderStatusLine() {
	y := e.height - 1
	style := tcell.StyleDefault.Reverse(true)
	
	for x := 0; x < e.width; x++ {
		e.term.SetCell(x, y, ' ', style)
	}
	
	if e.mode == ModeCommand {
		cmd := ":" + e.commandBuf
		e.term.DrawText(0, y, cmd, style)
		return
	}
	
	if e.message != "" {
		e.term.DrawText(0, y, e.message, style)
		return
	}
	
	mode := e.mode.String()
	e.term.DrawText(0, y, " "+mode+" ", style)
	
	filename := e.buffer.Filename()
	if filename == "" {
		filename = "[No Name]"
	}
	modified := ""
	if e.buffer.IsModified() {
		modified = " [+]"
	}
	info := filename + modified
	e.term.DrawText(len(mode)+2, y, info, style)
	
	pos := " " + string(rune('0'+e.cursorY+1)) + "," + string(rune('0'+e.cursorX+1)) + " "
	e.term.DrawText(e.width-len(pos), y, pos, style)
}
