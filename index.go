package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type PrList struct {
	Number   int       `json:"number"`
	MergedAt time.Time `json:"mergedAt"`
}

type Author struct {
	ID    string `json:"id"`
	IsBot bool   `json:"is_bot"`
	Login string `json:"login"`
	Name  string `json:"name"`
}

type MergedPrList struct {
	Author Author `json:"author"`
	PrList
	Title string `json:"title"`
}

func checkGh() error {
	out, err := exec.Command("git", "rev-parse", "--is-inside-work-tree").CombinedOutput()
	if err != nil || strings.TrimSpace(string(out)) != "true" {
		return errors.New("Opps, This command must be run inside a Git repository")
	}

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

func getLastMergedDate() (time.Time, error) {
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
		fmt.Printf("failed: %v\n", err)
		panic(err)
	}

	var prs []PrList
	if err := json.Unmarshal(out, &prs); err != nil {
		panic(err)
	}

	sort.Slice(prs, func(a, b int) bool {
		return prs[a].MergedAt.Before(prs[b].MergedAt)
	})

	if len(prs) == 0 {
		return time.Time{}, fmt.Errorf("no merged PRs")
	}

	last := prs[len(prs)-1]
	fmt.Printf("last merged: #%d at %s\n", last.Number, last.MergedAt)

	return last.MergedAt, nil
}

func getPrList() ([]MergedPrList, error) {
	getPrListCmdStr := `
	pr list
	--state merged
	--base develop
	--search merged:%s..%s
	--limit 300
	--json number,title,author,mergedAt
	`

	_from, err := getLastMergedDate()
	if err != nil {
		return []MergedPrList{}, fmt.Errorf("getLastMergedDate failed: %w", err)
	}

	from := _from.UTC().Format(time.RFC3339)
	to := time.Now().UTC().Format(time.RFC3339)

	arg := strings.Fields(fmt.Sprintf(getPrListCmdStr, from, to))
	out, err := exec.Command("gh", arg...).CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("gh pr list failed: %v\n%s", err, string(out))
	}

	var mergedPr []MergedPrList
	if err := json.Unmarshal(out, &mergedPr); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	if len(mergedPr) == 0 {
		return nil, fmt.Errorf("No Merged PRs")
	}

	return mergedPr, nil
}

func main() {
	if err := checkGh(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	mergedPr, err := getPrList()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println(mergedPr)
}
