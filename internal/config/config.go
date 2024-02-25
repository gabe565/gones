package config

import (
	"os"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
)

type Config struct {
	Audio Audio `toml:"audio"`
	Debug Debug `toml:"debug"`
	State State `toml:"state"`
	UI    UI    `toml:"ui"`
	Input Input `toml:"input"`
}

type Audio struct {
	Enabled bool `toml:"enabled"`
}

type Debug struct {
	Enabled bool `toml:"enabled"`
	Trace   bool `toml:"trace"`
}

type State struct {
	Resume   bool     `toml:"resume"`
	Interval Duration `toml:"interval"`
}

type UI struct {
	Fullscreen     bool    `toml:"fullscreen"`
	Scale          float64 `toml:"scale"`
	PauseUnfocused bool    `toml:"pause_unfocused"`
}

type Input struct {
	Keys Keys `toml:"keys"`
}

type Keys struct {
	Reset       ebiten.Key `toml:"reset"`
	State1Save  ebiten.Key `toml:"state1_save"`
	State1Load  ebiten.Key `toml:"state1_load"`
	FastForward ebiten.Key `toml:"fast_forward"`
	Fullscreen  ebiten.Key `toml:"fullscreen"`
	Player1     Keymap     `toml:"player1"`
	Player2     Keymap     `toml:"player2"`
}

var configDir = "gones"

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
