package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, arg []string) {
		fmt.Printf("gh-pr-formatter version: %s\n", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
