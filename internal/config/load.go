//go:build !wasm

package config

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"

	"github.com/knadh/koanf/providers/posflag"
	"github.com/knadh/koanf/providers/rawbytes"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
	"github.com/pelletier/go-toml/v2"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func Load(cmd *cobra.Command) (*Config, error) {
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

	if err := k.UnmarshalWithConf("", &conf, koanf.UnmarshalConf{Tag: "toml"}); err != nil {
		return nil, err
	}

	newCfg, err := toml.Marshal(conf)
	if err != nil {
		return nil, err
	}

	if !bytes.Equal(cfgContents, newCfg) {
		if cfgNotExists {
			log.WithField("file", cfgFile).Info("Creating config")

			if err := os.MkdirAll(filepath.Dir(cfgFile), 0o777); err != nil {
				return nil, err
			}
		} else {
			log.WithField("file", cfgFile).Info("Updating config")
		}

		if err := os.WriteFile(cfgFile, newCfg, 0o666); err != nil {
			return nil, err
		}
	}

	// Load flags
	err = k.Load(posflag.ProviderWithValue(cmd.Flags(), ".", k, func(key string, value string) (string, any) {
		if k, ok := flagConfigTable[key]; ok {
			key = k
		}
		for _, name := range excludeFromConfig {
			if key == name {
				return "", value
			}
		}
		return key, value
	}), nil)
	if err != nil {
		return nil, err
	}

	if err := k.UnmarshalWithConf("", &conf, koanf.UnmarshalConf{Tag: "toml"}); err != nil {
		return nil, err
	}

	log.WithField("file", cfgFile).Info("Loaded config")
	return &conf, err
}
