package wrap

// Line represents a wrapped line segment
type Line struct {
	Text       string
	StartCol   int // Column in original line where this segment starts
	LineNum    int // Original line number
	IsWrapped  bool // True if this is a continuation (not first segment)
}
