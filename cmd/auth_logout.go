package cmd

import (
	"fmt"

	"github.com/4okimi7uki/gh-pr-formatter/internal/auth"
	"github.com/spf13/cobra"
)

var authLogoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Delete the saved GitHub token",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := auth.DeleteGitHubToken(); err != nil {
			return err
		}
		fmt.Println("Token deleted.")
		return nil
	},
}

func init() {
	authCmd.AddCommand(authLogoutCmd)
}
