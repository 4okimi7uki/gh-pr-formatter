package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/4okimi7uki/gh-pr-formatter/internal/auth"
	"github.com/4okimi7uki/gh-pr-formatter/internal/client"
	"github.com/4okimi7uki/gh-pr-formatter/internal/gh"
	"github.com/4okimi7uki/gh-pr-formatter/internal/pkg/repository"
	"github.com/4okimi7uki/gh-pr-formatter/internal/pr"
	"github.com/4okimi7uki/gh-pr-formatter/internal/template"
	"github.com/4okimi7uki/gh-pr-formatter/internal/ui"
	"github.com/briandowns/spinner"
)

var Version = "v0.0.0-dev"

type StringSlice []string

func (s *StringSlice) String() string {
	return strings.Join(*s, ",")
}

func (s *StringSlice) Set(v string) error {
	*s = append(*s, v)
	return nil
}

func main() {
	ui.DisplayLogo()
	token, err := auth.LoadGitHubToken()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for _, arg := range os.Args[1:] {
		if strings.HasPrefix(arg, "-") && !strings.HasPrefix(arg, "--") {
			fmt.Fprintf(os.Stderr,
				"Error: '%s' is not a valid option. Use '--' prefix for all flags.\n", arg)
			os.Exit(1)
		}
	}
	// get options
	versionFlag := flag.Bool("version", false, "Print version information")
	repo := flag.String("repo", "", "Repository to operate on (e.g. owner/repo)")
	var excludePrefixes StringSlice
	flag.Var(&excludePrefixes, "x", "exclude branch prefix (repeatable)")
	flag.Parse()

	// version check
	if *versionFlag {
		fmt.Printf("gh-pr-formatter version %s\n", Version)
		os.Exit(0)
	}

	// loader start
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Suffix = " PR fetching..."
	s.Writer = os.Stderr
	s.Start()
	defer s.Stop()

	owner, repoName, err := repository.ResolveRepo(*repo)
	if err != nil {
		fmt.Fprintf(os.Stderr,
			"Error: '%s'\n", err)
		os.Exit(1)
	}

	c := client.NewClient(token)
	prs, from, err := c.ListMergedPullRequests(owner, repoName, 50, excludePrefixes...)
	if err != nil {
		fmt.Fprintf(os.Stderr,
			"Error: '%s'\n", err)
		os.Exit(1)
	}

	if len(prs) == 0 {
		s.Stop()
		fmt.Fprintln(os.Stdout)
		fmt.Fprintln(os.Stdout, "No merged pull requests found on 'develop' since last 'main' merge.")
		os.Exit(0)
	}

	if err := gh.CheckEnvironment(*repo); err != nil {
		s.Stop()
		fmt.Fprintln(os.Stderr)
		gh.PrintHelp(err)
		os.Exit(1)
	}

	groupedPrs := pr.GroupedPrsByAuthor(prs)
	fileName, err := template.BuildMarkdown(groupedPrs, from)
	if err != nil {
		s.Stop()
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// display result
	s.Stop()
	fmt.Fprintln(os.Stdout)
	ui.SuccessBox("🎉 SUCCESS:", "  Release PR Markdown created successfully!", "  ↳ Output: "+fileName)

	// version check
	if msg, err := gh.CheckLatestVersion("4okimi7uki", "gh-pr-formatter", Version); err == nil && msg != "" {
		ui.Boxed(msg, "Download: https://github.com/4okimi7uki/gh-pr-formatter/releases")
	} else if err != nil {
		_ = err // or log.Printf("version check failed: %v", err)
	}

}
