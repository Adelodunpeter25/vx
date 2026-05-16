package git

import (
	"bytes"
	"context"
	"errors"
	"os/exec"
	"path/filepath"
	"strings"
)

type Collector struct {
	root string
}

type Snapshot struct {
	Root        string
	RepoRoot    string
	Statuses    map[string]Status
	DirStatuses map[string]Status
}

func NewCollector(root string) *Collector {
	return &Collector{root: filepath.Clean(root)}
}

func (c *Collector) Collect(ctx context.Context) (*Snapshot, error) {
	repoRoot, err := c.repoRoot(ctx)
	if err != nil {
		return nil, err
	}

	statuses, err := c.collectStatuses(ctx, repoRoot)
	if err != nil {
		return nil, err
	}

	snap := &Snapshot{
		Root:        c.root,
		RepoRoot:    repoRoot,
		Statuses:    statuses,
		DirStatuses: make(map[string]Status),
	}
	snap.DirStatuses = Aggregate(statuses)
	return snap, nil
}

func (c *Collector) repoRoot(ctx context.Context) (string, error) {
	cmd := exec.CommandContext(ctx, "git", "-C", c.root, "rev-parse", "--show-toplevel")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	root := strings.TrimSpace(string(out))
	if root == "" {
		return "", errors.New("git repo root not found")
	}
	return filepath.Clean(root), nil
}

func (c *Collector) collectStatuses(ctx context.Context, repoRoot string) (map[string]Status, error) {
	cmd := exec.CommandContext(ctx, "git", "-C", repoRoot, "status", "--porcelain=v1", "-z")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	statuses := make(map[string]Status)
	entries := bytes.Split(out, []byte{0})
	for _, raw := range entries {
		if len(raw) == 0 {
			continue
		}
		line := string(raw)
		if len(line) < 3 {
			continue
		}
		code := line[:2]
		path := strings.TrimSpace(line[3:])
		if path == "" {
			continue
		}
		statuses[filepath.Clean(path)] = parsePorcelainStatus(code)
	}
	return statuses, nil
}

func parsePorcelainStatus(code string) Status {
	switch {
	case strings.Contains(code, "U"):
		return StatusConflicted
	case strings.Contains(code, "?"):
		return StatusUntracked
	case strings.Contains(code, "!"):
		return StatusIgnored
	case strings.Contains(code, "D"):
		return StatusDeleted
	case strings.Contains(code, "A"):
		return StatusAdded
	case strings.Contains(code, "M"):
		return StatusModified
	default:
		return StatusClean
	}
}

