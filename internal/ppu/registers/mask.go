package registers

import "github.com/gabe565/gones/internal/bitflags"

const (
	Grayscale bitflags.Flags = 1 << iota
	BgLeftColEnable
	SpriteLeftColEnable
	BackgroundEnable
	SpriteEnable
	EmphasizeRed
	EmphasizeGreen
	EmphasizeBlue
)
