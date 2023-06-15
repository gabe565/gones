package cartridge

import (
	"fmt"
)

func (c *Cartridge) SramPath() (string, error) {
	return fmt.Sprintf("%s.sav", c.hash), nil
}

func (c *Cartridge) StatePath(num uint8) (string, error) {
	return fmt.Sprintf("%s.%d.state.gz", c.hash, num), nil
}
