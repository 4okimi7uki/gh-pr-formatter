package app

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

var ErrNoPRs = errors.New("no prs")

func IsAuthCommand(cmd *cobra.Command) bool {
	for c := cmd; c != nil; c = c.Parent() {
		if c.Name() == "auth" {
			return true
		}
	}
	return false
}

func RunMain(_ *cobra.Command, version string, repo string, limit int, excludePrefix []string) error {
	start := time.Now()
	token, err := auth.LoadGitHubToken()
	if err != nil {
		if err == auth.ErrTokenNotFound {
			return fmt.Errorf("not logged in. %s", ui.Mastered("Run: gh-pr-formatter auth login"))
		}
		return err
	}

	ui.PrintLogo()

	var (
		prs  []models.MergedPrList
		from time.Time
	)

	err = ui.WithSpinner(" Resolving repo...", func(update func(string)) error {
		owner, repoName, err := repository.ResolveRepo(repo)
		if err != nil {
			return fmt.Errorf("resolve repo: %w", err)
		}
		err = gh.ValidateRepository(token, owner, repoName)
		if err != nil {
			return err
		}

		update(" Checking environment...")
		if err := gh.CheckEnvironment(owner + "/" + repoName); err != nil {
			return err
		}

		update(" PR fetching...")
		c := client.NewClient(token)
		prs, from, err = c.ListMergedPullRequests(owner, repoName, limit, excludePrefix...)
		if err != nil {
			return err
		}

		if len(prs) == 0 {
			return ErrNoPRs
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
		if errors.Is(err, ErrNoPRs) {
			fmt.Println("\nNo merged pull requests found on 'develop' since last 'main' merge.")
			return nil
		}

		return err
	}
	elapsed := time.Since(start)
	fmt.Printf("Done in %.1fs 📝✨\n", elapsed.Seconds())

	// version check
	if version != "v0.0.0-dev" {
		PrintCheckLatestVersion(version)
	}

	return nil
}

func PrintCheckLatestVersion(version string) {
	resolvedVersion := gh.ResolvedVersion(version)
	if msg, err := gh.CheckLatestVersion("4okimi7uki", "gh-pr-formatter", resolvedVersion); err == nil && msg != "" {
		_, _ = fmt.Fprintf(os.Stdout, "%s\n", ui.LimeYellow(msg))
		_, _ = fmt.Fprintf(os.Stdout, "%s\n\n", "https://github.com/4okimi7uki/gh-pr-formatter/releases")
	} else if err != nil {
		_ = err // or log.Printf("version check failed: %v", err)
	}
}
