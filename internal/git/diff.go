package git

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

type LineChange int

const (
	LineChangeNone LineChange = iota
	LineChangeAdded
	LineChangeModified
)

func DiffLineChanges(ctx context.Context, filePath string) (map[int]LineChange, error) {
	if filePath == "" {
		return map[int]LineChange{}, nil
	}
	absFile, err := filepath.Abs(filePath)
	if err != nil {
		return nil, err
	}
	dir := filepath.Dir(absFile)

	repoRoot, err := repoRoot(ctx, dir)
	if err != nil {
		return map[int]LineChange{}, err
	}

	rel, err := filepath.Rel(repoRoot, absFile)
	if err != nil {
		return map[int]LineChange{}, err
	}
	rel = filepath.Clean(rel)

	if isUntracked, err := untracked(ctx, repoRoot, rel); err == nil && isUntracked {
		return lineChangesFromFile(absFile), nil
	}

	cmd := exec.CommandContext(ctx, "git", "-C", repoRoot, "diff", "--unified=0", "--no-ext-diff", "HEAD", "--", rel)
	out, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && len(exitErr.Stderr) > 0 {
			return nil, fmt.Errorf("git diff: %s", strings.TrimSpace(string(exitErr.Stderr)))
		}
		return map[int]LineChange{}, nil
	}

	return parseUnifiedDiff(out), nil
}

func repoRoot(ctx context.Context, dir string) (string, error) {
	cmd := exec.CommandContext(ctx, "git", "-C", dir, "rev-parse", "--show-toplevel")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	root := strings.TrimSpace(string(out))
	if root == "" {
		return "", fmt.Errorf("git repo root not found")
	}
	return filepath.Clean(root), nil
}

func untracked(ctx context.Context, repoRoot, rel string) (bool, error) {
	cmd := exec.CommandContext(ctx, "git", "-C", repoRoot, "status", "--porcelain=v1", "--", rel)
	out, err := cmd.Output()
	if err != nil {
		return false, err
	}
	return strings.HasPrefix(string(out), "??"), nil
}

func lineChangesFromFile(path string) map[int]LineChange {
	data, err := os.ReadFile(path)
	if err != nil {
		return map[int]LineChange{}
	}
	lines := bytes.Count(data, []byte{'\n'}) + 1
	out := make(map[int]LineChange, lines)
	for i := 1; i <= lines; i++ {
		out[i] = LineChangeAdded
	}
	return out
}

func parseUnifiedDiff(out []byte) map[int]LineChange {
	changes := make(map[int]LineChange)
	lines := bytes.Split(out, []byte{'\n'})
	for _, raw := range lines {
		line := string(raw)
		if !strings.HasPrefix(line, "@@ ") {
			continue
		}
		// Format: @@ -oldStart,oldCount +newStart,newCount @@
		parts := strings.Split(line, " ")
		if len(parts) < 3 {
			continue
		}
		newSpec := strings.TrimPrefix(parts[2], "+")
		newStart, newCount := parseHunkSpec(newSpec)
		if newStart == 0 || newCount == 0 {
			continue
		}
		for i := 0; i < newCount; i++ {
			changes[newStart+i] = LineChangeModified
		}
	}
	return changes
}

func parseHunkSpec(spec string) (start, count int) {
	if spec == "" {
		return 0, 0
	}
	spec = strings.TrimSpace(strings.TrimSuffix(spec, "@@"))
	spec = strings.TrimPrefix(spec, "+")
	parts := strings.SplitN(spec, ",", 2)
	start, _ = strconv.Atoi(strings.TrimSpace(parts[0]))
	count = 1
	if len(parts) == 2 {
		if n, err := strconv.Atoi(strings.TrimSpace(parts[1])); err == nil && n > 0 {
			count = n
		}
	}
	return start, count
}

