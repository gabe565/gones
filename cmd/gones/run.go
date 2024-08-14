package gones

import (
	"context"
	"errors"
	"log/slog"
	"runtime"

	"github.com/gabe565/gones/internal/config"
	"github.com/gabe565/gones/internal/console"
	"github.com/gabe565/gones/internal/pprof"
	"github.com/gabe565/gones/internal/ppu"
	"github.com/hajimehoshi/ebiten/v2"
)

func run(ctx context.Context, conf *config.Config, path string) error {
	if pprof.Enabled {
		go func() {
			if err := pprof.ListenAndServe(); err != nil {
				slog.Error("Failed to start pprof", "error", err)
			}
		}()
	}

	c, err := newConsole(conf, path)
	if err != nil {
		return err
	}
	defer func() {
		if err := c.Close(); err != nil {
			slog.Error("Failed to close console", "error", err)
		}
	}()

	if runtime.GOOS != "js" {
		go func() {
			<-ctx.Done()
			slog.Info("Exiting...")
			c.SetUpdateAction(console.ActionExit)
		}()
	}

	scale := conf.UI.Scale
	ebiten.SetWindowSize(int(scale*ppu.Width), int(scale*ppu.TrimmedHeight))
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetFullscreen(conf.UI.Fullscreen)
	ebiten.SetScreenClearedEveryFrame(false)
	ebiten.SetRunnableOnUnfocused(!conf.UI.PauseUnfocused)
	setWindowIcons()

	if name := c.Cartridge.Name(); name != "" {
		ebiten.SetWindowTitle(name + " | GoNES")
	}

	if err := ebiten.RunGameWithOptions(c, &ebiten.RunGameOptions{
		SingleThread: true,
	}); err != nil && !errors.Is(err, console.ErrExit) {
		return err
	}

	return nil
}
