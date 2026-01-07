package replace

// State represents the current state of replace operation
type State int

const (
	StateSearchInput State = iota
	StateReplaceInput
	StateConfirm
	StateInactive
)
