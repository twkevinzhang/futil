package cmd

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"futil/internal/utils"
	"github.com/spf13/cobra"
	"io"
)

var (
	checksumCmdFile string
)

func init() {
	checksumCmd.Flags().StringVarP(&checksumCmdFile, "file", "f", "s", "the input file")
	linecountCmd.SetUsageTemplate(subcmdUsageTemplate)
	rootCmd.AddCommand(checksumCmd)
}

var checksumCmd = &cobra.Command{
	Use:   "checksum",
	Short: "Print checksum of file",
	RunE: func(cmd *cobra.Command, args []string) error {
		filename, _ := cmd.Flags().GetString("file")
		if filename == "" {
			return fmt.Errorf("please specify an input file using -f/--file")
		}
		md5, _ := cmd.Flags().GetBool("md5")
		sha1, _ := cmd.Flags().GetBool("sha1")
		sha256, _ := cmd.Flags().GetBool("sha256")
		s, e := checksum(filename, md5, sha1, sha256)
		if e != nil {
			return e
		}
		fmt.Println(s)
		return nil
	},
}

type hashWriter interface {
	io.Writer
	Sum([]byte) []byte
}

func checksum(filename string, md5Flag bool, sha1Flag bool, sha256Flag bool) (string, error) {
	algoCount := 0
	if md5Flag {
		algoCount++
	}
	if sha1Flag {
		algoCount++
	}
	if sha256Flag {
		algoCount++
	}
	if algoCount == 0 {
		return "", fmt.Errorf("please specify a checksum algorithm (--md5, --sha1, --sha256)")
	}
	if algoCount > 1 {
		return "", fmt.Errorf("please specify only one checksum algorithm")
	}

	f, _, err := utils.OpenInput(filename)
	if err != nil {
		return "", err
	}
	defer func() {
		if filename != "-" {
			f.Close()
		}
	}()

	var hashFunc hashWriter
	if md5Flag {
		hashFunc = md5.New()
	} else if sha1Flag {
		hashFunc = sha1.New()
	} else if sha256Flag {
		hashFunc = sha256.New()
	}

	if _, err := io.Copy(hashFunc, f); err != nil {
		return "", err
	}

	return hex.EncodeToString(hashFunc.Sum(nil)), nil
}
