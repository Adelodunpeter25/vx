package main

import (
	"fmt"
	"os"

	"github.com/Adelodunpeter25/vx/internal/editor"
	"github.com/Adelodunpeter25/vx/internal/terminal"
)

func main() {
	term, err := terminal.New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize terminal: %v\n", err)
		os.Exit(1)
	}
	defer term.Close()

	ed := editor.New(term)
	if err := ed.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
