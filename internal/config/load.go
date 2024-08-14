//go:build !js

package config

import (
	"bytes"
	"errors"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/knadh/koanf/providers/posflag"
	"github.com/knadh/koanf/providers/rawbytes"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
	"github.com/lmittmann/tint"
	"github.com/mattn/go-isatty"
	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/cobra"
)

func Load(cmd *cobra.Command) (*Config, error) {
	InitLog(cmd.ErrOrStderr())

	k := koanf.New(".")
	conf := NewDefault()

	// Load default config
	if err := k.Load(structs.Provider(conf, "toml"), nil); err != nil {
		return nil, err
	}

	// Find config file
	cfgFile, err := cmd.Flags().GetString("config")
	if err != nil {
		return nil, err
	}
	if cfgFile == "" {
		cfgDir, err := GetDir()
		if err != nil {
			return nil, err
		}

		cfgFile = filepath.Join(cfgDir, "config.toml")
	}
	logger := slog.With("file", cfgFile)

	var cfgNotExists bool
	// Load config file if exists
	cfgContents, err := os.ReadFile(cfgFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			cfgNotExists = true
		} else {
			return nil, err
		}
	}

	// Parse config file
	parser := TOMLParser{}
	if err := k.Load(rawbytes.Provider(cfgContents), parser); err != nil {
		return nil, err
	}

	if err := fixConfig(k); err != nil {
		return nil, err
	}

	if err := k.UnmarshalWithConf("", &conf, koanf.UnmarshalConf{Tag: "toml"}); err != nil {
		return nil, err
	}

	newCfg, err := toml.Marshal(conf)
	if err != nil {
		return nil, err
	}

	if !bytes.Equal(cfgContents, newCfg) {
		if cfgNotExists {
			logger.Info("Creating config")

			if err := os.MkdirAll(filepath.Dir(cfgFile), 0o777); err != nil {
				return nil, err
			}
		} else {
			logger.Info("Updating config")
		}

		if err := os.WriteFile(cfgFile, newCfg, 0o666); err != nil {
			return nil, err
		}
	}

	// Load flags
	flagTable := flagTable()
	err = k.Load(posflag.ProviderWithValue(cmd.Flags(), ".", k, func(key string, value string) (string, any) {
		if k, ok := flagTable[key]; ok {
			key = k
		} else {
			key = ""
		}
		return key, value
	}), nil)
	if err != nil {
		return nil, err
	}

	if err := k.UnmarshalWithConf("", &conf, koanf.UnmarshalConf{Tag: "toml"}); err != nil {
		return nil, err
	}

	paletteDir, err := GetPaletteDir()
	if err != nil {
		return nil, err
	}

	if err := os.MkdirAll(paletteDir, 0o777); err != nil && !errors.Is(err, os.ErrExist) {
		return nil, err
	}

	logger.Info("Loaded config")
	return &conf, err
}

func fixConfig(k *koanf.Koanf) error {
	// Migrate `input.keys` to `input`
	if k.Exists("input.keys") {
		inputKeys := k.Get("input.keys").(map[string]any)
		if err := k.Set("input", inputKeys); err != nil {
			return err
		}
		k.Delete("input.keys")
	}

	// Turbo duty cycle min
	if val := k.Int("input.turbo_duty_cycle"); val < 2 {
		slog.Warn("Turbo duty cycle must be 2 or greater. Setting value to 2.")
		if err := k.Set("input.turbo_duty_cycle", 2); err != nil {
			return err
		}
	}

	// Autosave interval min
	if val := k.Duration("state.autosave_interval"); val < 10*time.Second {
		slog.Warn("Autosave interval must be 10s or greater. Setting value to 10s.")
		if err := k.Set("state.interval", 10*time.Second); err != nil {
			return err
		}
	}

	// Volume min/max
	if val := k.Float64("audio.volume"); val < 0 {
		slog.Warn("Minimum volume is 0. Setting to 0.")
		if err := k.Set("audio.volume", 0); err != nil {
			return err
		}
	} else if val > 1 {
		slog.Warn("Maximum volume is 1. Setting to 1.")
		if err := k.Set("audio.volume", 1); err != nil {
			return err
		}
	}

	return nil
}

func InitLog(out io.Writer) {
	var color bool
	if f, ok := out.(*os.File); ok {
		color = isatty.IsTerminal(f.Fd()) || isatty.IsCygwinTerminal(f.Fd())
	}

	slog.SetDefault(slog.New(
		tint.NewHandler(out, &tint.Options{
			Level:      slog.LevelInfo,
			TimeFormat: time.Kitchen,
			NoColor:    !color,
		}),
	))
}
