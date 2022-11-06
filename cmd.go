package main

import (
	"errors"
	"fmt"
	"github.com/faiface/pixel/pixelgl"
	"github.com/gabe565/gones/internal/callbacks"
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

			callbackName, err := cmd.Flags().GetString("callback")
			if err != nil {
				return err
			}
			var callback callbacks.CallbackHandler
			if callbackName != "" {
				var ok bool
				callback, ok = callbacks.Callbacks[callbackName]
				if !ok {
					return fmt.Errorf("unknown callback: %s", callbackName)
				}
			}

			pixelgl.Run(func() {
				err = run(args[0], callback)
			})
			return err
		},
	}
	cmd.Flags().StringP("callback", "c", "", "Enable a callback function")
	return cmd
}
