package cmd

import (
	"fmt"
	"os"

	"github.com/4okimi7uki/gh-pr-formatter/internal/auth"
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
			fmt.Println("Already logged in (token available)")
			return nil
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

		fmt.Println("Token saved.")
		fmt.Printf("Tip: you can also set %s for CI.\n", auth.EnvTokenKey)
		return nil
	},
}

func init() {
	authCmd.AddCommand(authLoginCmd)
}
