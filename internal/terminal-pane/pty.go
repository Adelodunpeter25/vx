package terminalpane

import (
	"context"
	"os"
	"os/exec"
	"sync"

	"github.com/creack/pty"
)

type Pty struct {
	cmd  *exec.Cmd
	file *os.File
	mu   sync.Mutex
}

func (p *Pty) Start(shell string, env []string, dir string) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.cmd != nil {
		return nil
	}
	cmd := exec.Command(shell)
	if dir != "" {
		cmd.Dir = dir
	}
	if len(env) > 0 {
		cmd.Env = append(os.Environ(), env...)
	}
	ptmx, err := pty.Start(cmd)
	if err != nil {
		return err
	}
	p.cmd = cmd
	p.file = ptmx
	return nil
}

func (p *Pty) Stop() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.cmd == nil {
		return
	}
	_ = p.file.Close()
	_ = p.cmd.Process.Kill()
	_, _ = p.cmd.Process.Wait()
	p.cmd = nil
	p.file = nil
}

func (p *Pty) Resize(cols, rows int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.file == nil {
		return
	}
	_ = pty.Setsize(p.file, &pty.Winsize{Cols: uint16(cols), Rows: uint16(rows)})
}

func (p *Pty) Write(data []byte) {
	p.mu.Lock()
	file := p.file
	p.mu.Unlock()
	if file == nil || len(data) == 0 {
		return
	}
	_, _ = file.Write(data)
}

func (p *Pty) ReadLoop(ctx context.Context, onData func([]byte)) {
	buf := make([]byte, 4096)
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		p.mu.Lock()
		file := p.file
		p.mu.Unlock()
		if file == nil {
			return
		}
		n, err := file.Read(buf)
		if n > 0 {
			chunk := make([]byte, n)
			copy(chunk, buf[:n])
			onData(chunk)
		}
		if err != nil {
			return
		}
	}
}
