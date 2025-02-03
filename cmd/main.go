package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"io"
	"os"
)

// version of the futil tool
const version = "0.0.1"

type hashWriter interface {
	io.Writer
	Sum([]byte) []byte
}

func main() {
	app := &cli.App{
		Name:    "futil",
		Usage:   "File Utility",
		Version: version,
		Commands: []*cli.Command{
			{
				Name:      "linecount",
				Usage:     "Print line count of file",
				UsageText: "futil linecount [flags]",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "file",
						Aliases: []string{"f"},
						Usage:   "the input file",
					},
				},
				Action: lineCountCmd,
			},
			{
				Name:      "checksum",
				Usage:     "Print checksum of file",
				UsageText: "futil checksum [flags]",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "file",
						Aliases: []string{"f"},
						Usage:   "the input file",
					},
					&cli.BoolFlag{
						Name:  "md5",
						Usage: "calculate MD5 checksum",
					},
					&cli.BoolFlag{
						Name:  "sha1",
						Usage: "calculate SHA1 checksum",
					},
					&cli.BoolFlag{
						Name:  "sha256",
						Usage: "calculate SHA256 checksum",
					},
				},
				Action: checksumCmd,
			},
			{
				Name:      "version",
				Usage:     "Show the version info",
				UsageText: "futil version",
				Action: func(c *cli.Context) error {
					fmt.Printf("futil v%s\n", version)
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
