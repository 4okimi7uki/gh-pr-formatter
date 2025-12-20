package ui

import (
	"github.com/fatih/color"
)

func PrintlnNoErr(c *color.Color, s string) {
	_, _ = c.Println(s)
}

func DisplayLogo() {
	const logo = `        __                              ___                           __   __
 .-----|  |--.______.-----.----.______.'  _.-----.----.--------.---.-|  |_|  |_.-----.----.
 |  _  |     |______|  _  |   _|______|   _|  _  |   _|        |  _  |   _|   _|  -__|   _|
 |___  |__|__|      |   __|__|        |__| |_____|__| |__|__|__|___._|____|____|_____|__|
 |_____|            |__|
`

	PrintlnNoErr(Pink, logo)
	PrintlnNoErr(Pink, " make your release work a littele bit easier :P\n")
}
