//go:build !wasm

package config

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/fatih/color"
	"github.com/knadh/koanf/providers/posflag"
	"github.com/knadh/koanf/providers/rawbytes"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
	"github.com/pelletier/go-toml/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func Load(cmd *cobra.Command) (*Config, error) {
	InitLog()

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
			log.Info().Str("file", cfgFile).Msg("Creating config")

			if err := os.MkdirAll(filepath.Dir(cfgFile), 0o777); err != nil {
				return nil, err
			}
		} else {
			log.Info().Str("file", cfgFile).Msg("Updating config")
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

	log.Info().Str("file", cfgFile).Msg("Loaded config")
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
		log.Warn().Msg("Turbo duty cycle must be 2 or greater. Setting value to 2.")
		if err := k.Set("input.turbo_duty_cycle", 2); err != nil {
			return err
		}
	}

	// Autosave interval min
	if val := k.Duration("state.autosave_interval"); val < 10*time.Second {
		log.Warn().Msg("Autosave interval must be 10s or greater. Setting value to 10s.")
		if err := k.Set("state.interval", 10*time.Second); err != nil {
			return err
		}
	}

	// Volume min/max
	if val := k.Float64("audio.volume"); val < 0 {
		log.Warn().Msg("Minimum volume is 0. Setting to 0.")
		if err := k.Set("audio.volume", 0); err != nil {
			return err
		}
	} else if val > 1 {
		log.Warn().Msg("Maximum volume is 1. Setting to 1.")
		if err := k.Set("audio.volume", 1); err != nil {
			return err
		}
	}

	return nil
}

func InitLog() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	sprintf := fmt.Sprintf
	if !color.NoColor {
		sprintf = color.New(color.Bold).Sprintf
	}
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:     os.Stderr,
		NoColor: color.NoColor,
		FormatMessage: func(i interface{}) string {
			return sprintf("%-25s", i)
		},
	})
}
