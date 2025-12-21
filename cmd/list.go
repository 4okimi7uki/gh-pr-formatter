package cmd

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/4okimi7uki/gh-pr-formatter/internal/auth"
	"github.com/4okimi7uki/gh-pr-formatter/internal/client"
	"github.com/4okimi7uki/gh-pr-formatter/internal/gh"
	"github.com/4okimi7uki/gh-pr-formatter/internal/pkg/repository"
	"github.com/4okimi7uki/gh-pr-formatter/internal/ui"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List merged pull requests",
	RunE: func(cmd *cobra.Command, args []string) error {
		start := time.Now()
		token, err := auth.LoadGitHubToken()
		if err != nil {
			if err == auth.ErrTokenNotFound {
				return fmt.Errorf("not logged in. Run: gh-pr-formatter auth login")
			}
		}

		err = WithSpinner("Resolving repo...", func(update func(string)) error {
			owner, repoName, err := repository.ResolveRepo(repo)
			if err != nil {
				fmt.Fprintf(os.Stderr,
					"Error: '%s'\n", err)
				os.Exit(1)
			}

			update(" Checking environment...")
			if err := gh.CheckEnvironment(owner + "/" + repoName); err != nil {
				return err
			}

			update("PR fetching...")
			c := client.NewClient(token)
			prs, from, err := c.ListMergedPullRequests(owner, repoName, 50, excludePrefix...)
			if err != nil {
				return err
			}
			if len(prs) == 0 {
				return errNoPRs
			}

			loc, _ := time.LoadLocation("Asia/Tokyo")
			ui.PrintListSummary(os.Stdout, ui.ListSummary{
				Target: ui.Target{
					Repository:    owner + "/" + repoName,
					ExcludePrefix: excludePrefix,
					From:          from,
				},
				Prs: prs,
			}, loc)

			return nil
		})
		if err != nil {
			if errors.Is(err, errNoPRs) {
				fmt.Println("\nNo merged pull requests found on 'develop' since last 'main' merge.")
				return nil
			}

			return err
		}

		elapsed := time.Since(start)
		fmt.Printf("Done in %.1fs 🧸✨\n", elapsed.Seconds())

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
