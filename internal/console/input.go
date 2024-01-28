package console

import (
	"github.com/gabe565/gones/internal/controller"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	log "github.com/sirupsen/logrus"
)

func (c *Console) CheckInput() {
	c.Bus.UpdateInput()

	if inpututil.KeyPressDuration(controller.Reset) == 30 {
		c.Reset()
	}

	if inpututil.IsKeyJustPressed(controller.FastForward) {
		c.player.SetVolume(0.4)
		c.SetRate(3)
	} else if inpututil.IsKeyJustReleased(controller.FastForward) {
		c.SetRate(1)
		c.player.SetVolume(1)
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

	if inpututil.IsKeyJustPressed(controller.ToggleFullscreen) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}

	if inpututil.IsKeyJustPressed(controller.SaveState1) {
		if err := c.SaveState(1); err != nil {
			log.WithError(err).Error("Failed to save state")
		}
	}

	if inpututil.IsKeyJustPressed(controller.LoadState1) {
		if err := c.LoadState(1); err != nil {
			log.WithError(err).Error("Failed to load state")
		}
	}
}
