package controller

import (
	"github.com/gabe565/gones/internal/config"
	"github.com/gabe565/gones/internal/controller/button"
	"github.com/hajimehoshi/ebiten/v2"
)

type Keymap struct {
	Regular map[button.Button]ebiten.Key
	Turbo   map[button.Button]ebiten.Key
}

func NewKeymap(conf *config.Config, player Player) Keymap {
	var keymap config.Keymap
	switch player {
	case Player1:
		keymap = conf.Input.Keys.Player1
	case Player2:
		keymap = conf.Input.Keys.Player2
	default:
		panic("invalid player: " + player)
	}

	return Keymap{
		Regular: keymap.GetMap(),
		Turbo:   keymap.GetTurboMap(),
	}
}

func LoadKeys(conf *config.Config) {
	Reset = conf.Input.Keys.Reset
	SaveState1 = conf.Input.Keys.State1Save
	LoadState1 = conf.Input.Keys.State1Load
	FastForward = conf.Input.Keys.FastForward
	ToggleFullscreen = conf.Input.Keys.Fullscreen
}

var (
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
