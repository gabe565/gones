package controller

import (
	"github.com/gabe565/gones/internal/bitflags"
	"github.com/hajimehoshi/ebiten/v2"
)

type Keymap map[ebiten.Key]bitflags.Flags

var Player1Keymap = Keymap{
	ebiten.KeyW:          Up,
	ebiten.KeyD:          Right,
	ebiten.KeyS:          Down,
	ebiten.KeyA:          Left,
	ebiten.KeyEnter:      Start,
	ebiten.KeyShiftRight: Select,
	ebiten.KeyN:          ButtonA,
	ebiten.KeyM:          ButtonB,
}

var Player2Keymap = Keymap{
	ebiten.KeyHome:     Up,
	ebiten.KeyPageDown: Right,
	ebiten.KeyEnd:      Down,
	ebiten.KeyDelete:   Left,
	ebiten.KeyKPEnter:  Start,
	ebiten.KeyKPAdd:    Select,
	ebiten.KeyKP2:      ButtonA,
	ebiten.KeyKP3:      ButtonB,
}

var Joystick = map[bitflags.Flags]ebiten.StandardGamepadButton{
	Up:      ebiten.StandardGamepadButtonLeftTop,
	Right:   ebiten.StandardGamepadButtonLeftRight,
	Down:    ebiten.StandardGamepadButtonLeftBottom,
	Left:    ebiten.StandardGamepadButtonLeftLeft,
	Start:   ebiten.StandardGamepadButtonCenterRight,
	Select:  ebiten.StandardGamepadButtonCenterLeft,
	ButtonA: ebiten.StandardGamepadButtonRightRight,
	ButtonB: ebiten.StandardGamepadButtonRightBottom,
}

const (
	Reset = ebiten.KeyR

	FastForward      = ebiten.KeyF
	ToggleFullscreen = ebiten.KeyF11

	ToggleTrace = ebiten.KeyTab
	ToggleDebug = ebiten.KeyGraveAccent
	StepFrame   = ebiten.Key1
	RunToRender = ebiten.Key2
)
