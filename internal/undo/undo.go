package undo

// Action represents a single undoable action
type Action struct {
	Type     ActionType
	Line     int
	Col      int
	Text     string
	OldText  string
}

type ActionType int

const (
	ActionInsertRune ActionType = iota
	ActionDeleteRune
	ActionInsertLine
	ActionDeleteLine
	ActionSplitLine
	ActionJoinLine
)

// Stack manages undo/redo history
type Stack struct {
	actions []Action
	current int
}

func NewStack() *Stack {
	return &Stack{
		actions: make([]Action, 0),
		current: -1,
	}
}

// Push adds a new action and clears redo history
func (s *Stack) Push(action Action) {
	// Remove any actions after current position (clear redo history)
	if s.current < len(s.actions)-1 {
		s.actions = s.actions[:s.current+1]
	}
	
	s.actions = append(s.actions, action)
	s.current++
}

// Undo returns the action to undo, or nil if nothing to undo
func (s *Stack) Undo() *Action {
	if s.current < 0 {
		return nil
	}
	
	action := s.actions[s.current]
	s.current--
	return &action
}

// Redo returns the action to redo, or nil if nothing to redo
func (s *Stack) Redo() *Action {
	if s.current >= len(s.actions)-1 {
		return nil
	}
	
	s.current++
	return &s.actions[s.current]
}

// CanUndo returns true if there are actions to undo
func (s *Stack) CanUndo() bool {
	return s.current >= 0
}

// CanRedo returns true if there are actions to redo
func (s *Stack) CanRedo() bool {
	return s.current < len(s.actions)-1
}

// Clear clears all history
func (s *Stack) Clear() {
	s.actions = make([]Action, 0)
	s.current = -1
}
