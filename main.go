package main

import (
	"github.com/gabe565/gones/internal/console"
	"github.com/gabe565/gones/internal/pprof"
	"github.com/gabe565/gones/internal/ppu"
	"github.com/hajimehoshi/ebiten/v2"
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
	Scale float64
}

func (r Run) Run() error {
	pprof.Spawn()

	c, err := console.New(r.Path)
	if err != nil {
		return err
	}

	ebiten.SetWindowSize(int(r.Scale*ppu.Width), int(r.Scale*ppu.TrimmedHeight))
	ebiten.SetWindowTitle(filepath.Base(r.Path) + " | GoNES")
	ebiten.SetScreenFilterEnabled(false)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	return ebiten.RunGame(c)
}
