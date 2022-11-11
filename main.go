package main

import (
	"context"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/gabe565/gones/internal/console"
	"github.com/gabe565/gones/internal/controller"
	"github.com/gabe565/gones/internal/ppu"
	log "github.com/sirupsen/logrus"
	"image/color"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path/filepath"
)

func main() {
	if err := NewCommand().Execute(); err != nil {
		os.Exit(1)
	}
}

type Run struct {
	Path  string
	Trace bool
	Pprof string
}

func (r Run) Run() error {
	if r.Pprof != "" {
		go func() {
			log.WithField("address", r.Pprof).Info("starting pprof")
			if err := http.ListenAndServe(r.Pprof, nil); err != nil {
				log.WithError(err).Error("failed to start pprof")
			}
		}()
	}

	cfg := pixelgl.WindowConfig{
		Title:  "GoNES",
		Bounds: pixel.R(0, 0, 3*ppu.Width, 3*ppu.TrimmedHeight),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		return err
	}

	win.SetTitle(filepath.Base(r.Path) + " | GoNES")

	console, err := console.New(r.Path)
	if err != nil {
		return err
	}

	console.CPU.EnableTrace = r.Trace

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errCh := make(chan error)
	go func() {
		err := console.CPU.Run(ctx)
		errCh <- err
	}()

	for {
		select {
		case err := <-errCh:
			return err
		case <-console.Bus.RenderStart:
			if win.Closed() {
				return nil
			}

			for button, key := range controller.Keymap {
				if win.JustPressed(button) {
					console.Bus.Controller1.Set(key, true)
				} else if win.JustReleased(button) {
					console.Bus.Controller1.Set(key, false)
				}
			}

			win.Clear(color.Black)
			pic := console.PPU.Render()
			console.Bus.RenderDone <- struct{}{}

			sprite := pixel.NewSprite(pic, pic.Bounds())
			sprite.Draw(win, pixel.IM.Moved(win.Bounds().Center()).Scaled(win.Bounds().Center(), 3))
			win.Update()
		}
	}
}
