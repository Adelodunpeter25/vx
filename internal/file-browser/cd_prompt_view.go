package filebrowser

import (
	"strings"

	"github.com/Adelodunpeter25/vx/internal/terminal"
	"github.com/gdamore/tcell/v2"
)

func RenderCdPrompt(term *terminal.Terminal, prompt *CdPrompt, width int, promptY int, suggestY int) {
	if term == nil || prompt == nil || width <= 0 {
		return
	}
	promptStyle := tcell.StyleDefault.Foreground(tcell.ColorBlack).Background(tcell.ColorLightBlue)
	suggestStyle := tcell.StyleDefault.Foreground(tcell.ColorBlack).Background(tcell.ColorSilver)

	for x := 0; x < width; x++ {
		term.SetCell(x, promptY, ' ', promptStyle)
		term.SetCell(x, suggestY, ' ', suggestStyle)
	}

	line, cursorX := prompt.RenderLine(width)
	term.DrawText(0, promptY, line, promptStyle)
	if cursorX >= 0 && cursorX < width {
		r, _, style, _ := term.ScreenContent(cursorX, promptY)
		term.SetCell(cursorX, promptY, r, style.Reverse(true))
	}

	suggestions := prompt.Suggestions()
	if len(suggestions) == 0 {
		return
	}
	text := strings.Join(suggestions, " ")
	if len(text) > width {
		text = text[:width]
	}
	term.DrawText(0, suggestY, padRight(text, width), suggestStyle)
}
