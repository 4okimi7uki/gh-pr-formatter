package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var version = "v0.0.0-dev"

var rootCmd = &cobra.Command{
	Use:     "gh-pr-formatter",
	Short:   "Make your release work a little bit easier :P",
	Long:    "gh-pr-formatter formats merged GitHub pull requests for release notes.",
	Version: version,
}

func Excute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
