package cartridge

import (
	"github.com/gabe565/gones/internal/consts"
	log "github.com/sirupsen/logrus"
)

func NewMapper2(cartridge *Cartridge) Mapper {
	prgBanks := uint16(len(cartridge.Prg) / consts.PrgChunkSize)
	return &Mapper2{
		Cartridge: cartridge,
		prgBanks:  uint(prgBanks),
		prgBank2:  uint(prgBanks - 1),
	}
}

type Mapper2 struct {
	*Cartridge
	prgBanks uint
	prgBank1 uint
	prgBank2 uint
}

func (m *Mapper2) Step() {}

func (m *Mapper2) Read(addr uint16) byte {
	switch {
	case addr < 0x2000:
		return m.Chr[addr]
	case 0x6000 <= addr && addr < 0x8000:
		addr -= 0x6000
		return m.Sram[addr]
	case 0x8000 <= addr && addr < 0xC000:
		addr := uint(addr)
		addr -= 0x8000
		addr += m.prgBank1 * consts.PrgChunkSize
		return m.Prg[addr]
	case 0xC000 <= addr:
		addr := uint(addr)
		addr -= 0xC000
		addr += m.prgBank2 * consts.PrgChunkSize
		return m.Prg[addr]
	default:
		log.Fatalf("invalid mapper 2 read from $%04X", addr)
		return 0
	}
}

func (m *Mapper2) Write(addr uint16, data byte) {
	switch {
	case addr < 0x2000:
		m.Chr[addr] = data
	case 0x6000 <= addr && addr < 0x8000:
		addr := uint(addr)
		addr -= 0x6000
		m.Sram[addr] = data
	case 0x8000 <= addr:
		addr := uint(addr)
		addr %= m.prgBanks
		m.prgBank1 = addr
	default:
		log.Fatalf("invalid mapper 2 write to $%04X", addr)
	}
}
