package cmd

import (
	"context"
	"errors"
	"github.com/gabe565/gones/internal/console"
	"github.com/gabe565/gones/internal/pprof"
	"github.com/gabe565/gones/internal/ppu"
	"github.com/hajimehoshi/ebiten/v2"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
)

type ContextKey uint8

const (
	ConfigKey ContextKey = iota
)

type Config struct {
	Path       string
	Debug      bool
	Trace      bool
	Scale      float64
	Fullscreen bool
}

func New(version string) *cobra.Command {
	var config Config

	cmd := &cobra.Command{
		Use:     "gones ROM",
		Version: version,
		RunE:    run,
	}

	cmd.Flags().BoolVar(&config.Debug, "debug", false, "Start with step debugging enabled")
	cmd.Flags().BoolVar(&config.Trace, "trace", false, "Enable trace logging")
	cmd.Flags().Float64Var(&config.Scale, "scale", 3, "Default UI scale")
	cmd.Flags().BoolVarP(&config.Fullscreen, "fullscreen", "f", false, "Start in fullscreen")
	pprof.Flag(cmd)

	ctx := context.Background()
	ctx = context.WithValue(ctx, ConfigKey, &config)
	cmd.SetContext(ctx)

	return cmd
}

func run(cmd *cobra.Command, args []string) error {
	pprof.Spawn()

	config := cmd.Context().Value(ConfigKey).(*Config)

	if len(args) > 0 {
		config.Path = args[0]
	}
	cmd.SilenceUsage = true

	c, err := newConsole(config.Path)
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

	c.SetTrace(config.Trace)
	c.SetDebug(config.Debug)

	ebiten.SetWindowSize(int(config.Scale*ppu.Width), int(config.Scale*ppu.TrimmedHeight))
	ebiten.SetWindowTitle(strings.TrimSuffix(filepath.Base(config.Path), filepath.Ext(config.Path)) + " | GoNES")
	ebiten.SetScreenFilterEnabled(false)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowClosingHandled(true)
	ebiten.SetFullscreen(config.Fullscreen)

	if err := ebiten.RunGame(c); err != nil && !errors.Is(err, console.ErrExit) {
		return err
	}

	return nil
}
