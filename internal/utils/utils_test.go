package utils

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestIsBinary(t *testing.T) {
	// 測試一般文字資料，預期回傳 false
	text := "Hello, this is a text file.\nIt has multiple lines.\n"
	reader := bufio.NewReader(strings.NewReader(text))
	if IsBinary(reader) {
		t.Errorf("IsBinary() = true, want false for text input")
	}

	// 測試包含 null byte 的資料，預期回傳 true
	binaryData := []byte("Hello\x00World")
	reader = bufio.NewReader(bytes.NewReader(binaryData))
	if !IsBinary(reader) {
		t.Errorf("IsBinary() = false, want true for binary input with null byte")
	}

	// 測試當 Peek 發生錯誤時，預期回傳 true
	// 自訂一個 reader，使其 Peek 永遠回傳錯誤
	errReader := &errorReader{}
	bufReader := bufio.NewReader(errReader)
	if !IsBinary(bufReader) {
		t.Errorf("IsBinary() = false, want true when Peek error occurs")
	}
}

type errorReader struct{}

func (e *errorReader) Read(p []byte) (n int, err error) {
	return 0, io.ErrUnexpectedEOF
}

func TestOpenInput(t *testing.T) {
	// 測試當 filename 為 "-" 時，應回傳 os.Stdin
	f, name, err := OpenInput("-")
	if err != nil {
		t.Errorf("OpenInput(\"-\") returned error: %v", err)
	}
	if f != os.Stdin {
		t.Errorf("OpenInput(\"-\") did not return os.Stdin")
	}
	if name != "stdin" {
		t.Errorf("OpenInput(\"-\") name = %s, want 'stdin'", name)
	}

	// 測試開啟一個存在的暫存檔案
	tmpDir := os.TempDir()
	tmpFile := filepath.Join(tmpDir, "test_openinput.txt")
	content := []byte("temporary file content")
	if err := os.WriteFile(tmpFile, content, 0644); err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tmpFile)

	f, name, err = OpenInput(tmpFile)
	if err != nil {
		t.Errorf("OpenInput(%s) returned error: %v", tmpFile, err)
	}
	// 檢查回傳的檔案名稱是否為 basename
	wantName := filepath.Base(tmpFile)
	if name != wantName {
		t.Errorf("OpenInput(%s) returned name = %s, want %s", tmpFile, name, wantName)
	}
	// 測試讀取檔案內容是否正確
	data, err := io.ReadAll(f)
	f.Close()
	if err != nil {
		t.Errorf("Failed to read from file: %v", err)
	}
	if !bytes.Equal(data, content) {
		t.Errorf("File content mismatch: got %q, want %q", data, content)
	}

	// 測試傳入不存在的檔案，預期錯誤
	nonExistFile := "non_exist_file.txt"
	_, _, err = OpenInput(nonExistFile)
	if err == nil {
		t.Errorf("OpenInput(%s) did not return error for non-existent file", nonExistFile)
	}

	// 測試傳入一個目錄，預期錯誤
	_, _, err = OpenInput(tmpDir)
	if err == nil {
		t.Errorf("OpenInput(%s) did not return error for a directory", tmpDir)
	}
}
