package git

import "path/filepath"

func Aggregate(statuses map[string]Status) map[string]Status {
	out := make(map[string]Status)
	for path, status := range statuses {
		out[path] = status
		dir := filepath.Dir(path)
		for dir != "." && dir != string(filepath.Separator) {
			out[dir] = merge(out[dir], status)
			next := filepath.Dir(dir)
			if next == dir {
				break
			}
			dir = next
		}
	}
	return out
}

func merge(current, next Status) Status {
	if current == StatusConflicted || next == StatusConflicted {
		return StatusConflicted
	}
	if current == StatusModified || next == StatusModified {
		return StatusModified
	}
	if current == StatusAdded || next == StatusAdded {
		return StatusAdded
	}
	if current == StatusDeleted || next == StatusDeleted {
		return StatusDeleted
	}
	if current == StatusUntracked || next == StatusUntracked {
		return StatusUntracked
	}
	if current == StatusIgnored || next == StatusIgnored {
		return StatusIgnored
	}
	if current == StatusUnknown {
		return next
	}
	return current
}

