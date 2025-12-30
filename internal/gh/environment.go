package gh

import (
	"errors"
	"os/exec"
	"strings"
)

var ErrNotGitRepo = errors.New("not a git repository")

func CheckEnvironment(repo string) error {
	if repo == "" {
		// git repo check
		out, err := exec.Command("git", "rev-parse", "--is-inside-work-tree").CombinedOutput()
		if err != nil || strings.TrimSpace(string(out)) != "true" {
			return ErrNotGitRepo
		}
	}

	return nil
}
