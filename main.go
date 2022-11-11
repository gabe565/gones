package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/gabe565/gones/internal/console"
	"github.com/gabe565/gones/internal/joypad"
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

	c, err := console.New(r.Path, func(ppu *ppu.PPU, joypad1 *joypad.Joypad) {
		for button, key := range joypad.Keymap {
			if win.JustPressed(button) {
				joypad1.Set(key, true)
			} else if win.JustReleased(button) {
				joypad1.Set(key, false)
			}
		}

		win.Clear(color.Black)
		pic := ppu.Render()
		sprite := pixel.NewSprite(pic, pic.Bounds())
		sprite.Draw(win, pixel.IM.Moved(win.Bounds().Center()).Scaled(win.Bounds().Center(), 3))
		win.Update()

		if win.Closed() {
			os.Exit(0)
		}
	})
	if err != nil {
		return err
	}

	if err := c.Run(); err != nil {
		return err
	}

	return nil
}
