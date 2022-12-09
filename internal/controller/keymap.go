package controller

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Keymap[T ebiten.Key | ebiten.StandardGamepadButton] struct {
	Regular map[T]uint8
	Turbo   map[T]uint8
}

var Player1Keymap = Keymap[ebiten.Key]{
	Regular: map[ebiten.Key]uint8{
		ebiten.KeyW:          Up,
		ebiten.KeyD:          Right,
		ebiten.KeyS:          Down,
		ebiten.KeyA:          Left,
		ebiten.KeyEnter:      Start,
		ebiten.KeyShiftRight: Select,
		ebiten.KeyN:          ButtonA,
		ebiten.KeyM:          ButtonB,
	},
	Turbo: map[ebiten.Key]uint8{
		ebiten.KeyJ: ButtonA,
		ebiten.KeyK: ButtonB,
	},
}

var Player2Keymap = Keymap[ebiten.Key]{
	Regular: map[ebiten.Key]uint8{
		ebiten.KeyHome:     Up,
		ebiten.KeyPageDown: Right,
		ebiten.KeyEnd:      Down,
		ebiten.KeyDelete:   Left,
		ebiten.KeyKPEnter:  Start,
		ebiten.KeyKPAdd:    Select,
		ebiten.KeyKP2:      ButtonA,
		ebiten.KeyKP3:      ButtonB,
	},
	Turbo: map[ebiten.Key]uint8{
		ebiten.KeyKP5: ButtonA,
		ebiten.KeyKP6: ButtonB,
	},
}

var Joystick = Keymap[ebiten.StandardGamepadButton]{
	Regular: map[ebiten.StandardGamepadButton]uint8{
		ebiten.StandardGamepadButtonLeftTop:     Up,
		ebiten.StandardGamepadButtonLeftRight:   Right,
		ebiten.StandardGamepadButtonLeftBottom:  Down,
		ebiten.StandardGamepadButtonLeftLeft:    Left,
		ebiten.StandardGamepadButtonCenterRight: Start,
		ebiten.StandardGamepadButtonCenterLeft:  Select,
		ebiten.StandardGamepadButtonRightRight:  ButtonA,
		ebiten.StandardGamepadButtonRightBottom: ButtonB,
	},
	Turbo: map[ebiten.StandardGamepadButton]uint8{
		ebiten.StandardGamepadButtonRightTop:  ButtonA,
		ebiten.StandardGamepadButtonRightLeft: ButtonB,
	},
}

const (
	Reset = ebiten.KeyR

	SaveState1 = ebiten.KeyF1
	LoadState1 = ebiten.KeyF5

	FastForward      = ebiten.KeyF
	ToggleFullscreen = ebiten.KeyF11

	ToggleTrace = ebiten.KeyTab
	ToggleDebug = ebiten.KeyGraveAccent
	StepFrame   = ebiten.Key1
	RunToRender = ebiten.Key2
)
