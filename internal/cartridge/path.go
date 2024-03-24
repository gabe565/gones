//go:build !wasm

package cartridge

import (
	"fmt"
	"path/filepath"

	"github.com/gabe565/gones/internal/config"
)

func (c *Cartridge) SRAMPath() (string, error) {
	sramDir, err := config.GetSRAMDir()
	if err != nil {
		return "", err
	}

	sramName := c.hash + ".sav"
	return filepath.Join(sramDir, sramName), nil
}

func (c *Cartridge) StatePath(num uint8) (string, error) {
	statesDir, err := config.GetStatesDir()
	if err != nil {
		return "", err
	}

	stateName := fmt.Sprintf("%s.%d.state.gz", c.hash, num)
	return filepath.Join(statesDir, stateName), nil
}
