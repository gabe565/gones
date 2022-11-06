package main

import (
	"errors"
	"github.com/faiface/pixel/pixelgl"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "gones ROM",
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error

			if len(args) < 1 {
				return errors.New("No ROM provided")
			}

			pixelgl.Run(func() {
				err = run(args[0])
			})
			return err
		},
	}
	return cmd
}
