package config

import (
	"os"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
)

type Config struct {
	Audio Audio `toml:"audio"`
	Debug Debug `toml:"debug,omitempty"`
	State State `toml:"state"`
	UI    UI    `toml:"ui"`
	Input Input `toml:"input"`
}

type Audio struct {
	Enabled bool `toml:"enabled" comment:"Enables audio output."`
}

type Debug struct {
	Enabled bool `toml:"enabled"`
	Trace   bool `toml:"trace"`
}

type State struct {
	Resume   bool     `toml:"resume" comment:"Automatically resumes the previous game state."`
	Interval Duration `toml:"interval" comment:"Time interval to save the game state."`
}

type UI struct {
	Fullscreen     bool    `toml:"fullscreen" comment:"Default fullscreen state. Fullscreen can also be toggled with a key (F11 by default)."`
	Scale          float64 `toml:"scale" comment:"Multiplier used to scale the UI."`
	PauseUnfocused bool    `toml:"pause_unfocused" comment:"Pauses when the window loses focus. Optional, but audio will be glitchy when the game is running in the background."`
}

type Input struct {
	Keys Keys `toml:"keys" comment:"Global keys."`
}

type Keys struct {
	Reset       ebiten.Key `toml:"reset" comment:"Key to reset the game (must be held)."`
	State1Save  ebiten.Key `toml:"state1_save" comment:"Key to save the game state (separate from auto resume state)."`
	State1Load  ebiten.Key `toml:"state1_load" comment:"Key to load the last save state."`
	FastForward ebiten.Key `toml:"fast_forward" comment:"Key to fast-forward the game (must be held)."`
	Fullscreen  ebiten.Key `toml:"fullscreen" comment:"Key to toggle fullscreen."`
	Player1     Keymap     `toml:"player1" comment:"Player 1 keymap."`
	Player2     Keymap     `toml:"player2" comment:"Player 2 keymap."`
}

var configDir = "gones"

func GetDir() (string, error) {
	if xdgConfigHome := os.Getenv("XDG_CONFIG_HOME"); xdgConfigHome != "" {
		return filepath.Join(xdgConfigHome, configDir), nil
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
