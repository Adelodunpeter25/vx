package palette

import (
	"github.com/Adelodunpeter25/vx/internal/terminal"
	"github.com/gdamore/tcell/v2"
)

type Item struct {
	Label string
	Data  interface{}
}

type Palette struct {
	Active      bool
	Input       string
	Items       []Item
	Filtered    []Item
	SelectedIndex int
	Prompt      string
	OnSelect    func(Item)
	OnCancel    func()
	OnChange    func(string)
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

	// Overlay render at the bottom
	pWidth := width
	pHeight := 10
	if pHeight > height {
		pHeight = height
	}
	x := 0
	y := height - pHeight

	style := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorWhite)
	
	// Draw box
	for i := 0; i < pWidth; i++ {
		for j := 0; j < pHeight; j++ {
			term.SetCell(x+i, y+j, ' ', style)
		}
	}

	// Render input line
	inputLine := p.Prompt + " " + p.Input
	term.DrawText(x+1, y+1, inputLine, style.Bold(true))

	// Render results
	for i := 0; i < pHeight-3 && i < len(p.Filtered); i++ {
		itemStyle := style
		if i == p.SelectedIndex {
			itemStyle = itemStyle.Reverse(true)
		}
		term.DrawText(x+1, y+3+i, p.Filtered[i].Label, itemStyle)
	}
}
