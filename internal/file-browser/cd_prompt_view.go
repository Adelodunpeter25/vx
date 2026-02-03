package filebrowser

import (
	"strings"

	"github.com/Adelodunpeter25/vx/internal/terminal"
	"github.com/gdamore/tcell/v2"
)

func RenderCdPrompt(term *terminal.Terminal, prompt *CdPrompt, width int, promptY int, suggestY int, suggestRows int) {
	if term == nil || prompt == nil || width <= 0 {
		return
	}
	promptStyle := tcell.StyleDefault.Foreground(tcell.ColorGray).Background(tcell.ColorBlack)
	suggestStyle := tcell.StyleDefault.Foreground(tcell.ColorGray).Background(tcell.ColorBlack)

	for x := 0; x < width; x++ {
		term.SetCell(x, promptY, ' ', promptStyle)
		for row := 0; row < suggestRows; row++ {
			term.SetCell(x, suggestY+row, ' ', suggestStyle)
		}
	}

	line, cursorX := prompt.RenderLine(width)
	term.DrawText(0, promptY, line, promptStyle)
	if cursorX >= 0 && cursorX < width {
		r, _, style, _ := term.ScreenContent(cursorX, promptY)
		term.SetCell(cursorX, promptY, r, style.Reverse(true))
	}

	suggestions := prompt.Suggestions()
	if len(suggestions) == 0 || suggestRows <= 0 {
		return
	}
	text := strings.Join(suggestions, " ")
	lines := wrapText(text, width, suggestRows)
	for i, line := range lines {
		term.DrawText(0, suggestY+i, padRight(line, width), suggestStyle)
	}
}

func wrapText(text string, width int, rows int) []string {
	if width <= 0 || rows <= 0 || text == "" {
		return nil
	}
	var lines []string
	for len(text) > 0 && len(lines) < rows {
		if len(text) <= width {
			lines = append(lines, text)
			break
		}
		lines = append(lines, text[:width])
		text = text[width:]
	}
	return lines
}
