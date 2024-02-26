package config

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

func NewDefault() Config {
	return Config{
		UI: UI{
			Fullscreen:     false,
			Scale:          3,
			PauseUnfocused: true,
		},
		Audio: Audio{
			Enabled: true,
		},
		State: State{
			Resume:   true,
			Interval: Duration(time.Minute),
		},
		Input: Input{
			Reset:      ebiten.KeyR,
			ResetHold:  Duration(500 * time.Millisecond),
			State1Save: ebiten.KeyF1,
			State1Load: ebiten.KeyF5,

			FastForward:     ebiten.KeyF,
			FastForwardRate: 3,
			Fullscreen:      ebiten.KeyF11,

			TurboDutyCycle: 4,

			Player1: Keymap{
				A:      ebiten.KeyM,
				B:      ebiten.KeyN,
				Start:  ebiten.KeyEnter,
				Select: ebiten.KeyShiftRight,
				Up:     ebiten.KeyW,
				Down:   ebiten.KeyS,
				Left:   ebiten.KeyA,
				Right:  ebiten.KeyD,
				ATurbo: ebiten.KeyK,
				BTurbo: ebiten.KeyJ,
			},

			Player2: Keymap{
				A:      ebiten.KeyKP3,
				B:      ebiten.KeyKP2,
				Start:  ebiten.KeyKPEnter,
				Select: ebiten.KeyKPAdd,
				Up:     ebiten.KeyHome,
				Down:   ebiten.KeyEnd,
				Left:   ebiten.KeyDelete,
				Right:  ebiten.KeyPageDown,
				ATurbo: ebiten.KeyKP6,
				BTurbo: ebiten.KeyKP5,
			},
		},
	}
}