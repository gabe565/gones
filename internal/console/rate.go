//go:build !js

package console

import "gabe565.com/gones/internal/apu"

func (c *Console) SetRate(rate uint8) {
	c.rate = rate
	c.APU.Clear()
	c.APU.SampleRate = apu.DefaultSampleRate * float64(rate)
}
