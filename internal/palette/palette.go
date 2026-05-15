package palette

import (
	"strings"

	"github.com/Adelodunpeter25/vx/internal/terminal"
	"github.com/gdamore/tcell/v2"
)

type Item struct {
	Label string
	Data  interface{}
}

type Palette struct {
	Active        bool
	Input         string
	Items         []Item
	Filtered      []Item
	SelectedIndex int
	Prompt        string
	OnSelect      func(Item)
	OnCancel      func()
	OnChange      func(string)
}

func New(prompt string) *Palette {
	return &Palette{
		Prompt: prompt,
	}
}

func (p *Palette) SetItems(items []Item) {
	p.Items = items
	p.filter()
}

func (p *Palette) filter() {
	// Simple prefix filter for now, can be improved to fuzzy later
	p.Filtered = nil
	for _, item := range p.Items {
		p.Filtered = append(p.Filtered, item)
	}
	if p.SelectedIndex >= len(p.Filtered) {
		p.SelectedIndex = 0
	}
}

func (p *Palette) HandleKey(ev *terminal.Event) {
	switch ev.Key {
	case tcell.KeyEscape:
		if p.OnCancel != nil {
			p.OnCancel()
		}
	case tcell.KeyEnter:
		if len(p.Filtered) > 0 && p.SelectedIndex >= 0 && p.SelectedIndex < len(p.Filtered) {
			if p.OnSelect != nil {
				p.OnSelect(p.Filtered[p.SelectedIndex])
			}
		}
	case tcell.KeyUp:
		if p.SelectedIndex > 0 {
			p.SelectedIndex--
		}
	case tcell.KeyDown:
		if p.SelectedIndex < len(p.Filtered)-1 {
			p.SelectedIndex++
		}
	case tcell.KeyBackspace, tcell.KeyBackspace2:
		if len(p.Input) > 0 {
			p.Input = p.Input[:len(p.Input)-1]
			p.filter()
			if p.OnChange != nil {
				p.OnChange(p.Input)
			}
		}
	default:
		if ev.Rune != 0 {
			p.Input += string(ev.Rune)
			p.filter()
			if p.OnChange != nil {
				p.OnChange(p.Input)
			}
		}
	}
}

func (p *Palette) Render(term *terminal.Terminal, width, height int) {
	if !p.Active {
		return
	}

	numResults := len(p.Filtered)
	maxResults := 9 // Max results to show
	if numResults > maxResults {
		numResults = maxResults
	}

	// Dynamic height: results + separator + input
	pHeight := numResults + 2
	if pHeight > height {
		pHeight = height
	}
	x := 0
	y := height - pHeight

	bgStyle := tcell.StyleDefault.Background(tcell.NewRGBColor(30, 30, 30)).Foreground(tcell.ColorWhite)
	
	// Draw background
	for i := 0; i < width; i++ {
		for j := 0; j < pHeight; j++ {
			term.SetCell(x+i, y+j, ' ', bgStyle)
		}
	}

	// Render results (top part of palette)
	for i := 0; i < numResults; i++ {
		itemStyle := bgStyle
		if i == p.SelectedIndex {
			itemStyle = itemStyle.Background(tcell.NewRGBColor(60, 60, 100)).Bold(true)
		}
		label := " " + p.Filtered[i].Label
		term.DrawText(x, y+i, padRight(label, width), itemStyle)
	}

	// Separator line
	sepY := y + pHeight - 2
	sepStyle := tcell.StyleDefault.Foreground(tcell.ColorGray)
	if pHeight >= 2 {
		for i := 0; i < width; i++ {
			term.SetCell(x+i, sepY, '─', sepStyle)
		}
	}

	// Render input line at the very bottom
	inputY := y + pHeight - 1
	inputStyle := bgStyle.Bold(true)
	prompt := "> "
	if p.Prompt != "" {
		prompt = p.Prompt + " "
	}
	term.DrawText(x, inputY, prompt+p.Input, inputStyle)
	
	// Show blinking cursor
	term.ShowCursor(x+len(prompt)+len(p.Input), inputY)
}

func padRight(s string, width int) string {
	if len(s) >= width {
		return s
	}
	return s + strings.Repeat(" ", width-len(s))
}
