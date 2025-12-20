package ui

import (
	"fmt"
	"os"
	"strings"

	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/termenv"
)

func gradientLine(s string, leftHex, rightHex string) string {
	profile := termenv.ColorProfile()
	c1, _ := colorful.Hex(leftHex)
	c2, _ := colorful.Hex(rightHex)

	runes := []rune(s)
	n := len(runes)
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

	// fmt.Fprint(os.Stdout, Pink.Sprint(logo)+"\n")
	for t := range strings.SplitSeq(logo, "\n") {
		fmt.Fprint(os.Stdout, gradientLine(string(t)+"\n", "#F4307F", "#F9F339"))
	}
	fmt.Fprint(os.Stdout, LightOrange.Sprint(" make your release work a little bit easier :P\n"))
	fmt.Println()
}
