package config

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

func NewDefault() Config {
	return Config{
		Audio: Audio{
			Enabled: true,
		},
		Debug: Debug{
			Enabled: false,
			Trace:   false,
		},
		State: State{
			Resume:   true,
			Interval: Duration(time.Minute),
		},
		UI: UI{
			Fullscreen:     false,
			Scale:          3,
			PauseUnfocused: true,
		},
		Input: Input{
			Keys: Keys{
				Reset:      ebiten.KeyR,
				State1Save: ebiten.KeyF1,
				State1Load: ebiten.KeyF5,

				FastForward: ebiten.KeyF,
				Fullscreen:  ebiten.KeyF11,

				Player1: Keymap{
					Up:     ebiten.KeyW,
					Left:   ebiten.KeyA,
					Down:   ebiten.KeyS,
					Right:  ebiten.KeyD,
					Start:  ebiten.KeyEnter,
					Select: ebiten.KeyShiftRight,
					A:      ebiten.KeyM,
					B:      ebiten.KeyN,
					ATurbo: ebiten.KeyK,
					BTurbo: ebiten.KeyJ,
				},
				Player2: Keymap{
					Up:     ebiten.KeyHome,
					Left:   ebiten.KeyDelete,
					Down:   ebiten.KeyEnd,
					Right:  ebiten.KeyPageDown,
					Start:  ebiten.KeyKPEnter,
					Select: ebiten.KeyKPAdd,
					A:      ebiten.KeyKP3,
					B:      ebiten.KeyKP2,
					ATurbo: ebiten.KeyKP6,
					BTurbo: ebiten.KeyKP5,
				},
			},
		},
	}
}
