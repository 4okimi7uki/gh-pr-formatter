package cmd

import (
	"fmt"

	"github.com/4okimi7uki/gh-pr-formatter/internal/auth"
	"github.com/spf13/cobra"
)

var authStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show Authentication status",
	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := auth.LoadGitHubToken()
		if err == nil {
			fmt.Println("Auth status: Logged in (token found)")
			return nil
		}
		if err == auth.ErrTokenNotFound {
			fmt.Println("Auth status: Logged out")
			fmt.Println("Tip: Run `gh-pr-formatter auth login`")
			return nil
		}
		return err
	},
}

func init() {
	authCmd.AddCommand(authStatusCmd)
}
