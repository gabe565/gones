package controller

import (
	"gabe565.com/gones/internal/config"
	"gabe565.com/gones/internal/controller/button"
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
		keymap = conf.Input.Player1
	case Player2:
		keymap = conf.Input.Player2
	default:
		panic("invalid player: " + player)
	}

	return Keymap{
		Regular: keymap.GetMap(),
		Turbo:   keymap.GetTurboMap(),
	}
}

const (
	ToggleTrace = ebiten.KeyTab
	ToggleDebug = ebiten.KeyGraveAccent
	StepFrame   = ebiten.Key1
	RunToRender = ebiten.Key2
)
