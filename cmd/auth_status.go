package cmd

import (
	"fmt"

	"github.com/4okimi7uki/gh-pr-formatter/internal/auth"
	"github.com/4okimi7uki/gh-pr-formatter/internal/gh"
	"github.com/4okimi7uki/gh-pr-formatter/internal/ui"
	"github.com/spf13/cobra"
)

var authStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show authentication status",
	RunE: func(cmd *cobra.Command, args []string) error {
		token, err := auth.LoadGitHubToken()
		if err == nil {
			result, err := gh.ValidateGitHubToken(token)
			if err != nil {
				return err
			}

			if !result.Valid {
				fmt.Println("\033[1mgithub.com\033[0m")
				fmt.Printf("　%s Authentication failed\n", ui.Red("✗"))
				fmt.Printf("  - Detail: %s\n", result.Detail)
				fmt.Printf("  - Tip: Your token may be invalid or expired. Please check your token.\n")
				return nil
			}
			fmt.Println("\033[1mgithub.com\033[0m")
			fmt.Printf("　%s Logged in as %s\n", ui.Green("✔"), result.Login)
			fmt.Printf("　 - Token status: valid\n")
			fmt.Printf("　 - Token scopes: %s\n", result.TokenScopes)
			return nil
		}
		if err == auth.ErrTokenNotFound {
			fmt.Println("\033[1mgithub.com\033[0m")
			fmt.Printf("  %s Logged out\n", "○")
			fmt.Printf("  - Tip: Run `gh-pr-formatter auth login`\n")
			return nil
		}
		return err
	},
}

func init() {
	authCmd.AddCommand(authStatusCmd)
}
