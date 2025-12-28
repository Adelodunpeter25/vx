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

func (rc *RenderCache) hasChanged(e *Editor) bool {
	if rc.needsRedraw {
		return true
	}
	
	if rc.cursorX != e.cursorX || rc.cursorY != e.cursorY {
		return true
	}
	
	if rc.offsetY != e.offsetY {
		return true
	}
	
	if rc.width != e.width || rc.height != e.height {
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

func (rc *RenderCache) update(e *Editor) {
	rc.cursorX = e.cursorX
	rc.cursorY = e.cursorY
	rc.offsetY = e.offsetY
	rc.width = e.width
	rc.height = e.height
	rc.needsRedraw = false
}

func (rc *RenderCache) invalidate() {
	rc.needsRedraw = true
	rc.lines = make(map[int]string)
}

func (rc *RenderCache) invalidateLine(lineNum int) {
	delete(rc.lines, lineNum)
}
