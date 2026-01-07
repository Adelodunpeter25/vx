package editor

type Mode int

const (
	ModeNormal Mode = iota
	ModeInsert
	ModeCommand
	ModeSearch
	ModeReplace
	ModeBufferPrompt
)

func (m Mode) String() string {
	switch m {
	case ModeNormal:
		return "NORMAL"
	case ModeInsert:
		return "INSERT"
	case ModeCommand:
		return "COMMAND"
	case ModeSearch:
		return "SEARCH"
	case ModeReplace:
		return "REPLACE"
	case ModeBufferPrompt:
		return "PROMPT"
	default:
		return "UNKNOWN"
	}
}
