package console

import (
	"fmt"
)

func (c *Console) SRAMPath() (string, error) {
	return fmt.Sprintf("%s.sav", c.Cartridge.Hash()), nil
}

func (c *Console) StatePath(num uint8) (string, error) {
	return fmt.Sprintf("%s.%d.state.gz", c.Cartridge.Hash(), num), nil
}
