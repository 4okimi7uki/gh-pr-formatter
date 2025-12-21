package auth

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/zalando/go-keyring"
	"golang.org/x/term"
)

const (
	Service     = "gh-pr-formatter"
	User        = "default"
	EnvTokenKey = "GH_PR_FORMATTER_TOKEN"
)

var ErrTokenNotFound = errors.New("token not found")

func SaveToken(token string) error {
	token = strings.TrimSpace(token)
	if token == "" {
		return errors.New("token is empty")
	}
	return keyring.Set(Service, User, token)
}

func LoadGitHubToken() (string, error) {
	// Env
	_ = godotenv.Load()
	if v := strings.TrimSpace(os.Getenv(EnvTokenKey)); v != "" {
		return v, nil
	}

	v, err := keyring.Get(Service, User)
	if err != nil {
		// go-keyring は OS ごとにエラー文言が違うので、文字列でも吸収する
		msg := strings.ToLower(err.Error())
		if strings.Contains(msg, "no found") ||
			strings.Contains(msg, "no such") ||
			strings.Contains(msg, "could not be found") {
			return "", ErrTokenNotFound
		}

		return "", fmt.Errorf("failed to read token from keyring: %w", err)
	}
	v = strings.TrimSpace(v)
	if v == "" {
		return "", ErrTokenNotFound
	}

	return v, nil
}

func DeleteGitHubToken() error {
	err := keyring.Delete(Service, User)

	if err != nil {
		// 既に無い場合は成功扱いにしてUXを良くする
		msg := strings.ToLower(err.Error())
		if strings.Contains(msg, "not found") ||
			strings.Contains(msg, "no such") ||
			strings.Contains(msg, "could not be found") {
			return nil
		}
		return fmt.Errorf("failed to delete token from keyring: %w", err)
	}
	return nil
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
