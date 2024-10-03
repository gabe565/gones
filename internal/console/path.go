//go:build !js

package console

import (
	"fmt"
	"path/filepath"

	"github.com/gabe565/gones/internal/config"
)

func (c *Console) SRAMPath() (string, error) {
	sramDir, err := config.GetSRAMDir()
	if err != nil {
		return "", err
	}

	sramName := c.Cartridge.Hash() + ".sav"
	return filepath.Join(sramDir, sramName), nil
}

func (c *Console) StatePath(num uint8) (string, error) {
	statesDir, err := config.GetStatesDir()
	if err != nil {
		return "", err
	}

	stateName := fmt.Sprintf("%s.%d.state.gz", c.Cartridge.Hash(), num)
	return filepath.Join(statesDir, stateName), nil
}
