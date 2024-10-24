//go:build !js

package config

import (
	"bytes"
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"gabe565.com/gones/internal/consts"
	"gabe565.com/gones/internal/log"
	"gabe565.com/utils/must"
	"github.com/knadh/koanf/providers/posflag"
	"github.com/knadh/koanf/providers/rawbytes"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/cobra"
)

func (conf *Config) Load(cmd *cobra.Command, name, hash string) error {
	log.Init(cmd.ErrOrStderr())

	k := koanf.New(".")

	// Load default config
	if err := k.Load(structs.Provider(conf, "toml"), nil); err != nil {
		return err
	}

	// Find config file
	var gameCfgFile string
	cfgFile := must.Must2(cmd.Flags().GetString("config"))
	if cfgFile == "" {
		cfgDir, err := GetDir()
		if err != nil {
			return err
		}

		cfgFile = filepath.Join(cfgDir, "config.toml")
		gameCfgFile = filepath.Join(cfgDir, "games", hash+".toml")
	}

	if err := conf.loadMainConfig(k, cfgFile); err != nil {
		return err
	}

	if gameCfgFile != "" {
		if err := conf.loadGameOverrides(k, gameCfgFile, name); err != nil {
			return err
		}
	}

	if err := conf.loadFlags(k, cmd); err != nil {
		return err
	}

	paletteDir, err := GetPaletteDir()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(paletteDir, 0o777); err != nil {
		return err
	}

	return err
}

func (conf *Config) loadMainConfig(k *koanf.Koanf, path string) error {
	logger := slog.With("file", path)

	var cfgNotExists bool
	// Load config file if exists
	cfgContents, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			cfgNotExists = true
		} else {
			return err
		}
	}

	// Parse config file
	if err := k.Load(rawbytes.Provider(cfgContents), TOMLParser{}); err != nil {
		return err
	}

	if err := fixConfig(k); err != nil {
		return err
	}

	if err := k.UnmarshalWithConf("", conf, koanf.UnmarshalConf{Tag: "toml"}); err != nil {
		return err
	}

	// Update config if necessary
	newCfg, err := toml.Marshal(conf)
	if err != nil {
		return err
	}

	if !bytes.Equal(cfgContents, newCfg) {
		if cfgNotExists {
			logger.Info("Creating main config")

			if err := os.MkdirAll(filepath.Dir(path), 0o777); err != nil {
				return err
			}
		} else {
			logger.Info("Updating main config")
		}

		if err := os.WriteFile(path, newCfg, 0o666); err != nil {
			return err
		}
	}

	logger.Info("Loaded main config")
	return nil
}

func (conf *Config) loadGameOverrides(k *koanf.Koanf, path, name string) error {
	logger := slog.With("file", path)

	b, err := os.ReadFile(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}

		logger.Info("Creating game config")
		if err := os.MkdirAll(filepath.Dir(path), 0o777); err != nil {
			return err
		}

		if err := os.WriteFile(path, []byte("# Overrides for "+name), 0o666); err != nil {
			return err
		}

		return nil
	}

	if err := k.Load(rawbytes.Provider(b), TOMLParser{}); err != nil {
		return err
	}

	if err := k.UnmarshalWithConf("", conf, koanf.UnmarshalConf{Tag: "toml"}); err != nil {
		return err
	}

	if err := fixConfig(k); err != nil {
		return err
	}

	logger.Info("Loaded game config")
	return nil
}

func (conf *Config) loadFlags(k *koanf.Koanf, cmd *cobra.Command) error {
	lookup := flagTable()
	if err := k.Load(posflag.ProviderWithValue(cmd.Flags(), ".", k, func(key string, value string) (string, any) {
		if k, ok := lookup[key]; ok {
			key = k
		} else {
			key = ""
		}
		return key, value
	}), nil); err != nil {
		return err
	}

	return k.UnmarshalWithConf("", conf, koanf.UnmarshalConf{Tag: "toml"})
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

	// Overscan min/max
	if val := k.Int("ui.trim.top"); val < 0 || val >= consts.Height/2 {
		slog.Warn("Invalid top trim. Setting to default.")
		if err := k.Set("ui.trim.top", NewDefault().UI.Overscan.Top); err != nil {
			return err
		}
	}
	if val := k.Int("ui.trim.right"); val < 0 || val >= consts.Width/2 {
		slog.Warn("Invalid right trim. Setting to default.")
		if err := k.Set("ui.trim.right", NewDefault().UI.Overscan.Right); err != nil {
			return err
		}
	}
	if val := k.Int("ui.trim.bottom"); val < 0 || val >= consts.Height/2 {
		slog.Warn("Invalid bottom trim. Setting to default.")
		if err := k.Set("ui.trim.bottom", NewDefault().UI.Overscan.Bottom); err != nil {
			return err
		}
	}
	if val := k.Int("ui.trim.left"); val < 0 || val >= consts.Width/2 {
		slog.Warn("Invalid left trim. Setting to default.")
		if err := k.Set("ui.trim.left", NewDefault().UI.Overscan.Left); err != nil {
			return err
		}
	}

	return nil
}
