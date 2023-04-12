package config

import (
	"os"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/knadh/koanf/v2"
)

var K = koanf.New(".")

var Path string

var configDir = "gones"

func defaultConfig() map[string]any {
	return map[string]any{
		"audio": map[string]any{
			"enabled": true,
		},
		"debug": map[string]any{
			"enabled": false,
			"trace":   false,
		},
		"state": map[string]any{
			"resume": true,
		},
		"ui": map[string]any{
			"fullscreen": false,
			"scale":      3,
		},
		"input": map[string]any{
			"keys": map[string]any{
				"reset":       ebiten.KeyR,
				"state1_save": ebiten.KeyF1,
				"state1_load": ebiten.KeyF5,

				"fast_forward": ebiten.KeyF,
				"fullscreen":   ebiten.KeyF11,

				"player1": map[string]any{
					"up":      ebiten.KeyW,
					"left":    ebiten.KeyA,
					"down":    ebiten.KeyS,
					"right":   ebiten.KeyD,
					"start":   ebiten.KeyEnter,
					"select":  ebiten.KeyShiftRight,
					"a":       ebiten.KeyM,
					"b":       ebiten.KeyN,
					"a_turbo": ebiten.KeyK,
					"b_turbo": ebiten.KeyJ,
				},
				"player2": map[string]any{
					"up":      ebiten.KeyHome,
					"left":    ebiten.KeyDelete,
					"down":    ebiten.KeyEnd,
					"right":   ebiten.KeyPageDown,
					"start":   ebiten.KeyKPEnter,
					"select":  ebiten.KeyKPAdd,
					"a":       ebiten.KeyKP3,
					"b":       ebiten.KeyKP2,
					"a_turbo": ebiten.KeyKP6,
					"b_turbo": ebiten.KeyKP5,
				},
			},
		},
	}
}

func GetDir() (string, error) {
	if xdgConfigDir := os.Getenv("XDG_CONFIG_DIR"); xdgConfigDir != "" {
		return filepath.Join(xdgConfigDir, configDir), nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, ".config", configDir), nil
}

func GetStatesDir() (string, error) {
	configDir, err := GetDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(configDir, "states"), nil
}

func GetSramDir() (string, error) {
	configDir, err := GetDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(configDir, "sav"), nil
}
