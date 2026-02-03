package editor

// RenderCache tracks what was rendered to avoid unnecessary redraws
type RenderCache struct {
	lines       map[int]string
	cursorX     int
	cursorY     int
	offsetY     int
	width       int
	height      int
	statusLine  string
	needsRedraw bool
}

func newRenderCache() *RenderCache {
	return &RenderCache{
		lines:       make(map[int]string),
		needsRedraw: true,
	}
}

func (rc *RenderCache) hasChanged(p *Pane, width, height int) bool {
	if rc.needsRedraw {
		return true
	}

	if rc.cursorX != p.cursorX || rc.cursorY != p.cursorY {
		return true
	}

	if rc.offsetY != p.offsetY {
		return true
	}

	if rc.width != width || rc.height != height {
		return true
	}

	return false
}

func (rc *RenderCache) lineChanged(lineNum int, content string) bool {
	cached, exists := rc.lines[lineNum]
	return !exists || cached != content
}

func (rc *RenderCache) updateLine(lineNum int, content string) {
	rc.lines[lineNum] = content
}

func (rc *RenderCache) update(p *Pane, width, height int) {
	rc.cursorX = p.cursorX
	rc.cursorY = p.cursorY
	rc.offsetY = p.offsetY
	rc.width = width
	rc.height = height
	rc.needsRedraw = false
}

func (rc *RenderCache) invalidate() {
	rc.needsRedraw = true
	rc.lines = make(map[int]string)
}

func (rc *RenderCache) invalidateLine(lineNum int) {
	delete(rc.lines, lineNum)
}
