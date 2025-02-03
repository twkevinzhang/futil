package utils

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// IsBinary checks whether the beginning of the stream looks like binary.
// It peeks at up to 512 bytes and returns true if any null byte is found.
func IsBinary(r *bufio.Reader) bool {
	const sniffLen = 512
	buf, err := r.Peek(sniffLen)
	if err != nil && err != io.EOF {
		// 無法 Peek 時，認定為 binary
		return true
	}
	for _, b := range buf {
		if b == 0 {
			return true
		}
	}
	return false
}

// OpenInput opens the input file or returns os.Stdin if filename == "-"
func OpenInput(filename string) (io.ReadCloser, string, error) {
	if filename == "-" {
		return os.Stdin, "stdin", nil
	}

	// 判斷檔案是否存在
	stat, err := os.Stat(filename)
	if err != nil {
		return nil, filename, fmt.Errorf("error: No such file '%s'", filename)
	}
	if stat.IsDir() {
		return nil, filename, fmt.Errorf("error: Expected file got directory '%s'", filename)
	}
	f, err := os.Open(filename)
	if err != nil {
		return nil, filename, fmt.Errorf("error: %v", err)
	}
	return f, filepath.Base(filename), nil
}
