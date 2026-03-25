package cmd

import (
	"fmt"

	"github.com/4okimi7uki/gh-pr-formatter/internal/auth"
	"github.com/4okimi7uki/gh-pr-formatter/internal/ui"
	"github.com/spf13/cobra"
)

var authLogoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Delete the saved GitHub token",
	RunE: func(cmd *cobra.Command, args []string) error {
		// check already logged in
		v, _ := auth.LoadGitHubToken()
		if v == "" {
			fmt.Println("Already logged out — nothing to do.")
			return nil
		}

		if err := auth.DeleteGitHubToken(); err != nil {
			return err
		}
		fmt.Printf("%s Logged out.\n", ui.Mastered("○"))
		fmt.Printf("  - Run `gh-pr-formatter auth login` to sign in again.\n")
		return nil
	},
}

func init() {
	authCmd.AddCommand(authLogoutCmd)
}
