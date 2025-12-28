package clipboard

import (
	"golang.design/x/clipboard"
)

var initialized bool

// Init initializes the clipboard
func Init() error {
	if !initialized {
		err := clipboard.Init()
		if err != nil {
			return err
		}
		initialized = true
	}
	return nil
}

// Copy copies text to clipboard
func Copy(text string) error {
	if !initialized {
		if err := Init(); err != nil {
			return err
		}
	}
	clipboard.Write(clipboard.FmtText, []byte(text))
	return nil
}

// Paste gets text from clipboard
func Paste() (string, error) {
	if !initialized {
		if err := Init(); err != nil {
			return "", err
		}
	}
	data := clipboard.Read(clipboard.FmtText)
	return string(data), nil
}
