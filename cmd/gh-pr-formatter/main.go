package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/4okimi7uki/gh-pr-formatter/internal/gh"
	"github.com/4okimi7uki/gh-pr-formatter/internal/template"
	"github.com/4okimi7uki/gh-pr-formatter/internal/utils"
	"github.com/briandowns/spinner"
	"github.com/fatih/color"
)

func main() {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Suffix = " PR fetching..."
	s.Writer = os.Stderr
	s.Start()
	defer s.Stop()

	for _, arg := range os.Args[1:] {
		if strings.HasPrefix(arg, "-") && !strings.HasPrefix(arg, "--") {
			fmt.Fprintf(os.Stderr,
				"Error: '%s' is not a valid option. Use '--' prefix for all flags.\n", arg)
			os.Exit(1)
		}
	}

	repo := flag.String("repo", "", "Repository to operate on (e.g. owner/repo)")
	flag.Parse()

	if err := gh.CheckEnvironment(); err != nil {
		s.Stop()
		fmt.Fprintln(os.Stderr)
		gh.PrintHelp(err)
		os.Exit(1)
	}

	mergedPr, err := gh.GetPrList(*repo)
	if err != nil {
		s.Stop()
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	groupedPrs := utils.GroupedPrsByAuthor(mergedPr)
	fileName, err := template.BuildMarkdown(groupedPrs)
	if err != nil {
		s.Stop()
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	s.Stop()
	fmt.Fprintln(os.Stderr)

	successColor := color.RGB(67, 219, 88)

	// mergedColor := color.RGB(137, 87, 226)

	utils.PrintlnNoErr(successColor, "==========================================================")
	utils.PrintlnNoErr(successColor, " 🎉 SUCCESS: Release PR Markdown created successfully!")
	fmt.Println(successColor.Sprint(" -> Output: ") + fileName)
	utils.PrintlnNoErr(successColor, "==========================================================")

}
