package util

import "github.com/spf13/cobra"

func CompleteROM(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
	return []string{"nes"}, cobra.ShellCompDirectiveFilterFileExt
}
