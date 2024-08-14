package console

func (c *Console) SetRate(rate uint8) {
	c.rate = rate
	if rate == 1 {
		c.APU.Enabled = true
		c.player.SetVolume(1)
	} else {
		c.APU.Enabled = false
		c.player.SetVolume(0)
	}
}
