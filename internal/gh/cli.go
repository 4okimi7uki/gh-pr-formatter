package gh

import (
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/4okimi7uki/gh-pr-formatter/internal/pr"
)

func CheckEnvironment() error {
	// git repo check
	out, err := exec.Command("git", "rev-parse", "--is-inside-work-tree").CombinedOutput()
	if err != nil || strings.TrimSpace(string(out)) != "true" {
		return errors.New("oops, this command must be run inside a Git repository")
	}

	// gh exists check
	if _, err := exec.LookPath("gh"); err != nil {
		msg := `The 'gh' command was not found.

It looks like GitHub CLI isn't installed yet.
You can install it via Homebrew:
  brew install gh

Or check the official download page:
  https://cli.github.com/
		`
		return errors.New(msg)
	}

	if err := exec.Command("gh", "auth", "status").Run(); err != nil {
		msg := `GitHub CLI is installed, but you're not authenticated.
Please run 'gh auth login' first.`
		return errors.New(msg)
	}

	return nil
}

func GetLastMergedDate() (time.Time, error) {
	cmdStr := `
	pr list
	--state merged
	--base main
	--limit 100
	--search -head:hotfix/
	--json number,mergedAt
	`

	getDateCmd := strings.Fields(cmdStr)
	cmd := exec.Command("gh", getDateCmd...)
	out, err := cmd.CombinedOutput()

	if err != nil {
		return time.Time{}, fmt.Errorf("gh pr list failed: %v\n%s", err, string(out))
	}

	var prs []pr.PrList
	if err := json.Unmarshal(out, &prs); err != nil {
		return time.Time{}, fmt.Errorf("failed to decode JSON: %w", err)
	}

	if len(prs) == 0 {
		return time.Time{}, fmt.Errorf("no merged PRs")
	}

	var last time.Time
	for _, pr := range prs {
		if pr.MergedAt.After(last) {
			last = pr.MergedAt
		}
	}

	if last.IsZero() {
		return time.Time{}, fmt.Errorf("no merged PRs")
	}

	return last, nil
}

func GetPrList() ([]pr.MergedPrList, error) {
	getPrListCmdStr := `
	pr list
	--state merged
	--base develop
	--search merged:%s..%s
	--limit 300
	--json number,title,author,mergedAt
	`

	_from, err := GetLastMergedDate()
	if err != nil {
		return []pr.MergedPrList{}, fmt.Errorf("getLastMergedDate failed: %w", err)
	}

	from := _from.UTC().Format(time.RFC3339)
	to := time.Now().UTC().Format(time.RFC3339)

	arg := strings.Fields(fmt.Sprintf(getPrListCmdStr, from, to))
	out, err := exec.Command("gh", arg...).CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("gh pr list failed: %v\n%s", err, string(out))
	}

	var mergedPr []pr.MergedPrList
	if err := json.Unmarshal(out, &mergedPr); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	if len(mergedPr) == 0 {
		return nil, fmt.Errorf("No Merged PRs")
	}

	return mergedPr, nil
}
