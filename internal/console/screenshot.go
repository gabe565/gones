//go:build !js

package console

import (
	"image/png"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"gabe565.com/gones/internal/config"
	"github.com/hajimehoshi/ebiten/v2"
)

func (c *Console) writeScreenshot(screen *ebiten.Image) error {
	dir, err := config.GetScreenshotDir()
	if err != nil {
		return err
	}

	gameDir := filepath.Join(dir, c.Cartridge.Name())
	filename := filepath.Join(gameDir, time.Now().Format("2006-01-02_150405")+".png")

	if err := os.MkdirAll(gameDir, 0o777); err != nil {
		return err
	}

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	if err := png.Encode(f, screen); err != nil {
		return err
	}

	if err := f.Close(); err != nil {
		return err
	}

	slog.Info("Saved screenshot", "path", filename)
	return nil
}
