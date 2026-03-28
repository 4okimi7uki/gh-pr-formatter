package gh

import (
	"errors"
	"fmt"
	"net/http"
	"time"
)

var (
	ErrUnauthorized = errors.New("authentication failed (401 unauthorized), please check your GitHub token")
	ErrForbidden    = errors.New("access forbidden (403 forbidden), please check your GitHub token permissions")
	ErrRepoNotFound = errors.New("repository not found (404 not found), please check the owner/repository name")
)

func ValidateRepository(token string, owner string, repo string) error {
	endpoint := fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, repo)
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	switch resp.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusNotFound:
		return ErrRepoNotFound
	case http.StatusUnauthorized:
		return ErrUnauthorized
	case http.StatusForbidden:
		return ErrForbidden
	default:
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
}
