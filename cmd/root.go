package cmd

import (
	_ "embed"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

// version of the futil tool
const version = "0.0.1"

//go:embed rootcmd_usage_template.txt
var rootcmdUsageTemplate string

//go:embed subcmd_usage_template.txt
var subcmdUsageTemplate string

func init() {
	rootCmd.PersistentFlags().BoolP("help", "h", false, "help for futil")
	rootCmd.SetUsageTemplate(rootcmdUsageTemplate)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "futil",
	Short: "File Utility",
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true, // removes cmd
	},
}
