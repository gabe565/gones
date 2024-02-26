package config

import (
	"github.com/gabe565/gones/internal/controller/button"
	"github.com/hajimehoshi/ebiten/v2"
)

type Keymap struct {
	A      ebiten.Key `toml:"a"`
	B      ebiten.Key `toml:"b"`
	Start  ebiten.Key `toml:"start"`
	Select ebiten.Key `toml:"select"`
	Up     ebiten.Key `toml:"up"`
	Down   ebiten.Key `toml:"down"`
	Left   ebiten.Key `toml:"left"`
	Right  ebiten.Key `toml:"right"`

	ATurbo ebiten.Key `toml:"a_turbo" comment:"Key to press the A button repeatedly (must be held)."`
	BTurbo ebiten.Key `toml:"b_turbo" comment:"Key to press the B button repeatedly (must be held)."`
}

func (k Keymap) GetMap() map[button.Button]ebiten.Key {
	return map[button.Button]ebiten.Key{
		button.A:      k.A,
		button.B:      k.B,
		button.Start:  k.Start,
		button.Select: k.Select,
		button.Up:     k.Up,
		button.Down:   k.Down,
		button.Left:   k.Left,
		button.Right:  k.Right,
	}
}

func (k Keymap) GetTurboMap() map[button.Button]ebiten.Key {
	return map[button.Button]ebiten.Key{
		button.A: k.ATurbo,
		button.B: k.BTurbo,
	}
}
