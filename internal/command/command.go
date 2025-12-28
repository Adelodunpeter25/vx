package command

import (
	"fmt"
	"strings"

	"github.com/Adelodunpeter25/vx/internal/buffer"
)

type Result struct {
	Quit    bool
	Message string
	Error   error
}

func Execute(cmd string, buf *buffer.Buffer) Result {
	cmd = strings.TrimSpace(cmd)
	
	switch cmd {
	case "q":
		if buf.IsModified() {
			return Result{Error: fmt.Errorf("no write since last change (use :q! to override)")}
		}
		return Result{Quit: true}
	
	case "q!":
		return Result{Quit: true}
	
	case "w":
		if buf.Filename() == "" {
			return Result{Error: fmt.Errorf("no file name")}
		}
		if err := buf.Save(); err != nil {
			return Result{Error: err}
		}
		return Result{Message: fmt.Sprintf("\"%s\" written", buf.Filename())}
	
	case "wq":
		if buf.Filename() == "" {
			return Result{Error: fmt.Errorf("no file name")}
		}
		if err := buf.Save(); err != nil {
			return Result{Error: err}
		}
		return Result{Quit: true, Message: fmt.Sprintf("\"%s\" written", buf.Filename())}
	
	default:
		if strings.HasPrefix(cmd, "w ") {
			filename := strings.TrimSpace(cmd[2:])
			if filename == "" {
				return Result{Error: fmt.Errorf("no file name")}
			}
			// TODO: set filename in buffer
			if err := buf.Save(); err != nil {
				return Result{Error: err}
			}
			return Result{Message: fmt.Sprintf("\"%s\" written", filename)}
		}
		return Result{Error: fmt.Errorf("not an editor command: :%s", cmd)}
	}
}
