package terminalpane

import (
	"context"
	"os"
	"sync"
)

type Pane struct {
	cfg       Config
	pty       *Pty
	emu       *Emulator
	renderer  *Renderer
	scrollback *Scrollback
	limiter   *FrameLimiter

	mu     sync.Mutex
	width  int
	height int
	dirty  bool
	notify func()
	cancel context.CancelFunc
}

func NewPane(cfg Config) *Pane {
	if cfg.Shell == "" {
		cfg.Shell = os.Getenv("SHELL")
		if cfg.Shell == "" {
			cfg.Shell = "sh"
		}
	}
	if cfg.Scrollback <= 0 {
		cfg.Scrollback = 2000
	}
	return &Pane{
		cfg:        cfg,
		emu:        NewEmulator(80, 24),
		renderer:   NewRenderer(80, 24),
		scrollback: NewScrollback(cfg.Scrollback),
		limiter:    NewFrameLimiter(cfg.MaxFPS, cfg.MinRedraw),
		pty:        &Pty{},
	}
}

func (p *Pane) Start(cols, rows int) error {
	p.mu.Lock()
	p.width = cols
	p.height = rows
	p.mu.Unlock()
	if err := p.pty.Start(p.cfg.Shell, p.cfg.Env, p.cfg.Dir); err != nil {
		return err
	}
	p.pty.Resize(cols, rows)
	p.emu.Resize(cols, rows)
	p.renderer.Resize(cols, rows)
	ctx, cancel := context.WithCancel(context.Background())
	p.cancel = cancel
	go p.pty.ReadLoop(ctx, p.onPTYData)
	return nil
}

func (p *Pane) Stop() {
	if p.cancel != nil {
		p.cancel()
	}
	p.pty.Stop()
}

func (p *Pane) Resize(cols, rows int) {
	p.mu.Lock()
	p.width = cols
	p.height = rows
	p.mu.Unlock()
	p.pty.Resize(cols, rows)
	p.emu.Resize(cols, rows)
	p.renderer.Resize(cols, rows)
	p.markDirty()
}

func (p *Pane) Write(data []byte) {
	p.pty.Write(data)
}

func (p *Pane) Emulator() *Emulator {
	return p.emu
}

func (p *Pane) Renderer() *Renderer {
	return p.renderer
}

func (p *Pane) Scrollback() *Scrollback {
	return p.scrollback
}

func (p *Pane) SetNotify(fn func()) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.notify = fn
}

func (p *Pane) Dirty() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.dirty
}

func (p *Pane) ClearDirty() {
	p.mu.Lock()
	p.dirty = false
	p.mu.Unlock()
}

func (p *Pane) onPTYData(data []byte) {
	p.emu.Write(data)
	p.markDirty()
}

func (p *Pane) markDirty() {
	p.mu.Lock()
	p.dirty = true
	notify := p.notify
	p.mu.Unlock()
	if notify != nil && p.limiter.Allow() {
		notify()
	}
}
