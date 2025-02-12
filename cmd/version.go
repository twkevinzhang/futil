package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show the version info",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("futil v%s\n", version)
	},
}
