package cartridge

import (
	"errors"
	"fmt"

	"github.com/gabe565/gones/internal/memory"
	"github.com/gabe565/gones/internal/ppu/registers"
)

type Mapper interface {
	memory.ReadWrite8
	Cartridge() *Cartridge
	SetCartridge(*Cartridge)
}

type MapperOnCPUStep interface {
	OnCPUStep(uint)
}

type MapperOnScanline interface {
	OnScanline()
}

type MapperOnVRAMAddr interface {
	OnVRAMAddr(registers.Address)
}

type MapperIRQ interface {
	IRQ() bool
}

var ErrUnsupportedMapper = errors.New("unsupported mapper")

func NewMapper(cartridge *Cartridge) (Mapper, error) {
	switch cartridge.Mapper {
	case 0, 2:
		return NewMapper2(cartridge), nil
	case 1:
		return NewMapper1(cartridge), nil
	case 3:
		return NewMapper3(cartridge), nil
	case 4:
		return NewMapper4(cartridge), nil
	case 7:
		return NewMapper7(cartridge), nil
	case 69:
		return NewMapper69(cartridge), nil
	case 71:
		return NewMapper71(cartridge), nil
	default:
		return nil, fmt.Errorf("%w: %d", ErrUnsupportedMapper, cartridge.Mapper)
	}
}
