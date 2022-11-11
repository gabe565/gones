package main

import (
	"errors"
	"github.com/faiface/pixel/pixelgl"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	var action Run

	cmd := &cobra.Command{
		Use: "gones ROM",
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
	return cmd
}
