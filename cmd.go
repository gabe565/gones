package main

import (
	"errors"
	"github.com/faiface/pixel/pixelgl"
	"github.com/gabe565/gones/internal/pprof"
	"github.com/spf13/cobra"
)

var (
	Version = "next"
	Commit  = ""
)

func NewCommand() *cobra.Command {
	var action Run

	cmd := &cobra.Command{
		Use:     "gones ROM",
		Version: buildVersion(),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("No ROM provided")
			}
			action.Path = args[0]
			cmd.SilenceUsage = true
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			pixelgl.Run(func() {
				err = action.Run()
			})
			return err
		},
	}

	cmd.Flags().BoolVar(&action.Trace, "trace", false, "Enable trace logging")
	pprof.Flag(cmd)

	return cmd
}

func buildVersion() string {
	result := Version
	if Commit != "" {
		result += " (" + Commit + ")"
	}
	return result
}
