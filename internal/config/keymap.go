package config

import (
	"github.com/gabe565/gones/internal/controller/button"
	"github.com/hajimehoshi/ebiten/v2"
)

type Keymap struct {
	A      Key `toml:"a"`
	B      Key `toml:"b"`
	Start  Key `toml:"start"`
	Select Key `toml:"select"`
	Up     Key `toml:"up"`
	Down   Key `toml:"down"`
	Left   Key `toml:"left"`
	Right  Key `toml:"right"`

	ATurbo Key `toml:"a_turbo" comment:"Key to press the A button repeatedly (must be held)."`
	BTurbo Key `toml:"b_turbo" comment:"Key to press the B button repeatedly (must be held)."`
}

func (k Keymap) GetMap() map[button.Button]ebiten.Key {
	return map[button.Button]ebiten.Key{
		button.A:      ebiten.Key(k.A),
		button.B:      ebiten.Key(k.B),
		button.Start:  ebiten.Key(k.Start),
		button.Select: ebiten.Key(k.Select),
		button.Up:     ebiten.Key(k.Up),
		button.Down:   ebiten.Key(k.Down),
		button.Left:   ebiten.Key(k.Left),
		button.Right:  ebiten.Key(k.Right),
	}
}

func (k Keymap) GetTurboMap() map[button.Button]ebiten.Key {
	return map[button.Button]ebiten.Key{
		button.A: ebiten.Key(k.ATurbo),
		button.B: ebiten.Key(k.BTurbo),
	}
}
