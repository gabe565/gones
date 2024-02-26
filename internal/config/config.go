package config

import (
	"os"
	"path/filepath"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type Config struct {
	UI    UI    `toml:"ui"`
	Audio Audio `toml:"audio"`
	State State `toml:"state"`
	Input Input `toml:"input"`
	Debug Debug `toml:"debug,omitempty"`
}

type UI struct {
	Fullscreen     bool    `toml:"fullscreen" comment:"Default fullscreen state. Fullscreen can also be toggled with a key (F11 by default)."`
	Scale          float64 `toml:"scale" comment:"Multiplier used to scale the UI."`
	PauseUnfocused bool    `toml:"pause_unfocused" comment:"Pauses when the window loses focus. Optional, but audio will be glitchy when the game is running in the background."`
}

type Audio struct {
	Enabled bool `toml:"enabled" comment:"Enables audio output."`
}

type State struct {
	Resume   bool     `toml:"resume" comment:"Automatically resumes the previous game state."`
	Interval Duration `toml:"interval" comment:"Time interval to save the game state."`
}

type Input struct {
	Keys Keys `toml:"keys" comment:"Global keys."`
}

type Keys struct {
	Reset           ebiten.Key `toml:"reset" comment:"Key to reset the game (must be held)."`
	ResetHold       Duration   `toml:"reset_hold" comment:"Time the reset button must be held."`
	State1Save      ebiten.Key `toml:"state1_save" comment:"Key to save the game state (separate from auto resume state)."`
	State1Load      ebiten.Key `toml:"state1_load" comment:"Key to load the last save state."`
	FastForward     ebiten.Key `toml:"fast_forward" comment:"Key to fast-forward the game (must be held)."`
	FastForwardRate uint8      `toml:"fast_forward_rate" comment:"Fast-forward rate multiplier."`
	Fullscreen      ebiten.Key `toml:"fullscreen" comment:"Key to toggle fullscreen."`
	Player1         Keymap     `toml:"player1" comment:"Player 1 keymap."`
	Player2         Keymap     `toml:"player2" comment:"Player 2 keymap."`
}

func (k Keys) ResetHoldFrames() int {
	frames := int(time.Duration(k.ResetHold).Seconds() * 60)
	if frames == 0 {
		return 1
	}
	return frames
}

type Debug struct {
	Enabled bool `toml:"enabled"`
	Trace   bool `toml:"trace"`
}

var configDir = "gones"

func GetDir() (string, error) {
	if xdgConfigHome := os.Getenv("XDG_CONFIG_HOME"); xdgConfigHome != "" {
		return filepath.Join(xdgConfigHome, configDir), nil
	}

	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	dir = filepath.Join(dir, configDir)
	return dir, nil
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
