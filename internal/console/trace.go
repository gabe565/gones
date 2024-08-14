package console

import (
	"fmt"
	"runtime"
)

func (c *Console) Trace() string {
	if runtime.GOOS == "js" {
		return "DISABLED"
	}
	return fmt.Sprintf(
		"%s PPU:%3d,%3d CYC:%d",
		c.CPU.Trace(),
		c.PPU.Scanline,
		c.PPU.Cycles,
		c.CPU.GetCycles(),
	)
}
