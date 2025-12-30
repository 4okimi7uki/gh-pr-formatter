package template

import (
	_ "embed"
	"fmt"
	"sort"
	"strings"
	"text/template"
	"time"
)

type TemplateData struct {
	List string
}

//go:embed template.md
var templateSource string

func renderTemplate(data TemplateData) (string, error) {
	tpl, err := template.New("release").Parse(templateSource)
	if err != nil {
		return "", err
	}

	var out strings.Builder
	if err := tpl.Execute(&out, data); err != nil {
		return "", err
	}

	return out.String(), nil
}

func BuildMarkdown(groupedPrs map[string][]int, from time.Time) (string, string, error) {
	var prList strings.Builder

	authors := make([]string, 0, len(groupedPrs))
	for a := range groupedPrs {
		authors = append(authors, a)
	}
	sort.Strings(authors)

	for _, author := range authors {
		nums := groupedPrs[author]
		fmt.Fprintf(&prList, "@%s\n", author)
		for _, num := range nums {
			fmt.Fprintf(&prList, "- #%d\n", num)
		}
		prList.WriteString("\n")
	}

	rendered, err := renderTemplate(TemplateData{
		List: prList.String(),
	})
	if err != nil {
		return "", "", fmt.Errorf("failed to render template: %w", err)
	}

	return rendered, prList.String(), nil
}

func CountPRs(grouped map[string][]int) int {
	total := 0
	for _, prs := range grouped {
		total += len(prs)
	}
	return total
}
