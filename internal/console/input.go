package console

import (
	"runtime"

	"github.com/gabe565/gones/internal/controller"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	log "github.com/sirupsen/logrus"
)

func (c *Console) CheckInput() {
	c.Bus.UpdateInput()

	if duration := inpututil.KeyPressDuration(ebiten.Key(c.config.Input.Reset)); duration != 0 {
		if duration == c.config.Input.ResetHoldFrames() {
			c.Reset()
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.Key(c.config.Input.FastForward)) {
		if c.player != nil {
			c.player.SetVolume(c.config.Audio.Volume / 2)
		}
		c.SetRate(c.config.Input.FastForwardRate)
	} else if inpututil.IsKeyJustReleased(ebiten.Key(c.config.Input.FastForward)) {
		c.SetRate(1)
		if c.player != nil {
			c.player.SetVolume(c.config.Audio.Volume)
		}
	}

	if inpututil.IsKeyJustPressed(controller.ToggleDebug) {
		if c.debug == DebugDisabled {
			log.Info("Enable step debug")
			c.debug = DebugWait
			c.APU.Enabled = false
		} else {
			log.Info("Disable step debug")
			c.enableTrace = false
			c.debug = DebugDisabled
			c.APU.Enabled = true
		}
	}

	if c.debug != DebugDisabled {
		if inpututil.IsKeyJustPressed(controller.ToggleTrace) {
			log.Info("Toggle trace logs")
			c.enableTrace = !c.enableTrace
		}
		if inpututil.IsKeyJustPressed(controller.StepFrame) || inpututil.KeyPressDuration(controller.StepFrame) > 30 {
			c.debug = DebugStepFrame
		}
		if inpututil.IsKeyJustPressed(controller.RunToRender) || inpututil.KeyPressDuration(controller.RunToRender) > 30 {
			c.debug = DebugRunRender
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.Key(c.config.Input.Fullscreen)) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}

	if inpututil.IsKeyJustPressed(ebiten.Key(c.config.Input.State1Save)) {
		if ebiten.IsKeyPressed(ebiten.Key(c.config.Input.StateUndoModifier)) {
			if err := c.UndoSaveState(); err == nil {
				log.Info("Undo save state")
			} else {
				log.WithError(err).Error("Failed to undo save state")
			}
		} else {
			if err := c.SaveStateNum(1, true); err != nil {
				log.WithError(err).Error("Failed to save state")
			}
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.Key(c.config.Input.State1Load)) {
		if ebiten.IsKeyPressed(ebiten.Key(c.config.Input.StateUndoModifier)) {
			if err := c.UndoLoadState(); err == nil {
				log.Info("Undo load state")
			} else {
				log.WithError(err).Error("Failed to undo load state")
			}
		} else {
			if err := c.LoadStateNum(1); err != nil {
				log.WithError(err).Error("Failed to load state")
			}
		}
	}

	//goland:noinspection GoBoolExpressions
	if inpututil.IsKeyJustPressed(ebiten.Key(c.config.Input.Screenshot)) && runtime.GOOS != "js" {
		c.willScreenshot = true
	}
}
