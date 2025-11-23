package main

import (
	"fmt"
	"os"
	"time"

	"github.com/4okimi7uki/gh-pr-formatter/internal/gh"
	"github.com/4okimi7uki/gh-pr-formatter/internal/template"
	"github.com/4okimi7uki/gh-pr-formatter/internal/utils"
	"github.com/briandowns/spinner"
	"github.com/fatih/color"
)

func main() {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Writer = os.Stderr
	s.Start()
	defer s.Stop()

	if err := gh.CheckEnvironment(); err != nil {
		s.Stop()
		fmt.Fprintln(os.Stderr)
		gh.PrintHelp(err)
		os.Exit(1)
	}

	mergedPr, err := gh.GetPrList()
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

	successColor := color.RGB(44, 242, 0)

	utils.PrintlnNoErr(successColor, "==========================================================")
	utils.PrintlnNoErr(successColor, " 🎉 SUCCESS: Release PR Markdown created successfully!")
	utils.PrintlnNoErr(successColor, " -> Output: "+fileName)
	utils.PrintlnNoErr(successColor, "==========================================================")

}
