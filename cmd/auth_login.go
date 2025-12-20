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
		fmt.Fprintln(os.Stderr, "Github token (input hidden): ")

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

func Init() {
	authCmd.AddCommand(authLoginCmd)
}
