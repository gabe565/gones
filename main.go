package main

import (
	"errors"
	"github.com/gabe565/gones/internal/console"
	"github.com/gabe565/gones/internal/pprof"
	"github.com/gabe565/gones/internal/ppu"
	"github.com/hajimehoshi/ebiten/v2"
	log "github.com/sirupsen/logrus"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
)

func main() {
	if err := NewCommand().Execute(); err != nil {
		os.Exit(1)
	}
}

type Run struct {
	Path       string
	Trace      bool
	Scale      float64
	Fullscreen bool
}

func (r Run) Run() error {
	pprof.Spawn()

	c, err := console.New(r.Path)
	if err != nil {
		return err
	}
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, os.Interrupt)
		for range ch {
			log.Info("Exiting...")
			c.CloseOnUpdate()
		}
	}()

	c.SetTrace(r.Trace)

	ebiten.SetWindowSize(int(r.Scale*ppu.Width), int(r.Scale*ppu.TrimmedHeight))
	ebiten.SetWindowTitle(strings.TrimSuffix(filepath.Base(r.Path), filepath.Ext(r.Path)) + " | GoNES")
	ebiten.SetScreenFilterEnabled(false)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowClosingHandled(true)
	ebiten.SetFullscreen(r.Fullscreen)

	if err := ebiten.RunGame(c); err != nil && !errors.Is(err, console.ErrExit) {
		return err
	}

	return nil
}
