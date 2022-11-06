package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/gabe565/gones/internal/callbacks"
	"github.com/gabe565/gones/internal/console"
	"os"
)

func main() {
	if err := NewCommand().Execute(); err != nil {
		os.Exit(1)
	}
}

func run(path string, callback callbacks.CallbackHandler) error {
	cfg := pixelgl.WindowConfig{
		Title:  "GoNES",
		Bounds: pixel.R(0, 0, 10*32, 10*32),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		return err
	}

	c, err := console.New(path)
	if err != nil {
		return err
	}
	c.Reset()

	if callback != nil {
		c.Callback = callback(win)
	}

	if err := c.Run(); err != nil {
		return err
	}

	return nil
}
