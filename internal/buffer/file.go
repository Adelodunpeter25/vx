package buffer

import (
	"bufio"
	"fmt"
	"os"

	"github.com/Adelodunpeter25/vx/internal/utils"
)

func Load(filename string) (*Buffer, error) {
	// Check file size first
	tooLarge, size, err := utils.IsFileTooLarge(filename)
	if err != nil && !os.IsNotExist(err) {
		return nil, utils.NewFileError("load", filename, err)
	}
	
	if tooLarge {
		return nil, utils.NewFileError("load", filename, 
			fmt.Errorf("file too large (%d MB), maximum is %d MB", 
				size/(1024*1024), utils.MaxFileSize/(1024*1024)))
	}
	
	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			b := New()
			b.filename = filename
			return b, nil
		}
		return nil, utils.NewFileError("load", filename, err)
	}
	defer file.Close()

	b := &Buffer{
		filename: filename,
		lines:    []string{},
	}

	scanner := bufio.NewScanner(file)
	lineCount := 0
	
	for scanner.Scan() {
		lineCount++
		if lineCount > utils.MaxLines {
			return nil, utils.NewFileError("load", filename, 
				fmt.Errorf("too many lines (%d), maximum is %d", lineCount, utils.MaxLines))
		}
		
		// Validate and clean UTF-8
		line := scanner.Text()
		line = utils.ValidateUTF8(line)
		b.lines = append(b.lines, line)
	}

	if err := scanner.Err(); err != nil {
		// Try to recover - return what we loaded so far
		if len(b.lines) > 0 {
			b.modified = true
			return b, utils.NewFileError("load", filename, 
				fmt.Errorf("partial load: %v", err))
		}
		return nil, utils.NewFileError("load", filename, err)
	}

	if len(b.lines) == 0 {
		b.lines = []string{""}
	}

	return b, nil
}

func (b *Buffer) Save() error {
	if b.filename == "" {
		return fmt.Errorf("no filename set")
	}
	
	file, err := os.Create(b.filename)
	if err != nil {
		return utils.NewFileError("save", b.filename, err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for i, line := range b.lines {
		if i > 0 {
			if _, err := writer.WriteString("\n"); err != nil {
				return utils.NewFileError("save", b.filename, err)
			}
		}
		if _, err := writer.WriteString(line); err != nil {
			return utils.NewFileError("save", b.filename, err)
		}
	}
	
	if err := writer.Flush(); err != nil {
		return utils.NewFileError("save", b.filename, err)
	}

	b.modified = false
	return nil
}

// GetFileSize returns the size of the file on disk
func (b *Buffer) GetFileSize() (int64, error) {
	if b.filename == "" {
		return 0, nil
	}
	
	info, err := os.Stat(b.filename)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, nil
		}
		return 0, err
	}
	
	return info.Size(), nil
}
