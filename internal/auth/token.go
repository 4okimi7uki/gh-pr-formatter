package auth

import (
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/zalando/go-keyring"
	"golang.org/x/term"
)

const (
	service = "gh-pr-formatter"
	user    = "default"
)

func LoadGitHubToken() (string, error) {
	// #1 Keychain
	token, err := keyring.Get(service, user)
	if err == nil && token != "" {
		return token, nil
	}

	// #2-1 Env
	_ = godotenv.Load()
	token = os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return "", fmt.Errorf("GITHUB_TOKEN not found (keychain or env)")
	}

	// #2-2 Save to keychain
	if err := keyring.Set(service, user, token); err != nil {
		return "", fmt.Errorf("failed to save token to keychain: %w", err)
	}

	return token, nil
}

func DeleteGitHubToken() error {
	return keyring.Delete(service, user)
}

func GetTokenFromPrompt() (string, error) {
	fmt.Fprint(os.Stderr, "Enter your GitHub Token: ")
	b, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Fprintln(os.Stderr)
	if err != nil {
		return "", err
	}

	token := strings.TrimSpace(string(b))
	if token == "" {
		return "", fmt.Errorf("empty token")
	}
	return token, nil
}
