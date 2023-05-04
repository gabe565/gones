package cartridge

import (
	"errors"
	"fmt"

	"github.com/gabe565/gones/internal/interrupts"
	"github.com/gabe565/gones/internal/memory"
)

type CPU interface {
	interrupts.Interruptible
}

type Mapper interface {
	memory.ReadWrite8
	Cartridge() *Cartridge
	SetCartridge(*Cartridge)
}

type MapperInterrupts interface {
	SetCpu(CPU)
}

type MapperPPU interface {
	StepPPU(renderEnabled bool, scanline uint16, cycle uint)
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
	case 71:
		return NewMapper71(cartridge), nil
	default:
		return nil, fmt.Errorf("%w: %d", ErrUnsupportedMapper, cartridge.Mapper)
	}
}
