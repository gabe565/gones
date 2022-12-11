package config

import (
	"github.com/spf13/cobra"
)

var flagConfigTable = map[string]string{
	"debug":      "debug.enabled",
	"trace":      "debug.trace",
	"scale":      "ui.scale",
	"fullscreen": "ui.fullscreen",
}

var excludeFromConfig = [...]string{"config", "help", "version"}

func Flags(cmd *cobra.Command) {
	cmd.Flags().StringP("config", "c", "", "Config file (default is $HOME/.config/gones/config.yaml)")
	cmd.Flags().Bool("debug", false, "Start with step debugging enabled")
	cmd.Flags().Bool("trace", false, "Enable trace logging")
	cmd.Flags().Float64("scale", 3, "Default UI scale")
	cmd.Flags().BoolP("fullscreen", "f", false, "Start in fullscreen")
}
