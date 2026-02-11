package gh

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/4okimi7uki/gh-pr-formatter/internal/ui"
)

type githubRelease struct {
	TagName string `json:"tag_name"`
}

func ResolvedVersion(version string) string {
	if version != "" && version != "v0.0.0-dev" && version != "dev" {
		return version
	}

	if bi, ok := debug.ReadBuildInfo(); ok {
		v := bi.Main.Version
		if v != "" && v != "(devel)" {
			return v
		}
	}

	return "v0.0.0-dev"
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
			return fmt.Sprintf("* A new version is available: %s → %s", ui.Lime.Sprint(version), ui.Lime.Sprint(latest)), nil
		}
	}
	return "", nil
}
