package config

import (
	"github.com/spf13/cobra"
)

func Flags(cmd *cobra.Command) {
	cmd.Flags().StringP("config", "c", "", "Config file (default is $HOME/.config/gones/config.yaml)")
	if err := cmd.RegisterFlagCompletionFunc("config", func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
		return []string{"toml"}, cobra.ShellCompDirectiveFilterFileExt
	}); err != nil {
		panic(err)
	}

	cmd.Flags().Bool("debug", false, "Start with step debugging enabled")
	cmd.Flags().Bool("trace", false, "Enable trace logging")
	cmd.Flags().Float64("scale", 3, "Default UI scale")
	cmd.Flags().BoolP("fullscreen", "f", false, "Start in fullscreen")
	cmd.Flags().BoolP("audio", "a", true, "Enabled audio output")
	cmd.Flags().Bool("resume", true, "Automatically resume where you left off")
	cmd.Flags().String("palette", "", "Optional palette (.pal) file to use")
	if err := cmd.RegisterFlagCompletionFunc("palette", func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
		return []string{"pal"}, cobra.ShellCompDirectiveFilterFileExt
	}); err != nil {
		panic(err)
	}
	cmd.Flags().Bool("pause-unfocused", true, "Pauses when the window loses focus. Optional, but audio will be glitchy when the game is running in the background.")
}

func flagTable() map[string]string {
	return map[string]string{
		"debug":           "debug.enabled",
		"trace":           "debug.trace",
		"scale":           "ui.scale",
		"fullscreen":      "ui.fullscreen",
		"audio":           "audio.enabled",
		"resume":          "state.resume",
		"palette":         "ui.palette",
		"pause-unfocused": "ui.pause_unfocused",
	}
}
