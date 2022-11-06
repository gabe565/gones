package console

import (
	"github.com/gabe565/gones/internal/bus"
	"github.com/gabe565/gones/internal/cartridge"
	"github.com/gabe565/gones/internal/cpu"
)

func New(path string) (*cpu.CPU, error) {
	cart, err := cartridge.FromiNes(path)
	if err != nil {
		return nil, err
	}

	b := bus.New(cart)
	return cpu.New(b), nil
}
