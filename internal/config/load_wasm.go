package config

import (
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/posflag"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func Load(cmd *cobra.Command) error {
	if err := K.Load(confmap.Provider(defaultConfig(), ""), nil); err != nil {
		return err
	}

	err := K.Load(posflag.ProviderWithValue(cmd.Flags(), ".", K, func(key string, value string) (string, interface{}) {
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

	if err := K.Set("audio.enabled", false); err != nil {
		return err
	}

	log.Info("Loaded config")
	return nil
}
