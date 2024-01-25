package console

import "fmt"

func (c *Console) Trace() string {
	return fmt.Sprintf(
		"%s PPU:%3d,%3d CYC:%d",
		c.CPU.Trace(),
		c.PPU.Scanline,
		c.PPU.Cycles,
		c.CPU.Cycles(),
	)
}
