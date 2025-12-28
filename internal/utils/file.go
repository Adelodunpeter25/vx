package utils

import (
	"os"
	"unicode/utf8"
)

const (
	MaxFileSize = 100 * 1024 * 1024 // 100MB
	MaxLines    = 1000000            // 1 million lines
)

// IsFileTooLarge checks if a file exceeds safe size limits
func IsFileTooLarge(filename string) (bool, int64, error) {
	info, err := os.Stat(filename)
	if err != nil {
		return false, 0, err
	}
	
	size := info.Size()
	return size > MaxFileSize, size, nil
}

// ValidateUTF8 checks if a string is valid UTF-8 and returns a cleaned version
func ValidateUTF8(s string) string {
	if utf8.ValidString(s) {
		return s
	}
	
	// Replace invalid UTF-8 sequences with replacement character
	return cleanInvalidUTF8(s)
}

func cleanInvalidUTF8(s string) string {
	var result []rune
	for len(s) > 0 {
		r, size := utf8.DecodeRuneInString(s)
		if r == utf8.RuneError && size == 1 {
			// Invalid UTF-8, replace with replacement character
			result = append(result, '�')
		} else {
			result = append(result, r)
		}
		s = s[size:]
	}
	return string(result)
}

// CountLines estimates line count without loading entire file
func CountLines(filename string) (int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	buf := make([]byte, 32*1024)
	count := 0
	
	for {
		n, err := file.Read(buf)
		if n == 0 {
			break
		}
		
		for i := 0; i < n; i++ {
			if buf[i] == '\n' {
				count++
			}
		}
		
		if err != nil {
			break
		}
	}
	
	return count, nil
}
