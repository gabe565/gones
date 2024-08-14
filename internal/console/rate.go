//go:build !js

package console

import "github.com/gabe565/gones/internal/apu"

func (c *Console) SetRate(rate uint8) {
	c.rate = rate
	c.APU.Clear()
	c.APU.SampleRate = apu.DefaultSampleRate * float64(rate)
}
