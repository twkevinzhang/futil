package cmd

import (
	"os"
	"testing"
)

// 範例檔案內容（共 4 行）
var sampleContent = "how do\nyou\nturn this\non\n"

// createTempFile 於系統暫存目錄建立檔案，寫入指定內容，並回傳檔案路徑
func createTempFile(t *testing.T, prefix, content string) string {
	tmpfile, err := os.CreateTemp("", prefix)
	defer func() {
		err := tmpfile.Close()
		if err != nil {
			t.Errorf("關閉暫存檔案失敗: %v", err)
		}
	}()
	if err != nil {
		t.Fatalf("建立暫存檔案失敗: %v", err)
	}
	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatalf("寫入暫存檔案失敗: %v", err)
	}
	return tmpfile.Name()
}

// createBinaryTempFile 建立一個含有 binary 內容（包含 null byte）的暫存檔案
func createBinaryTempFile(t *testing.T, prefix string, content []byte) string {
	tmpfile, err := os.CreateTemp("", prefix)
	defer func() {
		err := tmpfile.Close()
		if err != nil {
			t.Errorf("關閉暫存檔案失敗: %v", err)
		}
	}()
	if err != nil {
		t.Fatalf("建立暫存檔案失敗: %v", err)
	}
	if _, err := tmpfile.Write(content); err != nil {
		t.Fatalf("寫入暫存檔案失敗: %v", err)
	}
	return tmpfile.Name()
}
