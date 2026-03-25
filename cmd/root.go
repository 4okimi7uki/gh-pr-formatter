package cmd

import (
	"fmt"
	"os"

	"github.com/4okimi7uki/gh-pr-formatter/internal/app"
	"github.com/4okimi7uki/gh-pr-formatter/internal/gh"
	"github.com/spf13/cobra"
)

var version = "v0.0.0-dev"
var showVersion bool

var (
	repo          string
	limit         int
	excludePrefix []string
)

var rootCmd = &cobra.Command{
	Use:          "gh-pr-formatter",
	SilenceUsage: true,
	Short:        "Make your release work a little bit easier :P",
	Long:         "gh-pr-formatter formats merged GitHub pull requests for release notes.",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if showVersion {
			return nil
		}

		if app.IsAuthCommand(cmd) {
			return nil
		}

		if repo == "" {
			return fmt.Errorf("required flag '--repo' or '-r' not set")
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if showVersion {
			resolvedVersion := gh.ResolvedVersion(version)
			fmt.Printf("gh-pr-formatter %s\n", resolvedVersion)

			// version check
			app.PrintCheckLatestVersion(version)
			return nil
		}

		return app.RunMain(cmd, version, repo, limit, excludePrefix)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolVarP(&showVersion, "version", "v", false, "Print version information")
	addCommonFlags(rootCmd)
}

func addCommonFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&repo, "repo", "r", "", "GitHub repository (owner/name) (required)")
	cmd.Flags().IntVar(&limit, "limit", 50, "Number of PRs to fetch")
	cmd.Flags().StringSliceVarP(&excludePrefix, "exclude-prefix", "x", nil, "Exclude branch prefixes (repeatable, e.g. -x fix -x chore)")
}
