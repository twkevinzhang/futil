package cmd

import (
	"bufio"
	_ "embed"
	"fmt"
	"futil/internal/utils"
	"github.com/spf13/cobra"
)

var (
	lineCountCmdFile string
)

func init() {
	linecountCmd.Flags().StringVarP(&lineCountCmdFile, "file", "f", "", "the input file")
	linecountCmd.SetUsageTemplate(subcmdUsageTemplate)
	rootCmd.AddCommand(linecountCmd)
}

var linecountCmd = &cobra.Command{
	Use:   "linecount",
	Short: "Print line count of file",
	RunE: func(cmd *cobra.Command, args []string) error {
		filename, _ := cmd.Flags().GetString("file")
		if filename == "" {
			return fmt.Errorf("please specify an input file using -f/--file")
		}
		i, e := linecount(filename)
		if e != nil {
			return e
		}
		fmt.Printf("%d\n", i)
		return nil
	},
}

func linecount(filename string) (int, error) {
	f, name, err := utils.OpenInput(filename)
	if err != nil {
		return 0, err
	}
	reader := bufio.NewReader(f)
	if filename != "-" && utils.IsBinary(reader) {
		f.Close()
		return 0, fmt.Errorf("cannot do linecount for binary file '%s'", name)
	}

	scanner := bufio.NewScanner(reader)
	count := 0
	for scanner.Scan() {
		count++
	}
	if err := scanner.Err(); err != nil {
		f.Close()
		return 0, err
	}
	if filename != "-" {
		f.Close()
	}
	return count, nil
}
