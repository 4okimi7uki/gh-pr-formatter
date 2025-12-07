package repository

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"strings"
)

func splitOwnerRepo(slug string) (string, string, error) {
	slug = strings.TrimSuffix(slug, ".git")
	owner, name, ok := strings.Cut(slug, "/")
	if !ok || owner == "" || name == "" {
		return "", "", fmt.Errorf("invalid repo slug: %q", slug)
	}

	return owner, name, nil
}

func repoFromEnv() (string, string, error) {
	if v := os.Getenv("GITHUB_REPOSITORY"); v != "" {
		return splitOwnerRepo(v)
	}
	return "", "", fmt.Errorf("GITHUB_REPOSITORY not set")
}

func repoFromGit() (string, string, error) {
	out, err := exec.Command("git", "config", "--get", "remote.origin.url").Output()
	if err != nil {
		return "", "", fmt.Errorf("failed to get remote.origin.url: %w", err)
	}
	raw := strings.TrimSpace(string(out))

	// Case: git@github.com:4okimi7uki/gh-pr-formatter.git
	if strings.HasPrefix(raw, "git@") {
		_, path, ok := strings.Cut(raw, ":")
		if !ok {
			return "", "", fmt.Errorf("unexpexted ssh remote: %q", raw)
		}
		return splitOwnerRepo(path)
	}

	// Case: https://github.com/4okimi7uki/gh-pr-formatter.git
	u, err := url.Parse(raw)
	if err == nil && u.Host != "" {
		path := strings.TrimPrefix(u.Path, "/")
		return splitOwnerRepo(path)
	}

	return splitOwnerRepo(raw)
}

func ResolveRepo(slugFlag string) (string, string, error) {
	if slugFlag != "" {
		return splitOwnerRepo(slugFlag)
	}

	if owner, name, err := repoFromEnv(); err == nil {
		return owner, name, nil
	}

	if owner, name, err := repoFromGit(); err == nil {
		return owner, name, nil
	}

	return "", "", fmt.Errorf("could not resolve repository (no --repo flag, GITHUB_REPOSITORY, or git remote.origin.url)")
}
