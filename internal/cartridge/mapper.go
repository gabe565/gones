package cartridge

import (
	"errors"
	"fmt"

	"gabe565.com/gones/internal/memory"
	"gabe565.com/gones/internal/ppu/registers"
)

type Mapper interface {
	memory.ReadWrite8
	Cartridge() *Cartridge
	SetCartridge(cartridge *Cartridge)
}

type MapperOnCPUStep interface {
	OnCPUStep(cycle uint)
}

type MapperOnScanline interface {
	OnScanline()
}

type MapperOnVRAMAddr interface {
	OnVRAMAddr(addr registers.Address)
}

type MapperIRQ interface {
	IRQ() bool
}

var ErrUnsupportedMapper = errors.New("unsupported mapper")

func NewMapper(cartridge *Cartridge) (Mapper, error) { //nolint:ireturn,nolintlint
	switch cartridge.Header.Mapper() {
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
		return nil, fmt.Errorf("%w: %d", ErrUnsupportedMapper, cartridge.Header.Mapper())
	}
}
