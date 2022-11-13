package controller

import (
	"github.com/faiface/pixel/pixelgl"
	"github.com/gabe565/gones/internal/bitflags"
)

var Player1Keymap = map[pixelgl.Button]bitflags.Flags{
	pixelgl.KeyW:          Up,
	pixelgl.KeyD:          Right,
	pixelgl.KeyS:          Down,
	pixelgl.KeyA:          Left,
	pixelgl.KeyEnter:      Start,
	pixelgl.KeyRightShift: Select,
	pixelgl.KeyN:          ButtonA,
	pixelgl.KeyM:          ButtonB,
}
