package editor

import (
	"github.com/Adelodunpeter25/vx/internal/buffer"
	"github.com/Adelodunpeter25/vx/internal/buffers"
	"github.com/Adelodunpeter25/vx/internal/preview"
	"github.com/Adelodunpeter25/vx/internal/replace"
	"github.com/Adelodunpeter25/vx/internal/search"
	"github.com/Adelodunpeter25/vx/internal/syntax"
	"github.com/Adelodunpeter25/vx/internal/visual"
)

type Pane struct {
	bufferMgr     *buffers.Manager
	buffer        *buffer.Buffer
	syntax        *syntax.Engine
	search        *search.Engine
	replace       *replace.Engine
	preview       *preview.Preview
	selection     *visual.Selection
	renderCache   *RenderCache
	msgManager    *MessageManager
	cursorX       int
	cursorY       int
	offsetX       int
	offsetY       int
	visualOffsetY int
	mode          Mode
	commandBuf    string
	searchBuf     string
	lastKey       rune
	mouseDownX    int
	mouseDownY    int
	mouseDragging bool
	viewX         int
	viewY         int
	viewWidth     int
	viewHeight    int
}

func NewPane(buf *buffer.Buffer, filename string) *Pane {
	bufMgr := buffers.New(buf, filename)
	return &Pane{
		bufferMgr:   bufMgr,
		buffer:      buf,
		syntax:      bufMgr.Current().Syntax,
		search:      search.New(),
		replace:     replace.New(),
		preview:     preview.New(),
		selection:   visual.New(),
		renderCache: newRenderCache(),
		msgManager:  NewMessageManager(),
		mode:        ModeNormal,
	}
}
