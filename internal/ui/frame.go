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

	top := Pink.Sprint("╭") +
		strings.Repeat(Pink.Sprint("─"), r30) +
		strings.Repeat(LightPurple.Sprint("─"), r30) +
		strings.Repeat(Purple.Sprint("─"), r25) +
		strings.Repeat(Blue.Sprint("─"), r15) +
		Blue.Sprint("╮")

	bottom := Pink.Sprint("╰") +
		strings.Repeat(Pink.Sprint("─"), r15) +
		strings.Repeat(LightPurple.Sprint("─"), r25) +
		strings.Repeat(Purple.Sprint("─"), r30) +
		strings.Repeat(Blue.Sprint("─"), r30) +
		Blue.Sprint("╯")

	fmt.Fprintln(os.Stdout, top)

	for _, l := range lines {
		padRightExtra := maxLen - len(l)

		fmt.Fprintf(os.Stdout, "%s%s%s%s%s\n",
			Pink.Sprint("│"),
			strings.Repeat(" ", paddingLeft),
			l,
			strings.Repeat(" ", paddingRight+padRightExtra),
			Blue.Sprint("│"),
		)
	}

	fmt.Fprintln(os.Stdout, bottom)
}

func SuccessBox(lines ...string) {
	const WIDTH = 60
	bar := strings.Repeat(Green.Sprint("═"), WIDTH)

	fmt.Fprintln(os.Stderr, bar)
	for _, l := range lines {
		fmt.Fprintf(os.Stderr, " %-56s \n", l)
	}
	fmt.Fprintln(os.Stderr, bar)
}
