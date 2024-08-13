package options

import "github.com/spf13/cobra"

type Option func(cmd *cobra.Command)

func WithVersion(version string) Option {
	return func(cmd *cobra.Command) {
		cmd.Version = buildVersion(version)
		cmd.InitDefaultVersionFlag()
	}
}
