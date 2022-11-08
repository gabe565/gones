package registers

import "github.com/gabe565/gones/internal/bitflags"

const (
	Grayscale bitflags.Flags = 1 << iota
	Leftmost8pxlBackground
	Leftmost8pxlSprite
	ShowBackground
	ShowSprites
	EmphasizeRed
	EmphasizeGreen
	EmphasizeBlue
)
