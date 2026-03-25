package cmd

import (
	"fmt"
	"os"

	"github.com/4okimi7uki/gh-pr-formatter/internal/auth"
	"github.com/4okimi7uki/gh-pr-formatter/internal/gh"
	"github.com/4okimi7uki/gh-pr-formatter/internal/ui"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var authLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Save a GitHub token securely (Keychain/credential store)",
	RunE: func(cmd *cobra.Command, args []string) error {
		// check already logged in
		v, _ := auth.LoadGitHubToken()
		if v != "" {
			result, _ := gh.ValidateGitHubToken(v)
			if result.Valid {
				fmt.Println("Already logged in (token available)")
				return nil
			} else {
				auth.DeleteGitHubToken()
			}
		}

		fmt.Fprintln(os.Stderr, "Enter your Github token: ")

		b, err := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Fprintln(os.Stderr)
		if err != nil {
			return fmt.Errorf("failed to read token: %w", err)
		}

		if err := auth.SaveToken(string(b)); err != nil {
			return err
		}

		fmt.Printf("%s Token saved successfully.\n", ui.Green("✔"))
		fmt.Printf("  - Run `gh-pr-formatter auth status` to verify your login.\n")
		fmt.Printf("  - You can also set %s for CI.\n", auth.EnvTokenKey)
		return nil
	},
}

func init() {
	authCmd.AddCommand(authLoginCmd)
}
