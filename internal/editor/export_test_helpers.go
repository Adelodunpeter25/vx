package editor

// Exported helpers for tests in the tests/ package.

func NewRenderCache() *RenderCache { return newRenderCache() }

func (rc *RenderCache) NeedsRedraw() bool                    { return rc.needsRedraw }
func (rc *RenderCache) SetNeedsRedraw(v bool)                 { rc.needsRedraw = v }
func (rc *RenderCache) MarkNeedsRedraw()                      { rc.invalidate() }
func (rc *RenderCache) LineChanged(n int, s string) bool      { return rc.lineChanged(n, s) }
func (rc *RenderCache) UpdateLine(n int, s string)            { rc.updateLine(n, s) }
func (rc *RenderCache) InvalidateLine(n int)                  { rc.invalidateLine(n) }
