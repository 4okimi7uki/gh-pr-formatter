package cmd

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/4okimi7uki/gh-pr-formatter/internal/auth"
	"github.com/4okimi7uki/gh-pr-formatter/internal/client"
	"github.com/4okimi7uki/gh-pr-formatter/internal/gh"
	"github.com/4okimi7uki/gh-pr-formatter/internal/models"
	"github.com/4okimi7uki/gh-pr-formatter/internal/output"
	"github.com/4okimi7uki/gh-pr-formatter/internal/pkg/repository"
	"github.com/4okimi7uki/gh-pr-formatter/internal/pr"
	"github.com/4okimi7uki/gh-pr-formatter/internal/template"
	"github.com/4okimi7uki/gh-pr-formatter/internal/ui"
	"github.com/spf13/cobra"
)

var version = "v0.0.0-dev"
var showVersion bool

var rootCmd = &cobra.Command{
	Use:          "gh-pr-formatter",
	SilenceUsage: true,
	Short:        "Make your release work a little bit easier :P",
	Long:         "gh-pr-formatter formats merged GitHub pull requests for release notes.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if showVersion {
			resolvedVersion := gh.ResolvedVersion(version)
			fmt.Printf("gh-pr-formatter %s\n", resolvedVersion)

			// version check
			printCheckLatestVersion()
			return nil
		}

		return runMain(cmd)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolVarP(&showVersion, "version", "v", false, "Print version information")
	rootCmd.PersistentFlags().StringVarP(&repo, "repo", "r", "", "GitHub repository (owner/name)")
	rootCmd.PersistentFlags().IntVar(&limit, "limit", 50, "Number of PRs to fetch")
	rootCmd.PersistentFlags().StringSliceVarP(&excludePrefix, "exclude-prefix", "x", nil, "Exclude branch prefixes (repeatable, e.g. -x fix -x chore)")
}

var errNoPRs = errors.New("no prs")

func runMain(_ *cobra.Command) error {
	start := time.Now()
	token, err := auth.LoadGitHubToken()
	if err != nil {
		if err == auth.ErrTokenNotFound {
			return fmt.Errorf("not logged in. %s", ui.Mastered.Sprint("Run: gh-pr-formatter auth login"))
		}
		return err
	}

	ui.PrintLogo()

	var (
		prs  []models.MergedPrList
		from time.Time
	)

	err = WithSpinner("Resolving repo...", func(update func(string)) error {
		owner, repoName, err := repository.ResolveRepo(repo)
		if err != nil {
			return fmt.Errorf("resolve repo: %w", err)
		}

		update(" Checking environment...")
		if err := gh.CheckEnvironment(owner + "/" + repoName); err != nil {
			return err
		}

		update("PR fetching...")
		c := client.NewClient(token)
		prs, from, err = c.ListMergedPullRequests(owner, repoName, limit, excludePrefix...)
		if err != nil {
			return err
		}

		if len(prs) == 0 {
			return errNoPRs
		}

		update("PR formatting...")
		groupedPrs := pr.GroupedPrsByAuthor(prs)

		md, prList, err := template.BuildMarkdown(groupedPrs, from)
		if err != nil {
			return err
		}

		update("Build markdown...")
		writer := output.NewDefaultFilteWriter()
		path, err := writer.WriteReleaseMarkdown(md)
		if err != nil {
			return err
		}

		loc, _ := time.LoadLocation("Asia/Tokyo")
		ui.PrintSummary(os.Stdout, ui.Summary{
			Target: ui.Target{
				Repository:    owner + "/" + repoName,
				ExcludePrefix: excludePrefix,
				From:          from,
			},
			PRCount:     template.CountPRs(groupedPrs),
			AuthorCount: len(groupedPrs),
			OutputPath:  path,
			PRList:      prList,
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
	fmt.Printf("Done in %.1fs 📝✨\n", elapsed.Seconds())

	// version check
	if version != "v0.0.0-dev" {
		printCheckLatestVersion()
	}

	return nil
}

func printCheckLatestVersion() {
	if msg, err := gh.CheckLatestVersion("4okimi7uki", "gh-pr-formatter", version); err == nil && msg != "" {
		fmt.Fprintf(os.Stdout, "\n%s\n", ui.LimeYellow.Sprint(msg))
		fmt.Fprintf(os.Stdout, "%s\n", "https://github.com/4okimi7uki/gh-pr-formatter/releases")
	} else if err != nil {
		_ = err // or log.Printf("version check failed: %v", err)
	}
}
