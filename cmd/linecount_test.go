package cmd

import (
	"os"
	"strings"
	"testing"
)

func TestLinecount(t *testing.T) {
	// Test 1: File input with -f
	t.Run("file input with -f", func(t *testing.T) {
		filePath := createTempFile(t, "linecount", sampleContent)
		defer os.Remove(filePath)

		output, err := linecount(filePath)
		if err != nil {
			t.Errorf("linecount returned error: %v", err)
		}
		if output != 4 {
			t.Errorf("expected output to be 4, but got: %d", output)
		}
	})

	// Test 2: Stdin input
	t.Run("stdin input", func(t *testing.T) {
		r, w, err := os.Pipe()
		if err != nil {
			t.Fatalf("failed to create pipe: %v", err)
		}
		go func() {
			_, _ = w.Write([]byte(sampleContent))
			w.Close()
		}()

		oldStdin := os.Stdin
		defer func() { os.Stdin = oldStdin }()
		os.Stdin = r

		output, err := linecount("-")
		if err != nil {
			t.Errorf("linecount returned error: %v", err)
		}
		if output != 4 {
			t.Errorf("expected output to be 4, but got: %d", output)
		}
	})

	// Test 3: Non-existent file
	t.Run("non-existent file", func(t *testing.T) {
		nonExist := "non_exist_file.txt"
		_, err := linecount(nonExist)
		if err == nil || !strings.Contains(err.Error(), "No such file") {
			t.Errorf("expected non-existent file error, but got: %v", err)
		}
	})

	// Test 4: Directory as input
	t.Run("directory as input", func(t *testing.T) {
		dir := os.TempDir()
		_, err := linecount(dir)
		if err == nil || !strings.Contains(err.Error(), "Expected file got directory") {
			t.Errorf("expected directory input error, but got: %v", err)
		}
	})

	// Test 5: Binary file input for linecount
	t.Run("binary file input for linecount", func(t *testing.T) {
		binaryContent := "binary\x00data"
		filePath := createBinaryTempFile(t, "binary", []byte(binaryContent))
		defer os.Remove(filePath)

		_, err := linecount(filePath)
		if err == nil || !strings.Contains(err.Error(), "cannot do linecount for binary file") {
			t.Errorf("expected binary file error, but got: %v", err)
		}
	})
}
