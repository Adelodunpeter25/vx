package clipboard

import (
	"github.com/atotto/clipboard"
)

// Copy copies text to clipboard
func Copy(text string) error {
	return clipboard.WriteAll(text)
}

// Paste gets text from clipboard
func Paste() (string, error) {
	return clipboard.ReadAll()
}
