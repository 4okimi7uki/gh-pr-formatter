package cmd

import (
	"fmt"
	"os"

	"github.com/4okimi7uki/gh-pr-formatter/internal/app"
	"github.com/4okimi7uki/gh-pr-formatter/internal/gh"
	"github.com/4okimi7uki/gh-pr-formatter/internal/pkg/repository"
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

		owner, repoName, _ := repository.ResolveRepo(repo)
		if owner == "" && repoName == "" {
			return fmt.Errorf("required flag '--repo' or '-r' not set\n(use -r owner/repo or run inside a git repository)")
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
	rootCmd.Flags().BoolVarP(&showVersion, "version", "v", false, "print version information")
	addCommonFlags(rootCmd)
}

func addCommonFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&repo, "repo", "r", "", "GitHub repository in owner/name format (required when not in a git repository)")
	cmd.Flags().IntVar(&limit, "limit", 50, "number of PRs to fetch")
	cmd.Flags().StringSliceVarP(&excludePrefix, "exclude-prefix", "x", nil, "exclude branch prefixes (repeatable, e.g. -x fix -x chore)")
}
