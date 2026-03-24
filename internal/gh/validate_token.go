package gh

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type GitHubUser struct {
	Login string `json:"login"`
	ID    int64  `json:"id"`
}

type TokenValidationResult struct {
	Valid       bool
	Login       string
	Detail      string
	TokenScopes string
}

func summarizeScopes(scopes []string) string {
	priority := []string{"gist", "read:org", "repo", "workflow"}

	var picked []string
	used := map[string]bool{}

	for _, p := range priority {
		for _, s := range scopes {
			trimS := strings.TrimSpace(s)
			if trimS == p && !used[trimS] {
				picked = append(picked, trimS)
				used[trimS] = true
			}
		}
	}

	restCount := 0
	for _, s := range scopes {
		trimS := strings.TrimSpace(s)
		if !used[trimS] {
			restCount++
		}
	}

	if restCount > 0 {
		if len(picked) == 0 {
			return fmt.Sprintf("(%d scopes)", restCount)
		}
		return fmt.Sprintf("%s (+%d more)", strings.Join(picked, ", "), restCount)
	}

	return strings.Join(picked, ", ")
}

func ValidateGitHubToken(token string) (TokenValidationResult, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		return TokenValidationResult{Detail: "failed to create request"}, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return TokenValidationResult{Detail: "failed to send request"}, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	scopes := resp.Header.Get("X-OAuth-Scopes")
	var scopeSummary string
	if scopes != "" {
		scopeSummary = summarizeScopes(strings.Split(scopes, ","))
	}

	switch resp.StatusCode {
	case http.StatusOK:
		var user GitHubUser
		if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
			return TokenValidationResult{Detail: "failed to decode response"}, err
		}
		return TokenValidationResult{Valid: true, Login: user.Login, TokenScopes: scopeSummary}, nil
	case http.StatusUnauthorized:
		return TokenValidationResult{Detail: "invalid or expired token"}, nil
	default:
		return TokenValidationResult{Detail: fmt.Sprintf("unexpected response status: %d", resp.StatusCode)}, nil
	}
}
