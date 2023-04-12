package cmd

import (
	"errors"
	"os"
	"os/signal"
	"path/filepath"
	"strings"

	"github.com/gabe565/gones/internal/config"
	"github.com/gabe565/gones/internal/console"
	"github.com/gabe565/gones/internal/controller"
	"github.com/gabe565/gones/internal/ppu"
	"github.com/hajimehoshi/ebiten/v2"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func New(version string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "gones ROM",
		Version: version,
		RunE:    run,

		SilenceErrors: true,
	}
	config.Flags(cmd)

	return cmd
}

func run(cmd *cobra.Command, args []string) error {
	if err := config.Load(cmd); err != nil {
		return err
	}
	controller.LoadKeys()

	if len(args) > 0 {
		config.Path = args[0]
	}
	path := config.Path
	cmd.SilenceUsage = true

	c, err := newConsole(path)
	if err != nil {
		return err
	}
	defer func() {
		if err := c.Close(); err != nil {
			log.Error(err)
		}
	}()

	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, os.Interrupt)
		for range ch {
			log.Info("Exiting...")
			c.CloseOnUpdate()
		}
	}()

	c.SetTrace(config.K.Bool("debug.trace"))
	c.SetDebug(config.K.Bool("debug.enabled"))

	scale := config.K.Float64("ui.scale")
	ebiten.SetWindowSize(int(scale*ppu.Width), int(scale*ppu.TrimmedHeight))
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetFullscreen(config.K.Bool("ui.fullscreen"))
	ebiten.SetScreenClearedEveryFrame(false)

	name := c.Cartridge.Name()
	if name == "" {
		name = strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	}
	ebiten.SetWindowTitle(name + " | GoNES")

	if err := ebiten.RunGame(c); err != nil && !errors.Is(err, console.ErrExit) {
		return err
	}

	return nil
}
