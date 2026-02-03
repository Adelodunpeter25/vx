package splitpane

// Rect defines a pane rectangle in screen coordinates.
type Rect struct {
	X      int
	Y      int
	Width  int
	Height int
}

// LayoutSideBySide returns pane rectangles split horizontally with a 1-column divider.
// DividerX is the column where the divider should be drawn (only for 2 panes). -1 if none.
func LayoutSideBySide(width, height, panes int) ([]Rect, int) {
	if panes <= 1 {
		return []Rect{{X: 0, Y: 0, Width: width, Height: height}}, -1
	}

	divider := 1
	available := width - divider
	if available < panes {
		available = panes
	}

	base := available / panes
	extra := available % panes

	rects := make([]Rect, panes)
	x := 0
	dividerX := -1
	for i := 0; i < panes; i++ {
		w := base
		if i < extra {
			w++
		}
		rects[i] = Rect{X: x, Y: 0, Width: w, Height: height}
		x += w
		if i == 0 {
			dividerX = x
			x += divider
		}
	}

	return rects, dividerX
}
