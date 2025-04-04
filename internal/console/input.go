package console

import (
	"log/slog"
	"runtime"

	"gabe565.com/gones/internal/controller"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func (c *Console) CheckInput() {
	c.Bus.UpdateInput()

	if duration := inpututil.KeyPressDuration(ebiten.Key(c.Config.Input.Reset)); duration != 0 {
		if duration == c.Config.Input.ResetHoldFrames() {
			c.Reset()
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.Key(c.Config.Input.FastForward)) {
		if c.player != nil {
			c.player.SetVolume(c.Config.Audio.Volume / 2)
		}
		c.SetRate(c.Config.Input.FastForwardRate)
	} else if inpututil.IsKeyJustReleased(ebiten.Key(c.Config.Input.FastForward)) {
		c.SetRate(1)
		if c.player != nil {
			c.player.SetVolume(c.Config.Audio.Volume)
		}
	}

	if runtime.GOOS != "js" {
		if inpututil.IsKeyJustPressed(controller.ToggleDebug) {
			if c.debug == DebugDisabled {
				slog.Info("Enable step debug")
				c.debug = DebugWait
				c.APU.Enabled = false
			} else {
				slog.Info("Disable step debug")
				c.enableTrace = false
				c.debug = DebugDisabled
				c.APU.Enabled = true
			}
		}

		if c.debug != DebugDisabled {
			if inpututil.IsKeyJustPressed(controller.ToggleTrace) {
				slog.Info("Toggle trace logs")
				c.enableTrace = !c.enableTrace
			}
			if inpututil.IsKeyJustPressed(controller.StepFrame) || inpututil.KeyPressDuration(controller.StepFrame) > 30 {
				c.debug = DebugStepFrame
			}
			if inpututil.IsKeyJustPressed(controller.RunToRender) || inpututil.KeyPressDuration(controller.RunToRender) > 30 {
				c.debug = DebugRunRender
			}
		}

		if inpututil.IsKeyJustPressed(ebiten.Key(c.Config.Input.Screenshot)) {
			c.willScreenshot = true
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.Key(c.Config.Input.Fullscreen)) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}

	if inpututil.IsKeyJustPressed(ebiten.Key(c.Config.Input.State1Save)) {
		if ebiten.IsKeyPressed(ebiten.Key(c.Config.Input.StateUndoModifier)) {
			if err := c.UndoSaveState(); err == nil {
				slog.Info("Undo save state")
			} else {
				slog.Error("Failed to undo save state", "error", err)
			}
		} else {
			if err := c.SaveStateNum(1, true); err != nil {
				slog.Error("Failed to save state", "error", err)
			}
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.Key(c.Config.Input.State1Load)) {
		if ebiten.IsKeyPressed(ebiten.Key(c.Config.Input.StateUndoModifier)) {
			if err := c.UndoLoadState(); err == nil {
				slog.Info("Undo load state")
			} else {
				slog.Error("Failed to undo load state", "error", err)
			}
		} else {
			if err := c.LoadStateNum(1); err != nil {
				slog.Error("Failed to load state", "error", err)
			}
		}
	}
}
