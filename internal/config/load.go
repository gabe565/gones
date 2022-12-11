package config

import (
	"bytes"
	"errors"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/posflag"
	"github.com/knadh/koanf/providers/rawbytes"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

func Load(cmd *cobra.Command) error {
	cfgFile, err := cmd.Flags().GetString("config")
	if err != nil {
		return err
	}
	if cfgFile == "" {
		cfgDir, err := GetDir()
		if err != nil {
			return err
		}

		cfgFile = filepath.Join(cfgDir, "config.yaml")
	}

	var cfgNotExists bool
	cfgContents, err := os.ReadFile(cfgFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			cfgNotExists = true
		} else {
			return err
		}
	}

	parser := yaml.Parser()

	if err := K.Load(rawbytes.Provider(cfgContents), parser); err != nil {
		return err
	}

	err = K.Load(posflag.ProviderWithValue(cmd.Flags(), ".", K, func(key string, value string) (string, interface{}) {
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

	newCfg, err := K.Marshal(parser)
	if err != nil {
		return err
	}

	if !bytes.Equal(cfgContents, newCfg) {
		if cfgNotExists {
			log.WithField("file", cfgFile).Info("Creating config")
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

	log.WithField("file", cfgFile).Info("Loaded config")

	return nil
}
