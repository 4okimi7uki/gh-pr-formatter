package pr

func GroupedPrsByAuthor(prs []MergedPrList) map[string][]int {
	m := make(map[string][]int)
	for _, pr := range prs {
		m[pr.Author.Login] = append(m[pr.Author.Login], pr.Number)
	}
	return m
}
