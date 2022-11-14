package console

import (
	"context"
	"github.com/faiface/pixel/pixelgl"
	"github.com/gabe565/gones/internal/controller"
	log "github.com/sirupsen/logrus"
	"time"
)

type Debug uint8

const (
	DebugDisabled = iota
	DebugStepFrame
	DebugRunRender
)

func (c *Console) WaitDebug(ctx context.Context, win *pixelgl.Window) Debug {
	for {
		select {
		case <-ctx.Done():
			return DebugDisabled
		default:
			if win.Closed() {
				return DebugDisabled
			}

			win.UpdateInputWait(time.Second / 60)
			switch {
			case win.JustPressed(controller.ToggleDebug):
				log.Info("Disable step debug")
				c.EnableTrace = false
				win.UpdateInput()
				return DebugDisabled
			case win.JustPressed(controller.ToggleTrace):
				log.Info("Toggle trace logs")
				c.EnableTrace = !c.EnableTrace
			case win.JustPressed(controller.StepFrame) || win.Repeated(controller.StepFrame):
				return DebugStepFrame
			case win.JustPressed(controller.RunToRender) || win.Repeated(controller.RunToRender):
				win.UpdateInput()
				return DebugRunRender
			}
		}
	}
}
