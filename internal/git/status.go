package git

type Status int

const (
	StatusClean Status = iota
	StatusModified
	StatusAdded
	StatusDeleted
	StatusUntracked
	StatusIgnored
	StatusConflicted
	StatusUnknown
)

func (s Status) String() string {
	switch s {
	case StatusClean:
		return "clean"
	case StatusModified:
		return "modified"
	case StatusAdded:
		return "added"
	case StatusDeleted:
		return "deleted"
	case StatusUntracked:
		return "untracked"
	case StatusIgnored:
		return "ignored"
	case StatusConflicted:
		return "conflicted"
	default:
		return "unknown"
	}
}

func (s Status) Badge() string {
	switch s {
	case StatusModified:
		return "M"
	case StatusAdded:
		return "A"
	case StatusDeleted:
		return "D"
	case StatusUntracked:
		return "?"
	case StatusIgnored:
		return "!"
	case StatusConflicted:
		return "U"
	case StatusClean:
		return " "
	default:
		return "?"
	}
}

