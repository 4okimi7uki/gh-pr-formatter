package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/4okimi7uki/gh-pr-formatter/internal/auth"
	"github.com/4okimi7uki/gh-pr-formatter/internal/client"
	"github.com/4okimi7uki/gh-pr-formatter/internal/pkg/repository"
	"github.com/spf13/cobra"
)

var (
	repo          string
	limit         int
	excludePrefix []string
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List merged pull requests",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("repo: ", repo)
		fmt.Println("limit: ", limit)
		fmt.Println("exclude: ", excludePrefix)

		token, err := auth.LoadGitHubToken()
		if err != nil {
			if err == auth.ErrTokenNotFound {
				return fmt.Errorf("not logged in. Run: gh-pr-formatter auth login")
			}
		}

		owner, repoName, err := repository.ResolveRepo(repo)
		if err != nil {
			fmt.Fprintf(os.Stderr,
				"Error: '%s'\n", err)
			os.Exit(1)
		}

		c := client.NewClient(token)
		prs, from, err := c.ListMergedPullRequests(owner, repoName, 50, excludePrefix...)

		loc, _ := time.LoadLocation("Asia/Tokyo")

		fmt.Println("\nMerged PRs into develop")
		fmt.Printf("period: %s → Now\n", from.In(loc).Format("2006-01-02 15:04 MST"))

		const (
			noWidth     = 4
			authorWidth = 17
			branchWidth = 31
		)
		fmt.Printf(
			"\n%-*s  %-*s %-*s %s\n",
			noWidth, "NO",
			authorWidth, "AUTHOR",
			branchWidth, "BRANCH",
			"TITLE",
		)
		if len(prs) > 0 {
			for _, pr := range prs {
				fmt.Printf("#%-4d @%-16s %-30s %s\n",
					pr.Number,
					pr.Author.Login,
					pr.HeadRefName,
					pr.Title,
				)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().StringVar(&repo, "repo", "", "GitHub repository (owner/name)")
	listCmd.Flags().IntVar(&limit, "limit", 50, "Number of PRs to fetch")
	listCmd.Flags().StringSliceVar(&excludePrefix, "x", nil, "Exclude branch prefixes (e.g. fix)")
}
