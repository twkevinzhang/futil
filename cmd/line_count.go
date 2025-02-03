package main

import (
	"bufio"
	"fmt"
	"futil/internal/utils"
	"github.com/urfave/cli/v2"
)

func lineCountCmd(c *cli.Context) error {
	filename := c.String("file")
	if filename == "" {
		return cli.Exit("error: Please specify an input file using -f/--file", 1)
	}

	f, name, err := utils.OpenInput(filename)
	if err != nil {
		return cli.Exit(err.Error(), 1)
	}
	reader := bufio.NewReader(f)
	if filename != "-" && utils.IsBinary(reader) {
		f.Close()
		return cli.Exit(fmt.Sprintf("error: Cannot do linecount for binary file '%s'", name), 1)
	}

	scanner := bufio.NewScanner(reader)
	count := 0
	for scanner.Scan() {
		count++
	}
	if err := scanner.Err(); err != nil {
		f.Close()
		return cli.Exit(fmt.Sprintf("error: %v", err), 1)
	}
	if filename != "-" {
		f.Close()
	}
	fmt.Printf("%d\n", count)
	return nil
}
