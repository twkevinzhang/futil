package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/urfave/cli/v2"
)

// 範例檔案內容（共 4 行）
var sampleContent = "how do\nyou\nturn this\non\n"

// 題目文件所提供的 checksum 預期結果（sampleContent 包含最後的換行符號）
var expectedMD5 = "a8c5d553ed101646036a811772ffbdd8"
var expectedSHA1 = "a656582ca3143a5f48718f4a15e7df018d286521"
var expectedSHA256 = "495a3496cfd90e68a53b5e3ff4f9833b431fe996298f5a28228240ee2a25c09d"

// captureOutput 暫時攔截標準輸出，並返回執行 f() 後的輸出內容
func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

// createTempFile 於系統暫存目錄建立檔案，寫入指定內容，並回傳檔案路徑
func createTempFile(t *testing.T, prefix, content string) string {
	tmpfile, err := ioutil.TempFile("", prefix)
	if err != nil {
		t.Fatalf("建立暫存檔案失敗: %v", err)
	}
	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatalf("寫入暫存檔案失敗: %v", err)
	}
	tmpfile.Close()
	return tmpfile.Name()
}

// createBinaryTempFile 建立一個含有 binary 內容（包含 null byte）的暫存檔案
func createBinaryTempFile(t *testing.T, prefix string, content []byte) string {
	tmpfile, err := ioutil.TempFile("", prefix)
	if err != nil {
		t.Fatalf("建立暫存檔案失敗: %v", err)
	}
	if _, err := tmpfile.Write(content); err != nil {
		t.Fatalf("寫入暫存檔案失敗: %v", err)
	}
	tmpfile.Close()
	return tmpfile.Name()
}

// newTestContext 依傳入參數建立一個 cli.Context 用於單元測試
func newTestContext(cmdName string, args []string) *cli.Context {
	set := flag.NewFlagSet(cmdName, 0)
	set.String("file", "", "the input file")
	set.String("f", "", "the input file")
	set.Bool("md5", false, "calculate MD5 checksum")
	set.Bool("sha1", false, "calculate SHA1 checksum")
	set.Bool("sha256", false, "calculate SHA256 checksum")
	_ = set.Parse(args)

	if set.Lookup("file").Value.String() == "" && set.Lookup("f").Value.String() != "" {
		_ = set.Set("file", set.Lookup("f").Value.String())
	}

	app := cli.NewApp()
	return cli.NewContext(app, set, nil)
}

// 以下為針對 lineCountCmd 的單元測試
func TestLineCountCmd(t *testing.T) {
	// 測試 1：檔案輸入，使用短旗 -f
	t.Run("file input with -f", func(t *testing.T) {
		filePath := createTempFile(t, "linecount", sampleContent)
		defer os.Remove(filePath)

		ctx := newTestContext("linecount", []string{"-f", filePath})
		output := captureOutput(func() {
			if err := lineCountCmd(ctx); err != nil {
				t.Errorf("lineCountCmd 回傳錯誤: %v", err)
			}
		})
		if !strings.Contains(output, "4") {
			t.Errorf("預期輸出包含 4，但得到: %s", output)
		}
	})

	// 測試 2：檔案輸入，使用長旗 --file
	t.Run("file input with --file", func(t *testing.T) {
		filePath := createTempFile(t, "linecount", sampleContent)
		defer os.Remove(filePath)

		ctx := newTestContext("linecount", []string{"--file", filePath})
		output := captureOutput(func() {
			if err := lineCountCmd(ctx); err != nil {
				t.Errorf("lineCountCmd 回傳錯誤: %v", err)
			}
		})
		if !strings.Contains(output, "4") {
			t.Errorf("預期輸出包含 4，但得到: %s", output)
		}
	})

	// 測試 3：從標準輸入讀取（-f -）
	t.Run("stdin input", func(t *testing.T) {
		// 透過 os.Pipe 建立一組管道，用以模擬標準輸入
		r, w, err := os.Pipe()
		if err != nil {
			t.Fatalf("建立 pipe 失敗: %v", err)
		}
		// 寫入 sampleContent 至管道的寫入端，然後關閉寫入端
		go func() {
			_, _ = w.Write([]byte(sampleContent))
			w.Close()
		}()

		// 保存原本的 os.Stdin，並於測試結束後還原
		oldStdin := os.Stdin
		defer func() { os.Stdin = oldStdin }()
		os.Stdin = r

		ctx := newTestContext("linecount", []string{"-f", "-"})
		output := captureOutput(func() {
			if err := lineCountCmd(ctx); err != nil {
				t.Errorf("lineCountCmd 回傳錯誤: %v", err)
			}
		})
		if !strings.Contains(output, "4") {
			t.Errorf("預期輸出包含 4，但得到: %s", output)
		}
	})

	// 測試 4：錯誤情況
	t.Run("non-existent file", func(t *testing.T) {
		nonExist := "non_exist_file.txt"
		ctx := newTestContext("linecount", []string{"-f", nonExist})
		err := lineCountCmd(ctx)
		if err == nil || !strings.Contains(err.Error(), "No such file") {
			t.Errorf("預期非存在檔案錯誤，但得到: %v", err)
		}
	})

	t.Run("directory as input", func(t *testing.T) {
		// 使用系統暫存目錄作為範例
		dir := os.TempDir()
		ctx := newTestContext("linecount", []string{"-f", dir})
		err := lineCountCmd(ctx)
		if err == nil || !strings.Contains(err.Error(), "Expected file got directory") {
			t.Errorf("預期傳入目錄錯誤，但得到: %v", err)
		}
	})

	t.Run("binary file input for linecount", func(t *testing.T) {
		// 建立一個包含 null byte 的 binary 檔案
		binaryContent := "binary\x00data"
		filePath := createBinaryTempFile(t, "binary", []byte(binaryContent))
		defer os.Remove(filePath)

		ctx := newTestContext("linecount", []string{"-f", filePath})
		err := lineCountCmd(ctx)
		if err == nil || !strings.Contains(err.Error(), "Cannot do linecount for binary file") {
			t.Errorf("預期 binary 檔案回傳錯誤，但得到: %v", err)
		}
	})
}

// 以下為針對 checksumCmd 的單元測試
func TestChecksumCmd(t *testing.T) {
	// 測試 1：檔案輸入，使用 --md5
	t.Run("file input with md5", func(t *testing.T) {
		filePath := createTempFile(t, "checksum", sampleContent)
		defer os.Remove(filePath)

		ctx := newTestContext("checksum", []string{"-f", filePath, "--md5"})
		output := captureOutput(func() {
			if err := checksumCmd(ctx); err != nil {
				t.Errorf("checksumCmd 回傳錯誤: %v", err)
			}
		})
		outTrim := strings.TrimSpace(output)
		if outTrim != expectedMD5 {
			t.Errorf("md5 checksum 不符合，預期 %s，得到 %s", expectedMD5, outTrim)
		}
	})

	t.Run("file input with sha1", func(t *testing.T) {
		filePath := createTempFile(t, "checksum", sampleContent)
		defer os.Remove(filePath)

		ctx := newTestContext("checksum", []string{"-f", filePath, "--sha1"})
		output := captureOutput(func() {
			if err := checksumCmd(ctx); err != nil {
				t.Errorf("checksumCmd 回傳錯誤: %v", err)
			}
		})
		outTrim := strings.TrimSpace(output)
		if outTrim != expectedSHA1 {
			t.Errorf("sha1 checksum 不符合，預期 %s，得到 %s", expectedSHA1, outTrim)
		}
	})

	t.Run("file input with sha256", func(t *testing.T) {
		filePath := createTempFile(t, "checksum", sampleContent)
		defer os.Remove(filePath)

		ctx := newTestContext("checksum", []string{"-f", filePath, "--sha256"})
		output := captureOutput(func() {
			if err := checksumCmd(ctx); err != nil {
				t.Errorf("checksumCmd 回傳錯誤: %v", err)
			}
		})
		outTrim := strings.TrimSpace(output)
		if outTrim != expectedSHA256 {
			t.Errorf("sha256 checksum 不符合，預期 %s，得到 %s", expectedSHA256, outTrim)
		}
	})

	// 測試 2：從標準輸入讀取（-f -）
	t.Run("stdin input with sha256", func(t *testing.T) {
		// 利用 os.Pipe 模擬標準輸入
		r, w, err := os.Pipe()
		if err != nil {
			t.Fatalf("建立 pipe 失敗: %v", err)
		}
		go func() {
			_, _ = w.Write([]byte(sampleContent))
			w.Close()
		}()

		oldStdin := os.Stdin
		defer func() { os.Stdin = oldStdin }()
		os.Stdin = r

		ctx := newTestContext("checksum", []string{"-f", "-", "--sha256"})
		output := captureOutput(func() {
			if err := checksumCmd(ctx); err != nil {
				t.Errorf("checksumCmd 回傳錯誤: %v", err)
			}
		})
		outTrim := strings.TrimSpace(output)
		if outTrim != expectedSHA256 {
			t.Errorf("stdin sha256 checksum 不符合，預期 %s，得到 %s", expectedSHA256, outTrim)
		}
	})

	// 測試 3：錯誤處理，針對非存在檔案與目錄
	t.Run("non-existent file", func(t *testing.T) {
		nonExist := "non_exist_file.txt"
		ctx := newTestContext("checksum", []string{"-f", nonExist, "--sha256"})
		err := checksumCmd(ctx)
		if err == nil || !strings.Contains(err.Error(), "No such file") {
			t.Errorf("預期非存在檔案錯誤，但得到: %v", err)
		}
	})

	t.Run("directory as input", func(t *testing.T) {
		dir := os.TempDir()
		ctx := newTestContext("checksum", []string{"-f", dir, "--sha256"})
		err := checksumCmd(ctx)
		if err == nil || !strings.Contains(err.Error(), "Expected file got directory") {
			t.Errorf("預期傳入目錄錯誤，但得到: %v", err)
		}
	})

	// 測試 4：binary 檔案的 checksum 計算（binary 檔案可正常計算 checksum）
	t.Run("binary file input for checksum (sha256)", func(t *testing.T) {
		binaryContent := []byte("binary\x00data")
		filePath := createBinaryTempFile(t, "binaryChecksum", binaryContent)
		defer os.Remove(filePath)

		// 計算預期的 sha256 值
		hash := sha256.Sum256(binaryContent)
		expected := hex.EncodeToString(hash[:])

		ctx := newTestContext("checksum", []string{"-f", filePath, "--sha256"})
		output := captureOutput(func() {
			if err := checksumCmd(ctx); err != nil {
				t.Errorf("checksumCmd 回傳錯誤: %v", err)
			}
		})
		outTrim := strings.TrimSpace(output)
		if outTrim != expected {
			t.Errorf("binary file sha256 checksum 不符合，預期 %s，得到 %s", expected, outTrim)
		}
	})
}
