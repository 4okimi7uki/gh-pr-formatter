package pr

import "github.com/4okimi7uki/gh-pr-formatter/internal/models"

func GroupedPrsByAuthor(prs []models.MergedPrList) map[string][]int {
	m := make(map[string][]int)
	for _, pr := range prs {
		m[pr.Author.Login] = append(m[pr.Author.Login], pr.Number)
	}
	return m
}
