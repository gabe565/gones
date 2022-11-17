package cartridge

import (
	"crypto/md5"
	"fmt"
	"github.com/gabe565/gones/internal/consts"
)

type Cartridge struct {
	hash string

	Prg     []byte
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

	cart.Prg = make([]byte, consts.PrgRomAddr, consts.PrgChunkSize*2)
	cart.Prg = append(cart.Prg, b...)
	cart.Prg = cart.Prg[:cap(cart.Prg)]
	cart.Prg[consts.ResetAddr+1-consts.PrgChunkSize*2] = 0x86

	cart.Chr = make([]byte, consts.ChrChunkSize)

	return cart
}
