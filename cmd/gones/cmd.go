package gones

import (
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"runtime"

	"github.com/gabe565/gones/cmd/options"
	"github.com/gabe565/gones/internal/config"
	"github.com/gabe565/gones/internal/console"
	"github.com/gabe565/gones/internal/pprof"
	"github.com/gabe565/gones/internal/ppu"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/spf13/cobra"
)

func New(opts ...options.Option) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gones ROM",
		Short: "NES emulator written in Go",
		RunE:  run,
		ValidArgsFunction: func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
			return []string{"nes"}, cobra.ShellCompDirectiveFilterFileExt
		},

		SilenceErrors:     true,
		DisableAutoGenTag: true,
	}
	config.Flags(cmd)

	for _, opt := range opts {
		opt(cmd)
	}

	return cmd
}

func run(cmd *cobra.Command, args []string) error {
	conf, err := config.Load(cmd)
	if err != nil {
		return err
	}

	var path string
	if len(args) > 0 {
		path = args[0]
	}
	cmd.SilenceUsage = true

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

	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, os.Interrupt)
		for range ch {
			slog.Info("Exiting...")
			c.SetUpdateAction(console.ActionExit)
		}
	}()

	scale := conf.UI.Scale
	ebiten.SetWindowSize(int(scale*ppu.Width), int(scale*ppu.TrimmedHeight))
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetFullscreen(conf.UI.Fullscreen)
	ebiten.SetScreenClearedEveryFrame(false)
	ebiten.SetRunnableOnUnfocused(!conf.UI.PauseUnfocused)
	if runtime.GOOS != "darwin" {
		ebiten.SetWindowIcon(getWindowIcons())
	}

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
