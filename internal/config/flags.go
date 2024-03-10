package config

import (
	"github.com/spf13/cobra"
)

var flagConfigTable = map[string]string{
	"debug":           "debug.enabled",
	"trace":           "debug.trace",
	"scale":           "ui.scale",
	"fullscreen":      "ui.fullscreen",
	"audio":           "audio.enabled",
	"resume":          "state.resume",
	"palette":         "ui.palette",
	"pause-unfocused": "ui.pause_unfocused",
}

func Flags(cmd *cobra.Command) {
	cmd.Flags().StringP("config", "c", "", "Config file (default is $HOME/.config/gones/config.yaml)")
	_ = cmd.RegisterFlagCompletionFunc("config", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"toml"}, cobra.ShellCompDirectiveFilterFileExt
	})

	cmd.Flags().Bool("debug", false, "Start with step debugging enabled")
	cmd.Flags().Bool("trace", false, "Enable trace logging")
	cmd.Flags().Float64("scale", 3, "Default UI scale")
	cmd.Flags().BoolP("fullscreen", "f", false, "Start in fullscreen")
	cmd.Flags().BoolP("audio", "a", true, "Enabled audio output")
	cmd.Flags().Bool("resume", true, "Automatically resume where you left off")
	cmd.Flags().String("palette", "", "Optional palette (.pal) file to use")
	if err := cmd.RegisterFlagCompletionFunc("palette", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"pal"}, cobra.ShellCompDirectiveFilterFileExt
	}); err != nil {
		panic(err)
	}
	cmd.Flags().Bool("pause-unfocused", true, "Pauses when the window loses focus. Optional, but audio will be glitchy when the game is running in the background.")
}
