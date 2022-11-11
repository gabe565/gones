package controller

import (
	"github.com/faiface/pixel/pixelgl"
	"github.com/gabe565/gones/internal/bitflags"
)

var Keymap = map[pixelgl.Button]bitflags.Flags{
	pixelgl.KeyUp:         Up,
	pixelgl.KeyRight:      Right,
	pixelgl.KeyDown:       Down,
	pixelgl.KeyLeft:       Left,
	pixelgl.KeyEnter:      Start,
	pixelgl.KeyRightShift: Select,
	pixelgl.KeyA:          ButtonA,
	pixelgl.KeyS:          ButtonB,
}
