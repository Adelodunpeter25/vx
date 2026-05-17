package palette

import (
	"path/filepath"
	"sort"
	"strings"

	"github.com/Adelodunpeter25/vx/internal/terminal"
	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
)

type Item struct {
	Label string
	Data  interface{}
	Icon  string
}

const maxPaletteResults = 9

type Palette struct {
	Active        bool
	Input         string
	Items         []Item
	Filtered      []Item
	SelectedIndex int
	scrollOffset  int
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
	type scoredItem struct {
		item  Item
		score int
	}

	filtered := make([]scoredItem, 0, len(p.Items))
	query := strings.ToLower(strings.TrimSpace(p.Input))
	for _, item := range p.Items {
		score, ok := scoreItem(item, query)
		if ok {
			filtered = append(filtered, scoredItem{item: item, score: score})
		}
	}
	sort.SliceStable(filtered, func(i, j int) bool {
		if filtered[i].score != filtered[j].score {
			return filtered[i].score < filtered[j].score
		}
		return strings.ToLower(filtered[i].item.Label) < strings.ToLower(filtered[j].item.Label)
	})

	p.Filtered = p.Filtered[:0]
	for _, scored := range filtered {
		p.Filtered = append(p.Filtered, scored.item)
	}
	if p.SelectedIndex >= len(p.Filtered) {
		p.SelectedIndex = 0
	}
	if p.SelectedIndex < 0 {
		p.SelectedIndex = 0
	}
	p.scrollOffset = 0
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
		if p.SelectedIndex < p.scrollOffset {
			p.scrollOffset = p.SelectedIndex
		}
	case tcell.KeyDown:
		if p.SelectedIndex < len(p.Filtered)-1 {
			p.SelectedIndex++
		}
		if p.SelectedIndex >= p.scrollOffset+maxPaletteResults {
			p.scrollOffset = p.SelectedIndex - maxPaletteResults + 1
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
	visibleResults := numResults
	if visibleResults > maxPaletteResults {
		visibleResults = maxPaletteResults
	}

	// Dynamic height: results + separator + input
	pHeight := visibleResults + 2
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
	for i := 0; i < visibleResults; i++ {
		idx := p.scrollOffset + i
		if idx >= len(p.Filtered) {
			break
		}
		itemStyle := bgStyle
		if idx == p.SelectedIndex {
			itemStyle = itemStyle.Background(tcell.NewRGBColor(60, 60, 100)).Bold(true)
		}
		iconStyle := itemStyle.Foreground(tcell.NewRGBColor(250, 179, 135))
		label := p.Filtered[idx].Label
		rowText := label
		if p.Filtered[idx].Icon != "" {
			rowText = p.Filtered[idx].Icon + " " + label
		}
		term.DrawText(x, y+i, padRight(rowText, width), itemStyle)
		if p.Filtered[idx].Icon != "" {
			term.DrawText(x, y+i, p.Filtered[idx].Icon, iconStyle)
		}
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
	if runewidth.StringWidth(s) >= width {
		return s
	}
	return s + strings.Repeat(" ", width-runewidth.StringWidth(s))
}

func scoreItem(item Item, query string) (int, bool) {
	if query == "" {
		return 0, true
	}
	label := strings.ToLower(item.Label)
	base := strings.ToLower(filepath.Base(item.Label))
	switch {
	case base == query:
		return 0, true
	case strings.HasPrefix(base, query):
		return 1, true
	case strings.Contains(base, query):
		return 2, true
	case strings.HasPrefix(label, query):
		return 3, true
	case strings.Contains(label, query):
		return 4, true
	default:
		return 0, false
	}
}
