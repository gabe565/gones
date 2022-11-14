package console

import (
	"github.com/gabe565/gones/internal/controller"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	log "github.com/sirupsen/logrus"
)

func (c *Console) CheckInput() {
	c.Bus.Controller1.UpdateInput()
	c.Bus.Controller2.UpdateInput()

	if ebiten.IsKeyPressed(controller.Reset) {
		c.Reset()
	}

	if inpututil.IsKeyJustPressed(controller.FastForward) {
		ebiten.SetTPS(3 * 60)
	} else if inpututil.IsKeyJustReleased(controller.FastForward) {
		ebiten.SetTPS(60)
	}

	if inpututil.IsKeyJustPressed(controller.ToggleDebug) {
		if c.Debug == DebugDisabled {
			log.Info("Enable step debug")
			c.Debug = DebugWait
		} else {
			log.Info("Disable step debug")
			c.EnableTrace = false
			c.Debug = DebugDisabled
		}
	}

	if c.Debug != DebugDisabled {
		if inpututil.IsKeyJustPressed(controller.ToggleTrace) {
			log.Info("Toggle trace logs")
			c.EnableTrace = !c.EnableTrace
		}
		if inpututil.IsKeyJustPressed(controller.StepFrame) || inpututil.KeyPressDuration(controller.StepFrame) > 30 {
			c.Debug = DebugStepFrame
		}
		if inpututil.IsKeyJustPressed(controller.RunToRender) || inpututil.KeyPressDuration(controller.RunToRender) > 30 {
			c.Debug = DebugRunRender
		}
	}
}
