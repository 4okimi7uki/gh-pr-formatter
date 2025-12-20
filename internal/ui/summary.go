package ui

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/4okimi7uki/gh-pr-formatter/internal/models"
)

type Target struct {
	Repository    string
	ExcludePrefix []string
	From          time.Time
}

type Summary struct {
	Target
	PRCount     int
	AuthorCount int
	OutputPath  string
	PRList      string
}

const (
	labelWidth  = 14
	noWidth     = 4
	authorWidth = 17
	branchWidth = 31
)

func PrintSummary(w io.Writer, s Summary, loc *time.Location) {
	period := fmt.Sprintf("%s → Now", s.From.In(loc).Format("2006-01-02 15:04 MST"))

	_, _ = fmt.Fprintln(w, Green.Sprint("\n\n✔ ")+"Release PR markdown generated")
	fmt.Println()

	_, _ = fmt.Fprintln(w, strings.Repeat("─", labelWidth*5))
	_, _ = fmt.Fprintf(w, " %-*s : %s\n", labelWidth, "Repository", s.Repository)
	_, _ = fmt.Fprintf(w, " %-*s : %s\n", labelWidth, "Range", period)
	_, _ = fmt.Fprintf(w, " %-*s : %s (%s)\n", labelWidth, "Exclude Prefix", s.ExcludePrefix, "applied to main")
	_, _ = fmt.Fprintf(w, " %-*s : %d (authors: %d)\n", labelWidth, "PRs", s.PRCount, s.AuthorCount)
	_, _ = fmt.Fprintf(w, " %-*s : %s\n", labelWidth, "Output", s.OutputPath)
	_, _ = fmt.Fprintln(w, strings.Repeat("─", labelWidth*5))
	_, _ = fmt.Println("\n" + s.PRList)
}

type ListSummary struct {
	Target
	Prs []models.MergedPrList
}

func PrintListSummary(w io.Writer, l ListSummary, loc *time.Location) {
	period := fmt.Sprintf("%s → Now", l.From.In(loc).Format("2006-01-02 15:04 MST"))

	fmt.Println()
	_, _ = fmt.Fprintln(w, strings.Repeat("─", labelWidth*5))
	_, _ = fmt.Fprintf(w, " %-*s : %s\n", labelWidth, "Repository", l.Repository)
	_, _ = fmt.Fprintf(w, " %-*s : %s\n", labelWidth, "Range", period)
	_, _ = fmt.Fprintf(w, " %-*s : %s (%s)\n", labelWidth, "Exclude Prefix", l.ExcludePrefix, "applied to main")
	_, _ = fmt.Fprintln(w, strings.Repeat("─", labelWidth*5))

	fmt.Println("\nMerged PRs into develop List:")

	fmt.Printf(
		"\n%-*s  %-*s %-*s %s\n",
		noWidth, "NO",
		authorWidth, "AUTHOR",
		branchWidth, "BRANCH",
		"TITLE",
	)

	for _, pr := range l.Prs {
		fmt.Printf("#%-4d @%-16s %-30s %s\n",
			pr.Number,
			pr.Author.Login,
			pr.HeadRefName,
			pr.Title,
		)
	}
	fmt.Println("---")
}
