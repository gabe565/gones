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
			Keys: Keys{
				Reset:      ebiten.KeyR,
				ResetHold:  Duration(500 * time.Millisecond),
				State1Save: ebiten.KeyF1,
				State1Load: ebiten.KeyF5,

				FastForward:     ebiten.KeyF,
				FastForwardRate: 3,
				Fullscreen:      ebiten.KeyF11,

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
