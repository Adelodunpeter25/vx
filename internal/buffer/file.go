package buffer

import (
	"bufio"
	"os"
)

func Load(filename string) (*Buffer, error) {
	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			b := New()
			b.filename = filename
			return b, nil
		}
		return nil, err
	}
	defer file.Close()

	b := &Buffer{
		filename: filename,
		lines:    []string{},
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		b.lines = append(b.lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if len(b.lines) == 0 {
		b.lines = []string{""}
	}

	return b, nil
}

func (b *Buffer) Save() error {
	file, err := os.Create(b.filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for i, line := range b.lines {
		if i > 0 {
			writer.WriteString("\n")
		}
		writer.WriteString(line)
	}
	writer.Flush()

	b.modified = false
	return nil
}
