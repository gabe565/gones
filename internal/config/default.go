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
		State: State{
			Resume:           true,
			AutosaveInterval: Duration(time.Minute),
			UndoStateCount:   5,
		},
		Input: Input{
			Reset:             Key(ebiten.KeyR),
			ResetHold:         Duration(500 * time.Millisecond),
			State1Save:        Key(ebiten.KeyF1),
			State1Load:        Key(ebiten.KeyF5),
			StateUndoModifier: Key(ebiten.KeyShiftLeft),

			FastForward:     Key(ebiten.KeyF),
			FastForwardRate: 3,
			Fullscreen:      Key(ebiten.KeyF11),

			Screenshot: Key(ebiten.KeyBackslash),

			TurboDutyCycle: 4,

			Player1: Keymap{
				A:      Key(ebiten.KeyM),
				B:      Key(ebiten.KeyN),
				Start:  Key(ebiten.KeyEnter),
				Select: Key(ebiten.KeyShiftRight),
				Up:     Key(ebiten.KeyW),
				Down:   Key(ebiten.KeyS),
				Left:   Key(ebiten.KeyA),
				Right:  Key(ebiten.KeyD),
				ATurbo: Key(ebiten.KeyK),
				BTurbo: Key(ebiten.KeyJ),
			},

			Player2: Keymap{
				A:      Key(ebiten.KeyKP3),
				B:      Key(ebiten.KeyKP2),
				Start:  Key(ebiten.KeyKPEnter),
				Select: Key(ebiten.KeyKPAdd),
				Up:     Key(ebiten.KeyHome),
				Down:   Key(ebiten.KeyEnd),
				Left:   Key(ebiten.KeyDelete),
				Right:  Key(ebiten.KeyPageDown),
				ATurbo: Key(ebiten.KeyKP6),
				BTurbo: Key(ebiten.KeyKP5),
			},
		},
		Audio: Audio{
			Enabled: true,
			Volume:  1,
			Channels: AudioChannels{
				Triangle: true,
				Square1:  true,
				Square2:  true,
				Noise:    true,
				PCM:      true,
			},
		},
	}
}
