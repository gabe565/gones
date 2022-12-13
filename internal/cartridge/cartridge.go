package cartridge

import (
	"crypto/md5"
	"fmt"
	"github.com/gabe565/gones/internal/consts"
	"github.com/gabe565/gones/internal/interrupts"
)

type Cartridge struct {
	hash string

	prg     []byte
	Chr     []byte
	Sram    []byte
	Mapper  byte
	Mirror  Mirror
	Battery bool
}

func New() *Cartridge {
	return &Cartridge{
		Sram: make([]byte, 0x2000),
	}
}

func FromBytes(b []byte) *Cartridge {
	cart := New()
	cart.hash = fmt.Sprintf("%x", md5.Sum(b))

	cart.prg = make([]byte, consts.PrgRomAddr, consts.PrgChunkSize*2)
	cart.prg = append(cart.prg, b...)
	cart.prg = cart.prg[:cap(cart.prg)]
	cart.prg[interrupts.Reset.VectorAddr+1-consts.PrgChunkSize*2] = 0x86

	cart.Chr = make([]byte, consts.ChrChunkSize)

	return cart
}
