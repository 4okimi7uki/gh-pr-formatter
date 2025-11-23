package template

import (
	_ "embed"
	"fmt"
	"os"
	"sort"
	"strings"
	"text/template"
	"time"

	"github.com/4okimi7uki/gh-pr-formatter/internal/utils"
	"github.com/fatih/color"
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

func BuildMarkdown(groupedPrs map[string][]int) (string, error) {
	var b strings.Builder

	authors := make([]string, 0, len(groupedPrs))
	for a := range groupedPrs {
		authors = append(authors, a)
	}
	sort.Strings(authors)

	for _, author := range authors {
		nums := groupedPrs[author]
		fmt.Fprintf(&b, "@%s\n", author)
		for _, num := range nums {
			fmt.Fprintf(&b, "- #%d\n", num)
		}
		b.WriteString("\n")
	}

	rendered, err := renderTemplate(TemplateData{
		List: b.String(),
	})
	if err != nil {
		return "", fmt.Errorf("failed to render template: %w", err)
	}

	dir := "./releasePrMarkdown"
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create dir: %w", err)
	}

	loc, _ := time.LoadLocation("Asia/Tokyo")
	now := time.Now().In(loc)
	fileName := fmt.Sprintf(`%s/release_%s.md`, dir, now.Format("20060102_1504"))

	err = os.WriteFile(fileName, []byte(rendered), 0644)
	if err != nil {
		return "", fmt.Errorf("failed to write %s: %w", fileName, err)
	}

	mergedColor := color.RGB(137, 87, 226)

	utils.PrintlnNoErr(mergedColor, "\n========================================================")
	fmt.Println(" Merged Pull Requests")
	utils.PrintlnNoErr(mergedColor, "========================================================")
	fmt.Println(b.String())

	return fileName, nil
}
