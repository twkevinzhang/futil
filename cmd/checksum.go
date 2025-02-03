package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"futil/internal/utils"
	"github.com/urfave/cli/v2"
	"io"
)

func checksumCmd(c *cli.Context) error {
	filename := c.String("file")
	if filename == "" {
		return cli.Exit("error: Please specify an input file using -f/--file", 1)
	}

	md5Flag := c.Bool("md5")
	sha1Flag := c.Bool("sha1")
	sha256Flag := c.Bool("sha256")
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
		return cli.Exit("error: Please specify a checksum algorithm (--md5, --sha1, --sha256)", 1)
	}
	if algoCount > 1 {
		return cli.Exit("error: Please specify only one checksum algorithm", 1)
	}

	f, _, err := utils.OpenInput(filename)
	if err != nil {
		return cli.Exit(err.Error(), 1)
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
		return cli.Exit(fmt.Sprintf("error: %v", err), 1)
	}

	fmt.Println(hex.EncodeToString(hashFunc.Sum(nil)))
	return nil
}
