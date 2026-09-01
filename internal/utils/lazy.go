package utils

import (
	"bufio"
	"os"
)

const (
	LazyLoadThreshold = 5000 // Lines threshold for lazy loading
	ChunkSize         = 1000 // Load this many lines at a time
)

// LazyFileReader reads file in chunks
type LazyFileReader struct {
	filename string
	file     *os.File
	scanner  *bufio.Scanner
	lines    []string
	loaded   int
	total    int
	eof      bool
}

// NewLazyFileReader creates a lazy file reader
func NewLazyFileReader(filename string) (*LazyFileReader, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	// Count total lines using the open file to avoid reopening
	total, err := CountLinesFromFile(file)
	if err != nil {
		file.Close()
		return nil, err
	}

	// File is already seeked to 0 by CountLinesFromFile, reuse it directly
	scanner := bufio.NewScanner(file)
	// Allow long lines up to 1MB
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	return &LazyFileReader{
		filename: filename,
		file:     file,
		scanner:  scanner,
		lines:    make([]string, 0, ChunkSize),
		total:    total,
	}, nil
}

// NewLazyFileReaderWithFile creates a lazy reader from an already-opened file
// and a pre-counted total. The file must be seeked to 0 and will be owned by the reader.
func NewLazyFileReaderWithFile(file *os.File, filename string, total int) *LazyFileReader {
	if file == nil {
		return nil
	}
	scanner := bufio.NewScanner(file)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)
	return &LazyFileReader{
		filename: filename,
		file:     file,
		scanner:  scanner,
		lines:    make([]string, 0, ChunkSize),
		total:    total,
	}
}

// LoadChunk loads next chunk of lines
func (r *LazyFileReader) LoadChunk() ([]string, error) {
	if r.eof {
		return nil, nil
	}

	chunk := make([]string, 0, ChunkSize)
	for i := 0; i < ChunkSize && r.scanner.Scan(); i++ {
		line := ValidateUTF8(r.scanner.Text())
		chunk = append(chunk, line)
		r.loaded++
	}

	if err := r.scanner.Err(); err != nil {
		return chunk, err
	}

	if len(chunk) < ChunkSize {
		r.eof = true
	}

	r.lines = append(r.lines, chunk...)
	return chunk, nil
}

// GetLine returns a line, loading if necessary
func (r *LazyFileReader) GetLine(n int) (string, error) {
	// Load chunks until we have the line
	for n >= len(r.lines) && !r.eof {
		_, err := r.LoadChunk()
		if err != nil {
			return "", err
		}
	}

	if n < len(r.lines) {
		return r.lines[n], nil
	}

	return "", nil
}

// GetLines returns all loaded lines
func (r *LazyFileReader) GetLines() []string {
	return r.lines
}

// LoadedCount returns number of loaded lines
func (r *LazyFileReader) LoadedCount() int {
	return r.loaded
}

// TotalCount returns total line count
func (r *LazyFileReader) TotalCount() int {
	return r.total
}

// IsFullyLoaded checks if all lines are loaded
func (r *LazyFileReader) IsFullyLoaded() bool {
	return r.eof
}

// Close closes the file
func (r *LazyFileReader) Close() error {
	if r.file != nil {
		return r.file.Close()
	}
	return nil
}

// ShouldUseLazyLoad determines if lazy loading should be used
func ShouldUseLazyLoad(filename string) (bool, error) {
	count, err := CountLines(filename)
	if err != nil {
		return false, err
	}
	return count > LazyLoadThreshold, nil
}
