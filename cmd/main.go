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
	cli.CommandHelpTemplate = readTemplate()
	cli.AppHelpTemplate = readTemplate()
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

func readTemplate() string {
	return `{{$v := offset .HelpName 6}}{{if .Usage}}{{wrap .Usage $v}}{{end}}

USAGE:
   {{if .UsageText}}{{wrap .UsageText 3}}{{else}}{{.HelpName}} {{if .VisibleFlags}}command [command options]{{end}}{{if .ArgsUsage}} {{.ArgsUsage}}{{else}}{{if .Args}} [arguments...]{{end}}{{end}}{{end}}{{if .Description}}

DESCRIPTION:
   {{template "descriptionTemplate" .}}{{end}}{{if .VisibleCommands}}

COMMANDS:{{template "visibleCommandCategoryTemplate" .}}{{end}}{{if .VisibleFlagCategories}}

OPTIONS:{{template "visibleFlagCategoryTemplate" .}}{{else if .VisibleFlags}}

OPTIONS:{{template "visibleFlagTemplate" .}}{{end}}
`
}
