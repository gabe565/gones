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

var Player2Keymap = map[pixelgl.Button]bitflags.Flags{
	pixelgl.KeyHome:     Up,
	pixelgl.KeyPageDown: Right,
	pixelgl.KeyEnd:      Down,
	pixelgl.KeyDelete:   Left,
	pixelgl.KeyKPEnter:  Start,
	pixelgl.KeyKPAdd:    Select,
	pixelgl.KeyKP2:      ButtonA,
	pixelgl.KeyKP3:      ButtonB,
}

const (
	ToggleTrace = pixelgl.KeyTab
	ToggleDebug = pixelgl.KeyGraveAccent
	StepFrame   = pixelgl.Key1
	RunToRender = pixelgl.Key2
)
