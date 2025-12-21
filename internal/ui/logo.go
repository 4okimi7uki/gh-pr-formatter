package ui

import (
	"fmt"
	"strings"

	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/termenv"
)

func gradientLine(s string, leftHex, rightHex string, maxLen int) string {
	profile := termenv.ColorProfile()
	c1, _ := colorful.Hex(leftHex)
	c2, _ := colorful.Hex(rightHex)

	runes := []rune(s)
	n := maxLen
	if n == 0 {
		return s
	}

	var b strings.Builder
	b.Grow(len(s) * 4)

	for i, r := range runes {
		// 空白は塗らない（AAの見た目が安定しやすい）
		if r == ' ' {
			b.WriteRune(r)
			continue
		}
		t := float64(i) / float64(max(1, n-1))
		c := c1.BlendLab(c2, t)

		hex := c.Hex()
		col := termenv.String(string(r)).Foreground(profile.Color(hex))
		b.WriteString(col.String())
	}
	return b.String()
}

func PrintLogo() {
	const logo string = `        __                              ___                           __   __
 .-----|  |--.______.-----.----.______.'  _.-----.----.--------.---.-|  |_|  |_.-----.----.
 |  _  |     |______|  _  |   _|______|   _|  _  |   _|        |  _  |   _|   _|  -__|   _|
 |___  |__|__|      |   __|__|        |__| |_____|__| |__|__|__|___._|____|____|_____|__|
 |_____|            |__|
`

	const (
		startColor = "#F4307F"
		endColor   = "#F9F339"
		tagline    = "make your release work a little bit easier :P"
		repoURL    = "https://github.com/4okimi7uki/gh-pr-formatter"
	)

	logoLines := strings.Split(strings.TrimRight(logo, "\n"), "\n")
	extraLines := []string{
		"",
		" " + tagline,
		"",
		"" + repoURL + "",
		"",
	}

	all := append(append([]string{}, logoLines...), extraLines...)
	width := maxLen(all)

	// logo AA
	for _, line := range logoLines {
		fmt.Print(gradientLine(line+"\n", startColor, endColor, width))
	}

	// tagline
	fmt.Print(gradientLine("\n "+tagline+"\n", startColor, endColor, width))

	// upperBar
	periodBar := strings.Repeat(".", len(repoURL)+3)
	fmt.Println(gradientLine(periodBar, startColor, endColor, width))

	// URL
	fmt.Print(gradientLine(" "+repoURL+"\n", startColor, endColor, len(repoURL)))

	// lowerBar
	dotBar := strings.Repeat("·", len(repoURL)+3)
	fmt.Println(gradientLine(dotBar, startColor, endColor, width))
	fmt.Println()

	fmt.Println()
}

func maxLen(lines []string) int {
	max := 0
	for _, s := range lines {
		s = strings.TrimRight(s, "\r")
		if max < len(s) {
			max = len(s)
		}
	}
	return max
}
