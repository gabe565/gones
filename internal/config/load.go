//go:build !wasm

package config

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"

	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/posflag"
	"github.com/knadh/koanf/providers/rawbytes"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func Load(cmd *cobra.Command) error {
	// Load default config
	if err := K.Load(confmap.Provider(defaultConfig(), ""), nil); err != nil {
		return err
	}

	// Find config file
	cfgFile, err := cmd.Flags().GetString("config")
	if err != nil {
		return err
	}
	if cfgFile == "" {
		cfgDir, err := GetDir()
		if err != nil {
			return err
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
			return err
		}
	}

	// Parse config file
	parser := TOMLParser{}
	if err := K.Load(rawbytes.Provider(cfgContents), parser); err != nil {
		return err
	}

	newCfg, err := K.Marshal(parser)
	if err != nil {
		return err
	}

	if !bytes.Equal(cfgContents, newCfg) {
		if cfgNotExists {
			log.WithField("file", cfgFile).Info("Creating config")

			if err := os.MkdirAll(filepath.Dir(cfgFile), 0o777); err != nil {
				return err
			}
		} else {
			log.WithField("file", cfgFile).Info("Updating config")
		}

		f, err := os.Create(cfgFile)
		if err != nil {
			return err
		}
		defer func(f *os.File) {
			_ = f.Close()
		}(f)

		if _, err := f.Write(newCfg); err != nil {
			return err
		}

		if err := f.Close(); err != nil {
			return err
		}
	}

	// Load flags
	err = K.Load(posflag.ProviderWithValue(cmd.Flags(), ".", K, func(key string, value string) (string, any) {
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
		return err
	}

	log.WithField("file", cfgFile).Info("Loaded config")

	return nil
}
