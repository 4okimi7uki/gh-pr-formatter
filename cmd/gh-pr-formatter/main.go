package main

import (
	"fmt"
	"os"

	"github.com/4okimi7uki/gh-pr-formatter/internal/gh"
	"github.com/4okimi7uki/gh-pr-formatter/internal/pr"
	"github.com/4okimi7uki/gh-pr-formatter/internal/template"
)

func main() {
	if err := gh.CheckEnvironment(); err != nil {
		gh.PrintHelp(err)
		os.Exit(1)
	}

	mergedPr, err := gh.GetPrList()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	groupedPrs := pr.GroupedPrsByAuthor(mergedPr)
	fileName, err := template.BuildMarkdown(groupedPrs)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println("========================================================")
	fmt.Println(" 🎉 SUCCESS: Release PR Markdown created successfully!")
	fmt.Println(" -> Output: " + fileName)
	fmt.Println("========================================================")

}
