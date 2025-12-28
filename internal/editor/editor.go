package editor

import (
	"github.com/Adelodunpeter25/vx/internal/buffer"
	"github.com/Adelodunpeter25/vx/internal/terminal"
	"github.com/gdamore/tcell/v2"
)

type Editor struct {
	term     *terminal.Terminal
	buffer   *buffer.Buffer
	width    int
	height   int
	cursorX  int
	cursorY  int
	offsetY  int
	mode     Mode
	quit     bool
}

func New(term *terminal.Terminal) *Editor {
	width, height := term.Size()
	return &Editor{
		term:    term,
		buffer:  buffer.New(),
		width:   width,
		height:  height,
		mode:    ModeNormal,
	}
}

func NewWithFile(term *terminal.Terminal, filename string) (*Editor, error) {
	buf, err := buffer.Load(filename)
	if err != nil {
		return nil, err
	}
	
	width, height := term.Size()
	return &Editor{
		term:    term,
		buffer:  buf,
		width:   width,
		height:  height,
		mode:    ModeNormal,
	}, nil
}

func (e *Editor) Run() error {
	e.render()
	
	for !e.quit {
		e.handleEvent()
	}
	return nil
}

func (e *Editor) handleEvent() {
	ev := e.term.ReadEvent()
	if ev == nil {
		return
	}

	switch ev.Type {
	case terminal.EventKey:
		e.handleKey(ev)
	case terminal.EventResize:
		e.width, e.height = e.term.Size()
		e.render()
	}
}

func (e *Editor) handleKey(ev *terminal.Event) {
	switch e.mode {
	case ModeNormal:
		e.handleNormalMode(ev)
	case ModeInsert:
		e.handleInsertMode(ev)
	case ModeCommand:
		e.handleCommandMode(ev)
	}
	e.render()
}

func (e *Editor) handleNormalMode(ev *terminal.Event) {
	switch ev.Rune {
	case 'q':
		e.quit = true
	case 'i':
		e.mode = ModeInsert
	case 'h':
		if e.cursorX > 0 {
			e.cursorX--
		}
	case 'j':
		if e.cursorY < e.buffer.LineCount()-1 {
			e.cursorY++
			e.adjustScroll()
			e.clampCursor()
		}
	case 'k':
		if e.cursorY > 0 {
			e.cursorY--
			e.adjustScroll()
			e.clampCursor()
		}
	case 'l':
		line := e.buffer.Line(e.cursorY)
		if e.cursorX < len(line) {
			e.cursorX++
		}
	}
	
	// Arrow keys
	switch ev.Key {
	case tcell.KeyLeft:
		if e.cursorX > 0 {
			e.cursorX--
		}
	case tcell.KeyRight:
		line := e.buffer.Line(e.cursorY)
		if e.cursorX < len(line) {
			e.cursorX++
		}
	case tcell.KeyUp:
		if e.cursorY > 0 {
			e.cursorY--
			e.adjustScroll()
			e.clampCursor()
		}
	case tcell.KeyDown:
		if e.cursorY < e.buffer.LineCount()-1 {
			e.cursorY++
			e.adjustScroll()
			e.clampCursor()
		}
	}
}

func (e *Editor) handleInsertMode(ev *terminal.Event) {
	if ev.Key == tcell.KeyEscape {
		e.mode = ModeNormal
		if e.cursorX > 0 {
			e.cursorX--
		}
		return
	}
	
	if ev.Key == tcell.KeyEnter {
		e.buffer.SplitLine(e.cursorY, e.cursorX)
		e.cursorY++
		e.cursorX = 0
		e.adjustScroll()
		return
	}
	
	if ev.Key == tcell.KeyBackspace || ev.Key == tcell.KeyBackspace2 {
		if e.cursorX > 0 {
			e.buffer.DeleteRune(e.cursorY, e.cursorX)
			e.cursorX--
		} else if e.cursorY > 0 {
			prevLen := len(e.buffer.Line(e.cursorY - 1))
			e.buffer.JoinLine(e.cursorY - 1)
			e.cursorY--
			e.cursorX = prevLen
			e.adjustScroll()
		}
		return
	}
	
	// Arrow keys in insert mode
	switch ev.Key {
	case tcell.KeyLeft:
		if e.cursorX > 0 {
			e.cursorX--
		}
		return
	case tcell.KeyRight:
		line := e.buffer.Line(e.cursorY)
		if e.cursorX < len(line) {
			e.cursorX++
		}
		return
	case tcell.KeyUp:
		if e.cursorY > 0 {
			e.cursorY--
			e.adjustScroll()
			e.clampCursor()
		}
		return
	case tcell.KeyDown:
		if e.cursorY < e.buffer.LineCount()-1 {
			e.cursorY++
			e.adjustScroll()
			e.clampCursor()
		}
		return
	}
	
	if ev.Rune != 0 {
		e.buffer.InsertRune(e.cursorY, e.cursorX, ev.Rune)
		e.cursorX++
	}
}

func (e *Editor) handleCommandMode(ev *terminal.Event) {
	// TODO: implement command mode
}

func (e *Editor) clampCursor() {
	line := e.buffer.Line(e.cursorY)
	maxX := len(line)
	if e.mode == ModeNormal && maxX > 0 {
		maxX--
	}
	if e.cursorX > maxX {
		e.cursorX = maxX
	}
	if e.cursorX < 0 {
		e.cursorX = 0
	}
}

func (e *Editor) adjustScroll() {
	contentHeight := e.height - 1
	if e.cursorY < e.offsetY {
		e.offsetY = e.cursorY
	}
	if e.cursorY >= e.offsetY+contentHeight {
		e.offsetY = e.cursorY - contentHeight + 1
	}
}

func (e *Editor) render() {
	e.term.Clear()
	
	// Render buffer content
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
	
	// Render status line
	e.renderStatusLine()
	
	// Position cursor
	screenY := e.cursorY - e.offsetY
	e.term.SetCell(e.cursorX, screenY, ' ', tcell.StyleDefault.Reverse(true))
	
	e.term.Show()
}

func (e *Editor) renderStatusLine() {
	y := e.height - 1
	style := tcell.StyleDefault.Reverse(true)
	
	// Clear status line
	for x := 0; x < e.width; x++ {
		e.term.SetCell(x, y, ' ', style)
	}
	
	// Mode indicator
	mode := e.mode.String()
	e.term.DrawText(0, y, " "+mode+" ", style)
	
	// File info
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
	
	// Cursor position
	pos := ""
	pos = " " + string(rune('0'+e.cursorY+1)) + "," + string(rune('0'+e.cursorX+1)) + " "
	e.term.DrawText(e.width-len(pos), y, pos, style)
}
