package editor

import (
	"context"
	"path/filepath"

	"github.com/Adelodunpeter25/vx/internal/git"
)

func (e *Editor) ensureGitLines(p *Pane) {
	if p == nil || p.buffer == nil {
		return
	}
	filename := p.buffer.Filename()
	if filename == "" {
		p.gitLines = nil
		p.gitLineFile = ""
		p.gitLineVersion = p.buffer.ModVersion()
		return
	}
	if p.gitLineFile == filename && p.gitLineVersion == p.buffer.ModVersion() {
		return
	}

	lines, err := git.DiffLineChanges(context.Background(), filename)
	if err != nil {
		p.gitLines = nil
		p.gitLineFile = filename
		p.gitLineVersion = p.buffer.ModVersion()
		return
	}
	p.gitLines = lines
	p.gitLineFile = filepath.Clean(filename)
	p.gitLineVersion = p.buffer.ModVersion()
}

func (p *Pane) gitLineChange(line int) git.LineChange {
	if p == nil || p.gitLines == nil {
		return git.LineChangeNone
	}
	if change, ok := p.gitLines[line]; ok {
		return change
	}
	return git.LineChangeNone
}

