package ui

import (
	"github.com/fatih/color"
)

func PrintlnNoErr(c *color.Color, s string) {
	_, _ = c.Println(s)
}
