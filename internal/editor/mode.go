package editor

type Mode int

const (
	ModeNormal Mode = iota
	ModeInsert
	ModeCommand
	ModeSearch
	ModeReplace
	ModeBufferPrompt
	ModeCdPrompt
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
	case ModeCdPrompt:
		return "CD"
	default:
		return "UNKNOWN"
	}
}
