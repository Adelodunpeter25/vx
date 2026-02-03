package terminal

import "github.com/gdamore/tcell/v2"

func (t *Terminal) SetCell(x, y int, r rune, style tcell.Style) {
	t.screen.SetContent(x, y, r, nil, style)
}

func (t *Terminal) DrawText(x, y int, text string, style tcell.Style) {
	for i, r := range text {
		t.SetCell(x+i, y, r, style)
	}
}

func (t *Terminal) DrawLine(y int, text string, style tcell.Style) {
	t.DrawText(0, y, text, style)
}

func (t *Terminal) ScreenContent(x, y int) (rune, []rune, tcell.Style, int) {
	return t.screen.GetContent(x, y)
}
