package main

import (
	"errors"
	"github.com/faiface/pixel/pixelgl"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "gones ROM",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("No ROM provided")
			}
			cmd.SilenceUsage = true
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			pixelgl.Run(func() {
				err = run(args[0])
			})
			return err
		},
	}
	return cmd
}
