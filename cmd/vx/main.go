package main

import (
	"fmt"
	"os"

	"github.com/Adelodunpeter25/vx/internal/editor"
	"github.com/Adelodunpeter25/vx/internal/terminal"
)

func main() {
	// Handle flags
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "h","-h", "--help":
			printHelp()
			return
		case "-v", "--version":
			printVersion()
			return
		}
		
		// Check if path is a directory before initializing terminal
		filename := os.Args[1]
		if info, err := os.Stat(filename); err == nil && info.IsDir() {
			fmt.Fprintf(os.Stderr, "Error: '%s' is a directory\n", filename)
			fmt.Fprintf(os.Stderr, "Usage: vx <filename>\n")
			os.Exit(1)
		}
	}
	
	term, err := terminal.New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize terminal: %v\n", err)
		os.Exit(1)
	}
	defer term.Close()

	var ed *editor.Editor
	if len(os.Args) > 1 {
		ed, err = editor.NewWithFile(term, os.Args[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to load file: %v\n", err)
			os.Exit(1)
		}
	} else {
		ed = editor.New(term)
	}

	if err := ed.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}