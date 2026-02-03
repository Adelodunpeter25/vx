package splitpane

// Rect defines a pane rectangle in screen coordinates.
type Rect struct {
	X      int
	Y      int
	Width  int
	Height int
}

// LayoutSideBySide returns pane rectangles split horizontally with a divider.
// splitRatio is between 0.1 and 0.9 and represents left pane width percentage.
// DividerX is the column where the divider should be drawn. -1 if none.
func LayoutSideBySide(width, height, panes int, splitRatio float64) ([]Rect, int) {
	if panes <= 1 {
		return []Rect{{X: 0, Y: 0, Width: width, Height: height}}, -1
	}
	if splitRatio < 0.1 {
		splitRatio = 0.1
	}
	if splitRatio > 0.9 {
		splitRatio = 0.9
	}

	divider := 1
	available := width - divider
	if available < 2 {
		return []Rect{{X: 0, Y: 0, Width: width, Height: height}}, -1
	}

	leftWidth := int(float64(available) * splitRatio)
	if leftWidth < 1 {
		leftWidth = 1
	}
	rightWidth := available - leftWidth
	if rightWidth < 1 {
		rightWidth = 1
		leftWidth = available - rightWidth
	}

	left := Rect{X: 0, Y: 0, Width: leftWidth, Height: height}
	dividerX := leftWidth
	right := Rect{X: leftWidth + divider, Y: 0, Width: rightWidth, Height: height}

	return []Rect{left, right}, dividerX
}
