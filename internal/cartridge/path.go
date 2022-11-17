package cartridge

import (
	"fmt"
	"github.com/gabe565/gones/internal/config"
	"path/filepath"
)

func (c *Cartridge) StatePath(num uint8) (string, error) {
	statesDir, err := config.GetStatesDir()
	if err != nil {
		return "", err
	}

	stateName := fmt.Sprintf("%s.%d.state", c.hash, num)
	return filepath.Join(statesDir, stateName), nil
}
