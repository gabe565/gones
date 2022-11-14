package main

import (
	"context"
	"errors"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/gabe565/gones/internal/console"
	"github.com/gabe565/gones/internal/controller"
	"github.com/gabe565/gones/internal/pprof"
	"github.com/gabe565/gones/internal/ppu"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"image/color"
	_ "net/http/pprof"
	"os"
	"os/signal"
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
}

func (r Run) Run() error {
	pprof.Spawn()

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

	c, err := console.New(r.Path)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	errGroup, ctx := errgroup.WithContext(ctx)

	errGroup.Go(func() error {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, os.Interrupt)
		select {
		case <-ctx.Done():
			return nil
		case <-ch:
			return errors.New("interrupt")
		}
	})

	errGroup.Go(func() error {
		var debug console.Debug

		for {
			select {
			case <-ctx.Done():
				return nil
			default:
				for {
					if debug == console.DebugStepFrame {
						debug = c.WaitDebug(ctx, win)
					}

					if err := c.Step(); err != nil {
						if errors.Is(err, console.ErrRender) {
							break
						}
						return err
					}
				}

				if win.Closed() {
					cancel()
					return nil
				}

				c.UpdateInput(win)

				if win.JustPressed(controller.Reset) {
					c.Reset()
				}

				if win.JustPressed(controller.FastForward) {
					win.SetVSync(false)
				} else if win.JustReleased(controller.FastForward) {
					win.SetVSync(true)
				}

				if win.JustPressed(controller.ToggleDebug) {
					log.Info("Enable step debug")
					debug = console.DebugStepFrame
				}

				win.Clear(color.Black)
				pic := c.PPU.Render()
				sprite := pixel.NewSprite(pic, pic.Bounds())
				sprite.Draw(win, pixel.IM.Moved(win.Bounds().Center()).Scaled(win.Bounds().Center(), 3))
				win.Update()

				if debug == console.DebugRunRender {
					debug = console.DebugStepFrame
				}
			}
		}
	})

	return errGroup.Wait()
}
