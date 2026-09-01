package buffer

import (
	"bufio"
	"fmt"
	"os"

	"github.com/Adelodunpeter25/vx/internal/undo"
	"github.com/Adelodunpeter25/vx/internal/utils"
)

func Load(filename string) (*Buffer, error) {
	// Single stat to avoid duplicate syscalls on startup path
	info, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			b := New()
			b.filename = filename
			return b, nil
		}
		return nil, utils.NewFileError("load", filename, err)
	}
	if info.IsDir() {
		return nil, utils.NewFileError("load", filename, fmt.Errorf("is a directory"))
	}
	size := info.Size()
	if size > utils.MaxFileSize {
		return nil, utils.NewFileError("load", filename,
			fmt.Errorf("file too large (%d MB), maximum is %d MB",
				size/(1024*1024), utils.MaxFileSize/(1024*1024)))
	}
	modTime := info.ModTime()

	file, err := os.Open(filename)
	if err != nil {
		return nil, utils.NewFileError("load", filename, err)
	}

	b := &Buffer{
		filename:  filename,
		lines:     []string{},
		undoStack: undo.NewStack(),
		modTime:   modTime,
		fileSize:  size,
	}

	// Fast path for small files: avoid counting lines (extra full read) entirely.
	// 256KB heuristic: with ~50 char avg line, 256KB ~ 5k lines (LazyLoadThreshold).
	// .zshrc and most dotfiles are <10KB, so this saves one full read.
	const smallFileThreshold = 256 * 1024
	useLazy := false
	var totalCount int

	if size > smallFileThreshold {
		// For larger files we need to know total lines to decide lazy vs full load
		// Count using the already-opened file (single fd, Seek back after)
		totalCount, err = utils.CountLinesFromFile(file)
		if err != nil {
			_ = file.Close()
			return nil, utils.NewFileError("load", filename, err)
		}
		if totalCount > utils.MaxLines {
			_ = file.Close()
			return nil, utils.NewFileError("load", filename,
				fmt.Errorf("too many lines (%d), maximum is %d", totalCount, utils.MaxLines))
		}
		if totalCount > utils.LazyLoadThreshold {
			useLazy = true
		}
		// file is now seeked to 0 and ready for either lazy or scanner
	}

	if useLazy {
		// Reuse the open file for lazy reader — no reopen, no recount
		reader := utils.NewLazyFileReaderWithFile(file, filename, totalCount)
		// file ownership transferred to reader; do not close here
		if reader.TotalCount() > utils.MaxLines {
			_ = reader.Close()
			return nil, utils.NewFileError("load", filename,
				fmt.Errorf("too many lines (%d), maximum is %d", reader.TotalCount(), utils.MaxLines))
		}
		chunk, err := reader.LoadChunk()
		if err != nil {
			_ = reader.Close()
			return nil, utils.NewFileError("load", filename, err)
		}
		if len(chunk) == 0 {
			b.lines = []string{""}
			b.totalLines = 1
			_ = reader.Close()
			return b, nil
		}
		b.lines = append(b.lines, chunk...)
		b.lazy = reader
		b.totalLines = reader.TotalCount()
		return b, nil
	}

	// Non-lazy: single-pass scan (no prior CountLines for small files)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	// Handle long lines (up to 1MB) without "token too long" error
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)
	lineCount := 0

	for scanner.Scan() {
		lineCount++
		if lineCount > utils.MaxLines {
			return nil, utils.NewFileError("load", filename,
				fmt.Errorf("too many lines (%d), maximum is %d", lineCount, utils.MaxLines))
		}
		line := scanner.Text()
		line = utils.ValidateUTF8(line)
		b.lines = append(b.lines, line)
	}

	if err := scanner.Err(); err != nil {
		if len(b.lines) > 0 {
			b.modified = true
			b.totalLines = len(b.lines)
			return b, utils.NewFileError("load", filename,
				fmt.Errorf("partial load: %v", err))
		}
		return nil, utils.NewFileError("load", filename, err)
	}

	if len(b.lines) == 0 {
		b.lines = []string{""}
		b.totalLines = 1
		return b, nil
	}

	b.totalLines = len(b.lines)
	return b, nil
}

func (b *Buffer) Save() error {
	if b.filename == "" {
		return fmt.Errorf("no filename set")
	}
	b.ensureAllLoaded()

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
	b.modVersion++
	if info, err := os.Stat(b.filename); err == nil {
		b.modTime = info.ModTime()
		b.fileSize = info.Size()
	}
	return nil
}

// GetFileSize returns the size of the file on disk (cached, no syscall if known)
func (b *Buffer) GetFileSize() (int64, error) {
	if b.filename == "" {
		return 0, nil
	}
	if b.fileSize > 0 {
		return b.fileSize, nil
	}
	info, err := os.Stat(b.filename)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, nil
		}
		return 0, err
	}
	b.fileSize = info.Size()
	return b.fileSize, nil
}
