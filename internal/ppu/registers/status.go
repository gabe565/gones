package registers

import "github.com/gabe565/gones/internal/bitflags"

const (
	SpriteOverflow bitflags.Flags = 1 << (iota + 5)
	SpriteZeroHit
	VblankStarted
)
