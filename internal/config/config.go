package config

import (
	"image"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"gabe565.com/gones/internal/consts"
)

type Config struct {
	UI    UI    `toml:"ui"`
	State State `toml:"state"`
	Input Input `toml:"input"`
	Audio Audio `toml:"audio"`
	Debug Debug `toml:"debug,omitempty"`
}

type UI struct {
	Fullscreen        bool     `toml:"fullscreen" comment:"Default fullscreen state. Fullscreen can also be toggled with a key (F11 by default)."`
	Scale             float64  `toml:"scale" comment:"Multiplier used to scale the UI."`
	PauseUnfocused    bool     `toml:"pause_unfocused" comment:"Pauses when the window loses focus. Optional, but audio will be glitchy when the game is running in the background."`
	Palette           string   `toml:"palette" comment:"Palette (.pal) file to use. An embedded palette will be used when blank."`
	RemoveSpriteLimit bool     `toml:"remove_sprite_limit" comment:"Removes the original hardware's 8 horizontal sprite limitation. When enabled, sprites will no longer flicker."`
	Overscan          Overscan `toml:"overscan,inline" comment:"Change the number of rows/cols of overscan."`
}

type Overscan struct {
	Top    int `toml:"top"`
	Right  int `toml:"right"`
	Bottom int `toml:"bottom"`
	Left   int `toml:"left"`
}

func (t Overscan) Rect() image.Rectangle {
	return image.Rect(t.Left, t.Top, consts.Width-t.Right, consts.Height-t.Bottom)
}

type State struct {
	Resume           bool     `toml:"resume" comment:"Automatically resumes the previous game state."`
	AutosaveInterval Duration `toml:"autosave_interval" comment:"If resume is enabled, the game state will be saved regularly at the configured interval."`
	UndoStateCount   int      `toml:"undo_state_count" comment:"Number of undo states to keep in memory."`
}

type Input struct {
	Reset             Key      `toml:"reset" comment:"Key to reset the game (must be held)."`
	ResetHold         Duration `toml:"reset_hold" comment:"Time the reset button must be held."`
	State1Save        Key      `toml:"state1_save" comment:"Key to save the game state (separate from auto resume state)."`
	State1Load        Key      `toml:"state1_load" comment:"Key to load the last save state."`
	StateUndoModifier Key      `toml:"state_undo_modifier" comment:"Hold this key and press the save/load state key, and the action will be undone."`
	FastForward       Key      `toml:"fast_forward" comment:"Key to fast-forward the game (must be held)."`
	FastForwardRate   uint8    `toml:"fast_forward_rate" comment:"Fast-forward rate multiplier."`
	Fullscreen        Key      `toml:"fullscreen" comment:"Key to toggle fullscreen."`
	Screenshot        Key      `toml:"screenshot" comment:"Key to take a screenshot."`
	TurboDutyCycle    uint16   `toml:"turbo_duty_cycle" comment:"Frame duty cycle when turbo key is held (minimum: 2)."`
	Player1           Keymap   `toml:"player1" comment:"Player 1 keymap."`
	Player2           Keymap   `toml:"player2" comment:"Player 2 keymap."`
}

func (i Input) ResetHoldFrames() int {
	frames := int(time.Duration(i.ResetHold).Seconds() * 60)
	if frames == 0 {
		return 1
	}
	return frames
}

type Audio struct {
	Enabled  bool          `toml:"enabled" comment:"Enables audio output."`
	Volume   float64       `toml:"volume" comment:"Output volume (between 0 and 1)."`
	Channels AudioChannels `toml:"channels" comment:"Toggles specific audio channels."`
}

type AudioChannels struct {
	Triangle bool `toml:"triangle"`
	Square1  bool `toml:"square_1"`
	Square2  bool `toml:"square_2"`
	Noise    bool `toml:"noise"`
	PCM      bool `toml:"pcm"`
}

type Debug struct {
	Enabled bool `toml:"enabled"`
	Trace   bool `toml:"trace"`
}

const configDir = "gones"

func GetDir() (string, error) {
	switch runtime.GOOS {
	case "darwin":
		if xdgConfigHome := os.Getenv("XDG_CONFIG_HOME"); xdgConfigHome != "" {
			return filepath.Join(xdgConfigHome, configDir), nil
		}
		fallthrough
	default:
		dir, err := os.UserConfigDir()
		if err != nil {
			return "", err
		}

		dir = filepath.Join(dir, configDir)
		return dir, nil
	}
}

func GetStatesDir() (string, error) {
	configDir, err := GetDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(configDir, "states"), nil
}

func GetSRAMDir() (string, error) {
	configDir, err := GetDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(configDir, "sav"), nil
}

func GetPaletteDir() (string, error) {
	configDir, err := GetDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(configDir, "palettes"), nil
}

func GetScreenshotDir() (string, error) {
	configDir, err := GetDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(configDir, "screenshots"), nil
}
