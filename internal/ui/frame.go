package ui

import (
	"fmt"
	"math"
	"os"
	"strings"
)

func Boxed(lines ...string) {
	maxLen := 0
	for _, line := range lines {
		if len(line) > maxLen {
			maxLen = len(line)
		}
	}

	paddingLeft := 1
	paddingRight := 1
	innerWidth := paddingLeft + maxLen + paddingRight
	horizWidth := innerWidth

	w := float64(horizWidth)
	r30 := int(math.Round(w * 0.30))
	r25 := int(math.Round(w * 0.25))
	r15 := horizWidth - r30 - r30 - r25

	top := Pink("╭") +
		strings.Repeat(Pink("─"), r30) +
		strings.Repeat(LightPurple("─"), r30) +
		strings.Repeat(Purple("─"), r25) +
		strings.Repeat(Blue("─"), r15) +
		Blue("╮")

	bottom := Pink("╰") +
		strings.Repeat(Pink("─"), r15) +
		strings.Repeat(LightPurple("─"), r25) +
		strings.Repeat(Purple("─"), r30) +
		strings.Repeat(Blue("─"), r30) +
		Blue("╯")

	_, _ = fmt.Fprintln(os.Stdout, top)

	for _, l := range lines {
		padRightExtra := maxLen - len(l)

		_, _ = fmt.Fprintf(os.Stdout, "%s%s%s%s%s\n",
			Pink("│"),
			strings.Repeat(" ", paddingLeft),
			l,
			strings.Repeat(" ", paddingRight+padRightExtra),
			Blue("│"),
		)
	}

	_, _ = fmt.Fprintln(os.Stdout, bottom)
}

// func SuccessBox(lines ...string) {
// 	const WIDTH = 60
// 	bar := strings.Repeat(Green("═"), WIDTH)

// 	fmt.Fprintln(os.Stderr, bar)
// 	for _, l := range lines {
// 		fmt.Fprintf(os.Stderr, " %-56s \n", l)
// 	}
// 	fmt.Fprintln(os.Stderr, bar)
// }
