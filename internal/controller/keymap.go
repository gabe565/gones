package controller

import (
	"github.com/faiface/pixel/pixelgl"
	"github.com/gabe565/gones/internal/bitflags"
)

type Keymap map[pixelgl.Button]bitflags.Flags

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

var Joystick = map[pixelgl.GamepadButton]bitflags.Flags{
	pixelgl.ButtonDpadUp:    Up,
	pixelgl.ButtonDpadRight: Right,
	pixelgl.ButtonDpadDown:  Down,
	pixelgl.ButtonDpadLeft:  Left,
	pixelgl.ButtonStart:     Start,
	pixelgl.ButtonGuide:     Select,
	pixelgl.ButtonA:         ButtonA,
	pixelgl.ButtonB:         ButtonB,
}

const (
	Reset = pixelgl.KeyR

	FastForward = pixelgl.KeyF

	ToggleTrace = pixelgl.KeyTab
	ToggleDebug = pixelgl.KeyGraveAccent
	StepFrame   = pixelgl.Key1
	RunToRender = pixelgl.Key2
)
