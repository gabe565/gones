package controller

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Keymap struct {
	Regular map[Button]ebiten.Key
	Turbo   map[Button]ebiten.Key
}

var Player1Keymap = Keymap{
	Regular: map[Button]ebiten.Key{
		Up:      ebiten.KeyW,
		Right:   ebiten.KeyD,
		Down:    ebiten.KeyS,
		Left:    ebiten.KeyA,
		Start:   ebiten.KeyEnter,
		Select:  ebiten.KeyShiftRight,
		ButtonA: ebiten.KeyM,
		ButtonB: ebiten.KeyN,
	},
	Turbo: map[Button]ebiten.Key{
		ButtonA: ebiten.KeyK,
		ButtonB: ebiten.KeyJ,
	},
}

var Player2Keymap = Keymap{
	Regular: map[Button]ebiten.Key{
		Up:      ebiten.KeyHome,
		Right:   ebiten.KeyPageDown,
		Down:    ebiten.KeyEnd,
		Left:    ebiten.KeyDelete,
		Start:   ebiten.KeyKPEnter,
		Select:  ebiten.KeyKPAdd,
		ButtonA: ebiten.KeyKP3,
		ButtonB: ebiten.KeyKP2,
	},
	Turbo: map[Button]ebiten.Key{
		ButtonA: ebiten.KeyKP6,
		ButtonB: ebiten.KeyKP5,
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
