package main

import (
	"context"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/gabe565/gones/internal/console"
	"github.com/gabe565/gones/internal/controller"
	"github.com/gabe565/gones/internal/pprof"
	"github.com/gabe565/gones/internal/ppu"
	log "github.com/sirupsen/logrus"
	"image/color"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"path/filepath"
	"time"
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

	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, os.Interrupt)
		for range ch {
			cancel()
		}
	}()

	var debugNextRender bool
	for {
		select {
		case err := <-errCh:
			return err
		case pic := <-console.Bus.Render:
			if win.Closed() {
				return nil
			}

			console.Bus.Controller1.UpdateInput(win)
			console.Bus.Controller2.UpdateInput(win)

			if win.JustPressed(controller.Reset) {
				console.CPU.ResetCh <- struct{}{}
			}

			if win.JustPressed(controller.FastForward) {
				win.SetVSync(false)
			} else if win.JustReleased(controller.FastForward) {
				win.SetVSync(true)
			}

			if win.JustPressed(controller.ToggleDebug) && !console.CPU.EnableDebug || debugNextRender {
				if !debugNextRender {
					log.Info("Enable step debug")
				}
				debugNextRender = false
				console.CPU.EnableDebug = true
				go func() {
					defer func() {
						console.CPU.EnableDebug = false
						console.CPU.DebugCh <- struct{}{}
					}()
					for {
						select {
						case <-ctx.Done():
							return
						default:
							if win.Closed() {
								return
							}

							win.UpdateInputWait(time.Millisecond)

							if win.JustPressed(controller.ToggleDebug) {
								win.UpdateInput()
								console.CPU.EnableTrace = false
								log.Info("Disable step debug")
								return
							} else if win.JustPressed(controller.StepFrame) || win.Repeated(controller.StepFrame) {
								console.CPU.DebugCh <- struct{}{}
							} else if win.JustPressed(controller.RunToRender) || win.Repeated(controller.RunToRender) {
								debugNextRender = true
								return
							} else if win.JustPressed(controller.ToggleTrace) {
								log.Info("Toggle trace logs")
								console.CPU.EnableTrace = !console.CPU.EnableTrace
							}
						}
					}
				}()
			}

			win.Clear(color.Black)
			sprite := pixel.NewSprite(pic, pic.Bounds())
			sprite.Draw(win, pixel.IM.Moved(win.Bounds().Center()).Scaled(win.Bounds().Center(), 3))
			win.Update()
		}
	}
}
