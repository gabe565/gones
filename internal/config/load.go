package config

import (
	"errors"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/posflag"
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

	parser := yaml.Parser()

	var writeCfg bool
	if err := K.Load(file.Provider(cfgFile), parser); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			writeCfg = true
		} else {
			return err
		}
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

	if writeCfg {
		log.WithField("file", cfgFile).Info("Writing config")

		f, err := os.Create(cfgFile)
		if err != nil {
			return err
		}
		defer func(f *os.File) {
			_ = f.Close()
		}(f)

		b, err := K.Marshal(parser)
		if err != nil {
			return err
		}

		if _, err := f.Write(b); err != nil {
			return err
		}

		if err := f.Close(); err != nil {
			return err
		}
	}

	log.WithField("file", cfgFile).Info("Loaded config")

	return nil
}
