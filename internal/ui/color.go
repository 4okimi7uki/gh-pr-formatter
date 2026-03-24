package ui

import "github.com/fatih/color"

var (
	Red         = color.New(color.FgRed).SprintfFunc()
	Green       = color.RGB(67, 219, 88).SprintfFunc()
	Yellow      = color.RGB(255, 219, 76).SprintfFunc()
	Mastered    = color.RGB(208, 175, 76).SprintfFunc()
	Pink        = color.RGB(255, 76, 180).SprintfFunc()
	LightPurple = color.RGB(231, 76, 255).SprintfFunc()
	Purple      = color.RGB(151, 76, 255).SprintfFunc()
	Blue        = color.RGB(103, 73, 255).SprintfFunc()
	LightOrange = color.RGB(255, 133, 81).SprintfFunc()
	LimeYellow  = color.RGB(202, 234, 119).SprintfFunc()
	Lime        = color.RGB(37, 198, 168).SprintfFunc()
)
