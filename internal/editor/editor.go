package editor

import (
	"github.com/Adelodunpeter25/vx/internal/buffer"
	"github.com/Adelodunpeter25/vx/internal/buffers"
	"github.com/Adelodunpeter25/vx/internal/preview"
	"github.com/Adelodunpeter25/vx/internal/replace"
	"github.com/Adelodunpeter25/vx/internal/search"
	"github.com/Adelodunpeter25/vx/internal/syntax"
	"github.com/Adelodunpeter25/vx/internal/terminal"
	"github.com/Adelodunpeter25/vx/internal/utils"
	"github.com/Adelodunpeter25/vx/internal/visual"
	"github.com/gdamore/tcell/v2"
)

type Editor struct {
	term        *terminal.Terminal
	bufferMgr   *buffers.Manager
	buffer      *buffer.Buffer
	syntax      *syntax.Engine
	search      *search.Engine
	replace     *replace.Engine
	preview     *preview.Preview
	selection   *visual.Selection
	renderCache *RenderCache
	width       int
	height      int
	cursorX     int
	cursorY     int
	offsetX     int // Horizontal scroll offset
	offsetY     int
	mode        Mode
	quit        bool
	commandBuf  string
	searchBuf   string
	message     string
	lastKey     rune // Track last key for multi-key commands like gg
	mouseDownX  int  // Track mouse button down position
	mouseDownY  int
	mouseDragging bool
}

func New(term *terminal.Terminal) *Editor {
	width, height := term.Size()
	buf := buffer.New()
	bufMgr := buffers.New(buf, "")
	
	return &Editor{
		term:        term,
		bufferMgr:   bufMgr,
		buffer:      buf,
		syntax:      bufMgr.Current().Syntax,
		search:      search.New(),
		replace:     replace.New(),
		preview:     preview.New(),
		selection:   visual.New(),
		renderCache: newRenderCache(),
		width:       width,
		height:      height,
		mode:        ModeNormal,
	}
}

func NewWithFile(term *terminal.Terminal, filename string) (*Editor, error) {
	buf, err := buffer.Load(filename)
	if err != nil {
		// Check if it's a partial load (recoverable error)
		if buf != nil {
			// File loaded with warnings
			width, height := term.Size()
			bufMgr := buffers.New(buf, filename)
			ed := &Editor{
				term:        term,
				bufferMgr:   bufMgr,
				buffer:      buf,
				syntax:      bufMgr.Current().Syntax,
				search:      search.New(),
				replace:     replace.New(),
				selection:   visual.New(),
				renderCache: newRenderCache(),
				width:       width,
				height:      height,
				mode:        ModeNormal,
				message:     "Warning: " + utils.FormatLoadError(filename, err),
			}
			return ed, nil
		}
		return nil, err
	}
	
	width, height := term.Size()
	bufMgr := buffers.New(buf, filename)
	ed := &Editor{
		term:        term,
		bufferMgr:   bufMgr,
		buffer:      buf,
		syntax:      bufMgr.Current().Syntax,
		search:      search.New(),
		replace:     replace.New(),
		preview:     preview.New(),
		selection:   visual.New(),
		renderCache: newRenderCache(),
		width:       width,
		height:      height,
		mode:        ModeNormal,
	}
	
	// Show file info message on load
	ed.showFileInfo()
	
	// Show warning if syntax highlighting disabled due to size
	if ed.syntax.IsTooLarge() {
		ed.message = "File too large for syntax highlighting"
	}
	
	return ed, nil
}

func (e *Editor) showFileInfo() {
	size, err := e.buffer.GetFileSize()
	if err != nil {
		return
	}
	
	filename := e.buffer.Filename()
	if filename == "" {
		filename = "[No Name]"
	}
	
	e.message = utils.FormatFileInfo(filename, size, e.buffer.LineCount())
}

func (e *Editor) Run() error {
	e.render()
	
	for !e.quit {
		e.handleEvent()
	}
	return nil
}

func (e *Editor) handleEvent() {
	ev := e.term.ReadEvent()
	if ev == nil {
		return
	}

	switch ev.Type {
	case terminal.EventKey:
		e.handleKey(ev)
	case terminal.EventResize:
		e.handleResize()
	case terminal.EventMouse:
		e.handleMouseEvent(ev)
		e.renderCache.invalidate()
		e.render()
	}
}

func (e *Editor) handleResize() {
	e.width, e.height = e.term.Size()
	
	// Ensure cursor stays visible after resize
	e.clampCursor()
	e.adjustScroll()
	
	e.renderCache.invalidate()
	e.render()
}

func (e *Editor) handleKey(ev *terminal.Event) {
	switch e.mode {
	case ModeNormal:
		e.handleNormalMode(ev)
	case ModeInsert:
		e.handleInsertMode(ev)
	case ModeCommand:
		e.handleCommandMode(ev)
	case ModeSearch:
		e.handleSearchMode(ev)
	case ModeReplace:
		// Convert terminal.Event to tcell.EventKey for replace mode
		tcellEv := tcell.NewEventKey(ev.Key, ev.Rune, tcell.ModNone)
		e.handleReplaceMode(tcellEv)
	case ModeBufferPrompt:
		tcellEv := tcell.NewEventKey(ev.Key, ev.Rune, tcell.ModNone)
		e.handleBufferPromptMode(tcellEv)
	}
	e.renderCache.invalidate()
	e.render()
}
