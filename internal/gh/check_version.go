package gh

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type githubRelease struct {
	TagName string `json:"tag_name"`
}

func fetchLatestVersion(owner, repo string) (string, error) {
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
	latest, err := fetchLatestVersion(owner, repo)
	if err == nil && latest != "" {
		latestTrimmed := strings.TrimPrefix(latest, "v")
		currentTrimmed := strings.TrimPrefix(version, "v")

		if latestTrimmed != currentTrimmed {
			return fmt.Sprintf("* a new version of gh-pr-formatter is version available: %s --> %s", version, latest), nil
		}
	}
	return "", nil
}
