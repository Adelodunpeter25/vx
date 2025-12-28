package main

import (
	"fmt"
	"os"

	"github.com/Adelodunpeter25/vx/internal/terminal"
	"github.com/gdamore/tcell/v2"
)

func main() {
	term, err := terminal.New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize terminal: %v\n", err)
		os.Exit(1)
	}
	defer term.Close()

	// Test: draw welcome message
	width, height := term.Size()
	msg := "vx - press 'q' to quit"
	x := (width - len(msg)) / 2
	y := height / 2

	term.Clear()
	term.DrawText(x, y, msg, tcell.StyleDefault)
	term.Show()

	// Event loop
	for {
		ev := term.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			if ev.Rune() == 'q' {
				return
			}
		case *tcell.EventResize:
			term.Clear()
			width, height = term.Size()
			x = (width - len(msg)) / 2
			y = height / 2
			term.DrawText(x, y, msg, tcell.StyleDefault)
			term.Show()
		}
	}
}
