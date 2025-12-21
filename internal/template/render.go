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

	// dir := "./releasePrMarkdown"
	// if err := os.MkdirAll(dir, 0755); err != nil {
	// 	return "", fmt.Errorf("failed to create dir: %w", err)
	// }

	// loc, _ := time.LoadLocation("Asia/Tokyo")
	// now := time.Now().In(loc)
	// fileName := fmt.Sprintf(`%s/release_%s.md`, dir, now.Format("20060102_1504"))

	// err = os.WriteFile(fileName, []byte(rendered), 0644)
	// if err != nil {
	// 	return "", fmt.Errorf("failed to write %s: %w", fileName, err)
	// }

	// mergedColor := color.RGB(137, 87, 226)
	// headTitle := " Merged Pull Requests \n"
	// period := fmt.Sprintf("%s → Now", from.In(loc).Format("2006-01-02 15:04 MST"))

	// fmt.Fprintln(os.Stdout, ui.Green.Sprint("\n\n✔ ")+"Release PR markdown generated")
	// fmt.Println()
	// const labelWidth = 12

	// fmt.Fprintf(os.Stdout, " %-*s : %s\n", labelWidth, "Repository", "")
	// fmt.Fprintf(os.Stdout, " %-*s : %s\n", labelWidth, "Range", period)
	// fmt.Fprintf(os.Stdout, " %-*s : %d (authors: %d)\n", labelWidth, "PRs", countPRs(groupedPrs), len(groupedPrs))
	// fmt.Fprintf(os.Stdout, " %-*s : %s\n", labelWidth, "Output", fileName)
	// fmt.Fprintln(os.Stdout, strings.Repeat("─", 65))
	// fmt.Println("\n" + prList.String())

	// barLength := len(max(headTitle, period))
	// border := mergedColor.Sprint(strings.Repeat("-", barLength))
	// fmt.Println("\n\n" + border)
	// fmt.Print(headTitle)
	// fmt.Print(period)
	// fmt.Println(border)

	// fmt.Println(prList.String())
	// fmt.Println(border)

	// return fileName, nil
}

func CountPRs(grouped map[string][]int) int {
	total := 0
	for _, prs := range grouped {
		total += len(prs)
	}
	return total
}
