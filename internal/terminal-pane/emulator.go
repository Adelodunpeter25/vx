package terminalpane

import (
	"sync"

	uv "github.com/charmbracelet/ultraviolet"
	"github.com/charmbracelet/x/vt"
)

type Emulator struct {
	mu  sync.RWMutex
	emu *vt.SafeEmulator
}

func NewEmulator(cols, rows int) *Emulator {
	return &Emulator{emu: vt.NewSafeEmulator(cols, rows)}
}

func (e *Emulator) Resize(cols, rows int) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.emu == nil {
		return
	}
	e.emu.Resize(cols, rows)
}

func (e *Emulator) Write(data []byte) {
	e.mu.RLock()
	emu := e.emu
	e.mu.RUnlock()
	if emu == nil {
		return
	}
	_, _ = emu.Write(data)
}

func (e *Emulator) CellAt(x, y int) *uv.Cell {
	e.mu.RLock()
	defer e.mu.RUnlock()
	if e.emu == nil {
		return nil
	}
	return e.emu.CellAt(x, y)
}

func (e *Emulator) Cursor() vt.CursorPosition {
	e.mu.RLock()
	defer e.mu.RUnlock()
	if e.emu == nil {
		return vt.CursorPosition{}
	}
	return e.emu.CursorPosition()
}

func (e *Emulator) Size() (int, int) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	if e.emu == nil {
		return 0, 0
	}
	return e.emu.Size()
}
