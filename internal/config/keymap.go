package config

import (
	"github.com/gabe565/gones/internal/controller/button"
	"github.com/hajimehoshi/ebiten/v2"
)

type Keymap struct {
	Up     ebiten.Key `toml:"up"`
	Left   ebiten.Key `toml:"left"`
	Down   ebiten.Key `toml:"down"`
	Right  ebiten.Key `toml:"right"`
	Start  ebiten.Key `toml:"start"`
	Select ebiten.Key `toml:"select"`
	A      ebiten.Key `toml:"a"`
	B      ebiten.Key `toml:"b"`
	ATurbo ebiten.Key `toml:"a_turbo" comment:"Key to press the A button repeatedly (must be held)."`
	BTurbo ebiten.Key `toml:"b_turbo" comment:"Key to press the B button repeatedly (must be held)."`
}

func (k Keymap) GetMap() map[button.Button]ebiten.Key {
	return map[button.Button]ebiten.Key{
		button.Up:     k.Up,
		button.Left:   k.Left,
		button.Down:   k.Down,
		button.Right:  k.Right,
		button.Start:  k.Start,
		button.Select: k.Select,
		button.A:      k.A,
		button.B:      k.B,
	}
}

func (k Keymap) GetTurboMap() map[button.Button]ebiten.Key {
	return map[button.Button]ebiten.Key{
		button.A: k.ATurbo,
		button.B: k.BTurbo,
	}
}
