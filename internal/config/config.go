package config

import (
	"github.com/knadh/koanf/v2"
	"os"
	"path/filepath"
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
