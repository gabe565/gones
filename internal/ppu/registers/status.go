package registers

import "github.com/gabe565/gones/internal/bitflags"

const (
	_ bitflags.Flags = 1 << iota
	_
	_
	_
	_
	SpriteOverflow
	SpriteZeroHit
	VblankStarted
)
