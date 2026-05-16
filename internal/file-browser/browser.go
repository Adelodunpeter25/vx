package filebrowser

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Adelodunpeter25/vx/internal/terminal"
	"github.com/Adelodunpeter25/vx/internal/utils"
	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
)

type Node struct {
	Name     string
	Path     string
	Icon     string
	IsDir    bool
	Expanded bool
	Loaded   bool
	Children []*Node
	Parent   *Node
}

type State struct {
	Open     bool
	Width    int
	Focused  bool
	RootPath string
	Root     *Node

	selected int
	scroll   int
	showHidden bool
}

type Action struct {
	OpenPath    string
	PreviewPath string
}

func New(root string) *State {
	if root == "" {
		if cwd, err := os.Getwd(); err == nil {
			root = cwd
		}
	}
	if root == "" {
		root = "."
	}
	root = filepath.Clean(root)
	state := &State{
		Open:     false,
		Width:    30,
		Focused:  false,
		RootPath: root,
	}
	state.Root = &Node{
		Name:  filepath.Base(root),
		Path:  root,
		Icon:  utils.FileIconFor(root, true),
		IsDir: true,
	}
	state.loadChildren(state.Root)
	state.Root.Expanded = true
	return state
}

func (s *State) SetRoot(root string) {
	if root == "" {
		if cwd, err := os.Getwd(); err == nil {
			root = cwd
		}
	}
	root = filepath.Clean(root)
	s.RootPath = root
	s.Root = &Node{
		Name:  filepath.Base(root),
		Path:  root,
		Icon:  utils.FileIconFor(root, true),
		IsDir: true,
	}
	s.selected = 0
	s.scroll = 0
	s.loadChildren(s.Root)
	s.Root.Expanded = true
}

func (s *State) SetSelectedPath(path string) {
	path = filepath.Clean(path)

	// Ensure all parents are expanded
	rel, err := filepath.Rel(s.RootPath, path)
	if err == nil && !strings.HasPrefix(rel, "..") {
		parts := strings.Split(rel, string(os.PathSeparator))
		curr := s.Root
		currPath := s.RootPath
		for i := 0; i < len(parts)-1; i++ {
			currPath = filepath.Join(currPath, parts[i])
			if !curr.Loaded {
				s.loadChildren(curr)
			}
			found := false
			for _, child := range curr.Children {
				if child.Path == currPath {
					child.Expanded = true
					curr = child
					found = true
					break
				}
			}
			if !found {
				break
			}
		}
	}

	nodes := s.Visible()
	for i, node := range nodes {
		if node.Path == path {
			s.selected = i
			return
		}
	}
}

func (s *State) Refresh() {
	if s.Root != nil {
		s.invalidateAll(s.Root)
	}
}

func (s *State) invalidateAll(node *Node) {
	if node == nil {
		return
	}
	node.Loaded = false
	for _, child := range node.Children {
		s.invalidateAll(child)
	}
}

func (s *State) Visible() []*Node {
	if s.Root == nil {
		return nil
	}
	var nodes []*Node
	for _, child := range s.Root.Children {
		s.appendVisible(&nodes, child, 0)
	}
	if s.selected >= len(nodes) {
		s.selected = len(nodes) - 1
	}
	if s.selected < 0 && len(nodes) > 0 {
		s.selected = 0
	}
	return nodes
}

func (s *State) appendVisible(out *[]*Node, node *Node, depth int) {
	if node == nil {
		return
	}
	*out = append(*out, node)
	if node.IsDir && node.Expanded {
		if !node.Loaded {
			s.loadChildren(node)
		}
		for _, child := range node.Children {
			s.appendVisible(out, child, depth+1)
		}
	}
}

func (s *State) loadChildren(node *Node) {
	if node == nil || !node.IsDir {
		return
	}
	entries, err := os.ReadDir(node.Path)
	if err != nil {
		node.Loaded = true
		return
	}
	children := make([]*Node, 0, len(entries))
	for _, ent := range entries {
		name := ent.Name()
		if strings.HasPrefix(name, ".") && !s.showHidden {
			continue
		}
		path := filepath.Join(node.Path, name)
		child := &Node{
			Name:   name,
			Path:   path,
			Icon:   utils.FileIconFor(path, ent.IsDir()),
			IsDir:  ent.IsDir(),
			Parent: node,
		}
		children = append(children, child)
	}
	sort.Slice(children, func(i, j int) bool {
		if children[i].IsDir != children[j].IsDir {
			return children[i].IsDir
		}
		return strings.ToLower(children[i].Name) < strings.ToLower(children[j].Name)
	})
	node.Children = children
	node.Loaded = true
}

func (s *State) SetShowHidden(show bool) {
	s.showHidden = show
	if s.Root != nil {
		s.Root.Loaded = false
		s.loadChildren(s.Root)
		s.Root.Expanded = true
	}
}

func (s *State) Render(term *terminal.Terminal, x, y, width, height int) {
	if !s.Open || term == nil || width <= 0 || height <= 0 {
		return
	}
	if s.Width > 0 {
		width = s.Width
	}

	bgStyle := tcell.StyleDefault
	for row := 0; row < height; row++ {
		for col := 0; col < width; col++ {
			term.SetCell(x+col, y+row, ' ', bgStyle)
		}
	}

	nodes := s.Visible()
	if len(nodes) == 0 {
		return
	}
	s.clampScroll(len(nodes), height)
	if s.selected < s.scroll {
		s.scroll = s.selected
	}
	if s.selected >= s.scroll+height {
		s.scroll = s.selected - height + 1
	}

	for row := 0; row < height; row++ {
		idx := s.scroll + row
		if idx >= len(nodes) {
			break
		}
		node := nodes[idx]
		depth := 0
		for p := node.Parent; p != nil && p != s.Root; p = p.Parent {
			depth++
		}
		iconStyle := tcell.StyleDefault.Foreground(tcell.ColorGray)
		icon := node.Icon
		if node.IsDir {
			iconStyle = tcell.StyleDefault.Foreground(tcell.NewRGBColor(250, 179, 135))
		} else {
			iconStyle = tcell.StyleDefault.Foreground(tcell.NewRGBColor(166, 227, 161))
		}
		indent := strings.Repeat("  ", depth)
		label := node.Name
		rowStyle := tcell.StyleDefault
		if idx == s.selected {
			rowStyle = rowStyle.Background(tcell.NewRGBColor(45, 46, 61)).Foreground(tcell.ColorWhite)
			if s.Focused {
				rowStyle = rowStyle.Bold(true)
			}
		}
		drawClippedText(term, x, y+row, width, strings.Repeat(" ", width), rowStyle)
		drawClippedText(term, x, y+row, width, indent, rowStyle)
		indentWidth := runewidth.StringWidth(indent)
		iconWidth := runewidth.StringWidth(icon + " ")
		if indentWidth < width {
			drawClippedText(term, x+indentWidth, y+row, width-indentWidth, icon+" ", iconStyle)
		}
		if runewidth.StringWidth(label) > width {
			label = runewidth.Truncate(label, width, "…")
		}
		nameX := x + indentWidth + iconWidth
		drawClippedText(term, nameX, y+row, width-indentWidth-iconWidth, label, rowStyle)
	}
}

func (s *State) HandleKey(ev *terminal.Event) Action {
	if s == nil || !s.Open || ev == nil {
		return Action{}
	}
	nodes := s.Visible()
	switch ev.Key {
	case tcell.KeyEscape:
		s.Focused = false
		return Action{}
	case tcell.KeyUp:
		if s.selected > 0 {
			s.selected--
		}
		return Action{}
	case tcell.KeyDown:
		if s.selected < len(nodes)-1 {
			s.selected++
		}
		return Action{}
	case tcell.KeyLeft:
		if s.selected >= 0 && s.selected < len(nodes) {
			node := nodes[s.selected]
			if node.IsDir && node.Expanded {
				node.Expanded = false
				return Action{}
			}
			if node.Parent != nil && node.Parent != s.Root {
				parent := node.Parent
				for i, n := range nodes {
					if n == parent {
						s.selected = i
						break
					}
				}
			}
		}
		return Action{}
	case tcell.KeyRight:
		if s.selected >= 0 && s.selected < len(nodes) {
			node := nodes[s.selected]
			if node.IsDir && !node.Expanded {
				node.Expanded = true
			}
		}
		return Action{}
	case tcell.KeyEnter:
		return s.activateSelection(nodes)
	}
	if ev.Rune != 0 && ev.Rune == ' ' {
		return s.activateSelection(nodes)
	}
	return Action{}
}

func (s *State) HandleMouse(ev *terminal.Event, x, y, width, height int) Action {
	if s == nil || !s.Open || ev == nil {
		return Action{}
	}
	if s.Width > 0 {
		width = s.Width
	}
	if ev.MouseX < x || ev.MouseX >= x+width || ev.MouseY < y || ev.MouseY >= y+height {
		return Action{}
	}
	if ev.Button != tcell.Button1 && ev.Button != tcell.WheelUp && ev.Button != tcell.WheelDown {
		return Action{}
	}
	if ev.Button == tcell.WheelUp {
		if s.scroll > 0 {
			s.scroll--
		}
		nodes := s.Visible()
		s.clampScroll(len(nodes), height)
		return Action{}
	}
	if ev.Button == tcell.WheelDown {
		s.scroll++
		nodes := s.Visible()
		s.clampScroll(len(nodes), height)
		return Action{}
	}

	nodes := s.Visible()
	row := ev.MouseY - y
	idx := s.scroll + row
	if idx < 0 || idx >= len(nodes) {
		return Action{}
	}
	s.selected = idx
	return s.previewSelection(nodes)
}

func (s *State) activateSelection(nodes []*Node) Action {
	if s.selected < 0 || s.selected >= len(nodes) {
		return Action{}
	}
	node := nodes[s.selected]
	if node.IsDir {
		node.Expanded = !node.Expanded
		return Action{}
	}
	return Action{OpenPath: node.Path}
}

func (s *State) previewSelection(nodes []*Node) Action {
	if s.selected < 0 || s.selected >= len(nodes) {
		return Action{}
	}
	node := nodes[s.selected]
	if node.IsDir {
		node.Expanded = !node.Expanded
		return Action{}
	}
	return Action{PreviewPath: node.Path}
}

func (s *State) clampScroll(nodesLen int, height int) {
	if height < 1 {
		height = 1
	}
	maxScroll := nodesLen - height
	if maxScroll < 0 {
		maxScroll = 0
	}
	if s.scroll < 0 {
		s.scroll = 0
	}
	if s.scroll > maxScroll {
		s.scroll = maxScroll
	}
}

func padRight(s string, width int) string {
	sw := runewidth.StringWidth(s)
	if sw >= width {
		return s
	}
	return s + strings.Repeat(" ", width-sw)
}

func drawClippedText(term *terminal.Terminal, x, y, width int, text string, style tcell.Style) {
	if term == nil || width <= 0 || text == "" {
		return
	}

	cursorX := x
	remaining := width
	for _, r := range text {
		rw := runewidth.RuneWidth(r)
		if rw <= 0 {
			rw = 1
		}
		if rw > remaining {
			break
		}
		term.SetCell(cursorX, y, r, style)
		cursorX += rw
		remaining -= rw
	}
}
