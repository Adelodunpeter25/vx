package command

import (
	"fmt"
	"strings"

	"github.com/Adelodunpeter25/vx/internal/buffer"
	"github.com/Adelodunpeter25/vx/internal/utils"
)

type Result struct {
	Quit       bool
	Message    string
	Error      error
	NewBuffer  *buffer.Buffer
	SwitchFile bool
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
		
		size, _ := buf.GetFileSize()
		msg := utils.FormatFileInfo(buf.Filename(), size, buf.LineCount())
		return Result{Message: msg}
	
	case "wq":
		if buf.Filename() == "" {
			return Result{Error: fmt.Errorf("no file name")}
		}
		
		if err := buf.Save(); err != nil {
			return Result{Error: err}
		}
		
		size, _ := buf.GetFileSize()
		msg := utils.FormatFileInfo(buf.Filename(), size, buf.LineCount())
		return Result{Quit: true, Message: msg}
	
	default:
		if strings.HasPrefix(cmd, "e ") {
			filename := strings.TrimSpace(cmd[2:])
			if filename == "" {
				return Result{Error: fmt.Errorf("no file name")}
			}
			
			// Check if current buffer is modified
			if buf.IsModified() {
				return Result{Error: fmt.Errorf("no write since last change (use :e! to override)")}
			}
			
			// Load new file
			newBuf, err := buffer.Load(filename)
			if err != nil {
				return Result{Error: err}
			}
			
			size, _ := newBuf.GetFileSize()
			msg := utils.FormatFileInfo(filename, size, newBuf.LineCount())
			return Result{NewBuffer: newBuf, SwitchFile: true, Message: msg}
		}
		
		if strings.HasPrefix(cmd, "w ") {
			filename := strings.TrimSpace(cmd[2:])
			if filename == "" {
				return Result{Error: fmt.Errorf("no file name")}
			}
			buf.SetFilename(filename)
			if err := buf.Save(); err != nil {
				return Result{Error: err}
			}
			
			size, _ := buf.GetFileSize()
			msg := utils.FormatFileInfo(filename, size, buf.LineCount())
			return Result{Message: msg}
		}
		
		if strings.HasPrefix(cmd, "wq ") {
			filename := strings.TrimSpace(cmd[3:])
			if filename == "" {
				return Result{Error: fmt.Errorf("no file name")}
			}
			buf.SetFilename(filename)
			if err := buf.Save(); err != nil {
				return Result{Error: err}
			}
			
			size, _ := buf.GetFileSize()
			msg := utils.FormatFileInfo(filename, size, buf.LineCount())
			return Result{Quit: true, Message: msg}
		}
		
		return Result{Error: fmt.Errorf("not an editor command: :%s", cmd)}
	}
}
