package config

import (
	"os"
	"path/filepath"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/knadh/koanf/v2"
)

var K = koanf.New(".")

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
			"resume":   true,
			"interval": time.Minute,
		},
		"ui": map[string]any{
			"fullscreen": false,
			"scale":      3,
		},
		"input": map[string]any{
			"keys": map[string]any{
				"reset":       ebiten.KeyR.String(),
				"state1_save": ebiten.KeyF1.String(),
				"state1_load": ebiten.KeyF5.String(),

				"fast_forward": ebiten.KeyF.String(),
				"fullscreen":   ebiten.KeyF11.String(),

				"player1": map[string]any{
					"up":      ebiten.KeyW.String(),
					"left":    ebiten.KeyA.String(),
					"down":    ebiten.KeyS.String(),
					"right":   ebiten.KeyD.String(),
					"start":   ebiten.KeyEnter.String(),
					"select":  ebiten.KeyShiftRight.String(),
					"a":       ebiten.KeyM.String(),
					"b":       ebiten.KeyN.String(),
					"a_turbo": ebiten.KeyK.String(),
					"b_turbo": ebiten.KeyJ.String(),
				},
				"player2": map[string]any{
					"up":      ebiten.KeyHome.String(),
					"left":    ebiten.KeyDelete.String(),
					"down":    ebiten.KeyEnd.String(),
					"right":   ebiten.KeyPageDown.String(),
					"start":   ebiten.KeyKPEnter.String(),
					"select":  ebiten.KeyKPAdd.String(),
					"a":       ebiten.KeyKP3.String(),
					"b":       ebiten.KeyKP2.String(),
					"a_turbo": ebiten.KeyKP6.String(),
					"b_turbo": ebiten.KeyKP5.String(),
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
