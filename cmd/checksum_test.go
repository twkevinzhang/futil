package cmd

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"strings"
	"testing"
)

// 題目文件所提供的 checksum 預期結果（sampleContent 包含最後的換行符號）
var expectedMD5 = "a8c5d553ed101646036a811772ffbdd8"
var expectedSHA1 = "a656582ca3143a5f48718f4a15e7df018d286521"
var expectedSHA256 = "495a3496cfd90e68a53b5e3ff4f9833b431fe996298f5a28228240ee2a25c09d"

func TestChecksum(t *testing.T) {
	// Test 1: File input with md5
	t.Run("file input with md5", func(t *testing.T) {
		filePath := createTempFile(t, "checksum", sampleContent)
		defer os.Remove(filePath)

		output, err := checksum(filePath, true, false, false)
		if err != nil {
			t.Errorf("checksum returned error: %v", err)
		}
		if output != expectedMD5 {
			t.Errorf("md5 checksum does not match, expected %s, got %s", expectedMD5, output)
		}
	})

	// Test 2: File input with sha1
	t.Run("file input with sha1", func(t *testing.T) {
		filePath := createTempFile(t, "checksum", sampleContent)
		defer os.Remove(filePath)

		output, err := checksum(filePath, false, true, false)
		if err != nil {
			t.Errorf("checksum returned error: %v", err)
		}
		if output != expectedSHA1 {
			t.Errorf("sha1 checksum does not match, expected %s, got %s", expectedSHA1, output)
		}
	})

	// Test 3: File input with sha256
	t.Run("file input with sha256", func(t *testing.T) {
		filePath := createTempFile(t, "checksum", sampleContent)
		defer os.Remove(filePath)

		output, err := checksum(filePath, false, false, true)
		if err != nil {
			t.Errorf("checksum returned error: %v", err)
		}
		if output != expectedSHA256 {
			t.Errorf("sha256 checksum does not match, expected %s, got %s", expectedSHA256, output)
		}
	})

	// Test 4: Stdin input with sha256
	t.Run("stdin input with sha256", func(t *testing.T) {
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

		output, err := checksum("-", false, false, true)
		if err != nil {
			t.Errorf("checksum returned error: %v", err)
		}
		if output != expectedSHA256 {
			t.Errorf("stdin sha256 checksum does not match, expected %s, got %s", expectedSHA256, output)
		}
	})

	// Test 5: Non-existent file
	t.Run("non-existent file", func(t *testing.T) {
		nonExist := "non_exist_file.txt"
		_, err := checksum(nonExist, false, false, true)
		if err == nil || !strings.Contains(err.Error(), "No such file") {
			t.Errorf("expected non-existent file error, but got: %v", err)
		}
	})

	// Test 6: Directory as input
	t.Run("directory as input", func(t *testing.T) {
		dir := os.TempDir()
		_, err := checksum(dir, false, false, true)
		if err == nil || !strings.Contains(err.Error(), "Expected file got directory") {
			t.Errorf("expected directory input error, but got: %v", err)
		}
	})

	// Test 7: Binary file input for checksum (sha256)
	t.Run("binary file input for checksum (sha256)", func(t *testing.T) {
		binaryContent := []byte("binary\x00data")
		filePath := createBinaryTempFile(t, "binaryChecksum", binaryContent)
		defer os.Remove(filePath)

		hash := sha256.Sum256(binaryContent)
		expected := hex.EncodeToString(hash[:])

		output, err := checksum(filePath, false, false, true)
		if err != nil {
			t.Errorf("checksum returned error: %v", err)
		}
		if output != expected {
			t.Errorf("binary file sha256 checksum does not match, expected %s, got %s", expected, output)
		}
	})
}
