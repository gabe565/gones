package cartridge

import (
	"fmt"
	"path/filepath"

	"github.com/gabe565/gones/internal/config"
)

func (c *Cartridge) SramPath() (string, error) {
	sramDir, err := config.GetSramDir()
	if err != nil {
		return "", err
	}

	sramName := fmt.Sprintf("%s.sav", c.hash)
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
