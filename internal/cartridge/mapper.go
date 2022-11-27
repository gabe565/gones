package cartridge

import (
	"errors"
	"fmt"
	"github.com/gabe565/gones/internal/memory"
)

type Mapper interface {
	memory.ReadWrite
	Cartridge() *Cartridge
	Step()
}

var ErrUnsupportedMapper = errors.New("unsupported mapper")

func NewMapper(cartridge *Cartridge) (Mapper, error) {
	switch cartridge.Mapper {
	case 0, 2:
		return NewMapper2(cartridge), nil
	case 1:
		return NewMapper1(cartridge), nil
	default:
		return nil, fmt.Errorf("%w: %d", ErrUnsupportedMapper, cartridge.Mapper)
	}
}
