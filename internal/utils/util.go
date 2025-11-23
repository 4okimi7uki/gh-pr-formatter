package utils

import (
	model "github.com/4okimi7uki/gh-pr-formatter/internal/model"
	"github.com/fatih/color"
)

func GroupedPrsByAuthor(prs []model.MergedPrList) map[string][]int {
	m := make(map[string][]int)
	for _, pr := range prs {
		m[pr.Author.Login] = append(m[pr.Author.Login], pr.Number)
	}
	return m
}

func PrintlnNoErr(c *color.Color, s string) {
	_, _ = c.Println(s)
}
