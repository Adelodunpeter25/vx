package terminal

import (
	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
)

func (t *Terminal) SetCell(x, y int, r rune, style tcell.Style) {
	t.screen.SetContent(x, y, r, nil, style)
}

func (t *Terminal) DrawText(x, y int, text string, style tcell.Style) {
	cursorX := x
	for _, r := range text {
		t.SetCell(cursorX, y, r, style)
		cursorX += runewidth.RuneWidth(r)
	}
}

func (t *Terminal) DrawLine(y int, text string, style tcell.Style) {
	t.DrawText(0, y, text, style)
}

func (t *Terminal) ScreenContent(x, y int) (rune, []rune, tcell.Style, int) {
	return t.screen.GetContent(x, y)
}
