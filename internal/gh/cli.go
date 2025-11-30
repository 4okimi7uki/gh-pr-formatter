package gh

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	model "github.com/4okimi7uki/gh-pr-formatter/internal/model"
)

var (
	ErrNotGitRepo  = errors.New("not a git repository")
	ErrGhNotFound  = errors.New("gh command not found")
	ErrGhNotAuthed = errors.New("gh not authenticated")
)

func CheckEnvironment(repo string) error {

	if repo == "" {
		// git repo check
		out, err := exec.Command("git", "rev-parse", "--is-inside-work-tree").CombinedOutput()
		if err != nil || strings.TrimSpace(string(out)) != "true" {
			return ErrNotGitRepo
		}
	}
	// gh exists check
	if _, err := exec.LookPath("gh"); err != nil {
		return ErrGhNotFound
	}

	if err := exec.Command("gh", "auth", "status").Run(); err != nil {
		return ErrGhNotAuthed
	}

	return nil
}

func PrintHelp(err error) {
	switch err {
	case ErrNotGitRepo:
		fmt.Fprintln(os.Stderr, "this command must be run in a Git repository or with the '--repo owner/repo' option")
	case ErrGhNotFound:
		fmt.Fprintln(os.Stderr, "gh command not found")
		fmt.Println(`
It looks like GitHub CLI isn't installed yet.
You can install it via Homebrew:
  brew install gh

Or check the official download page:
  https://cli.github.com/`)
	case ErrGhNotAuthed:
		fmt.Fprintln(os.Stderr, "gh is installed, but you're not authenticated")
		fmt.Fprintln(os.Stderr, "Please run 'gh auth login' first")
	}
}

func GetLastMergedDate(repo string) (time.Time, error) {
	cmdStr := `
	pr list
	--state merged
	--base main
	--limit 100
	--search -head:hotfix/
	--json number,mergedAt
	`

	getDateCmd := strings.Fields(cmdStr)

	if repo != "" {
		getDateCmd = append(getDateCmd, "--repo", repo)
	}

	cmd := exec.Command("gh", getDateCmd...)
	out, err := cmd.CombinedOutput()

	if err != nil {
		return time.Time{}, fmt.Errorf("gh pr list failed: %v\n%s", err, string(out))
	}

	var prs []model.PrList
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

func GetPrList(repo string) ([]model.MergedPrList, error) {
	getPrListCmdStr := `
	pr list
	--state merged
	--base develop
	--search merged:%s..%s
	--limit 300
	--json number,title,author,mergedAt
	`

	_from, err := GetLastMergedDate(repo)
	if err != nil {
		return []model.MergedPrList{}, fmt.Errorf("getLastMergedDate failed: %w", err)
	}

	from := _from.UTC().Format(time.RFC3339)
	to := time.Now().UTC().Format(time.RFC3339)

	arg := strings.Fields(fmt.Sprintf(getPrListCmdStr, from, to))

	if repo != "" {
		arg = append(arg, "--repo", repo)
	}

	out, err := exec.Command("gh", arg...).CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("gh pr list failed: %v\n%s", err, string(out))
	}

	var mergedPr []model.MergedPrList
	if err := json.Unmarshal(out, &mergedPr); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	if len(mergedPr) == 0 {
		return nil, fmt.Errorf("no Merged PRs")
	}

	return mergedPr, nil
}

type githubRelease struct {
	TagName string `json:"tag_name"`
}

func getLatestVersion(owner, repo string) (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", owner, repo)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", "gh-pr-formatter")

	client := &http.Client{
		Timeout: 3 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status: %s", resp.Status)
	}

	var release githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", err
	}
	return release.TagName, nil
}

func CheckLatestVersion(owner, repo, version string) (string, error) {
	latest, err := getLatestVersion(owner, repo)
	if err == nil && latest != "" {
		latestTrimmed := strings.TrimPrefix(latest, "v")
		currentTrimmed := strings.TrimPrefix(version, "v")

		if latestTrimmed != currentTrimmed {
			return fmt.Sprintf("a new version of gh-pr-formatter is version available: %s --> %s", version, latest), nil
		}
	}
	return "", nil
}
