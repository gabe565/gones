package bus

import (
	"fmt"
	"github.com/gabe565/gones/internal/apu"
	"github.com/gabe565/gones/internal/cartridge"
	"github.com/gabe565/gones/internal/controller"
	"github.com/gabe565/gones/internal/ppu"
	log "github.com/sirupsen/logrus"
)

func New(mapper cartridge.Mapper, ppu *ppu.PPU, apu *apu.APU) *Bus {
	return &Bus{
		mapper: mapper,
		apu:    apu,
		ppu:    ppu,
		controller1: controller.Controller{
			Keymap: controller.Player1Keymap,
		},
		controller2: controller.Controller{
			Keymap: controller.Player2Keymap,
		},
	}
}

type Bus struct {
	CpuVram     [0x800]byte
	mapper      cartridge.Mapper
	apu         *apu.APU
	ppu         *ppu.PPU
	controller1 controller.Controller
	controller2 controller.Controller
}

func (b *Bus) ReadMem(addr uint16) byte {
	switch {
	case addr < 0x2000:
		addr &= 0x07FF
		return b.CpuVram[addr]
	case addr == 0x2000, addr == 0x2001, addr == 0x2003, addr == 0x2005, addr == 0x2006, addr == 0x4014:
		return 0
	case addr == 0x2002:
		return b.ppu.ReadStatus()
	case addr == 0x2004:
		return b.ppu.ReadOam()
	case addr == 0x2007:
		return b.ppu.Read()
	case 0x2008 <= addr && addr < 0x4000:
		addr &= 0x2007
		return b.ReadMem(addr)
	case 0x4000 <= addr && addr < 0x4016:
		return b.apu.ReadMem(addr)
	case addr == 0x4016:
		return b.controller1.Read()
	case addr == 0x4017:
		return b.controller2.Read()
	case addr <= 0x4018 && addr < 0x4020:
		// Disabled
	default:
		return b.mapper.ReadMem(addr)
	}
	return 0
}

func (b *Bus) WriteMem(addr uint16, data byte) {
	switch {
	case addr < 0x2000:
		addr &= 0x07FF
		b.CpuVram[addr] = data
	case addr == 0x2000:
		b.ppu.WriteCtrl(data)
	case addr == 0x2001:
		b.ppu.WriteMask(data)
	case addr == 0x2002:
		log.WithField("address", fmt.Sprintf("%02X", addr)).
			Error("attempt to write to PPU status register")
	case addr == 0x2003:
		b.ppu.WriteOamAddr(data)
	case addr == 0x2004:
		b.ppu.WriteOam(data)
	case addr == 0x2005:
		b.ppu.WriteScroll(data)
	case addr == 0x2006:
		b.ppu.WriteAddr(data)
	case addr == 0x2007:
		b.ppu.Write(data)
	case 0x2008 <= addr && addr < 0x4000:
		addr &= 0x2007
		b.WriteMem(addr, data)
	case addr == 0x4014:
		var buf [256]byte
		hi := uint16(data) << 8
		for k := range buf {
			buf[k] = b.ReadMem(hi + uint16(k))
		}
		b.ppu.WriteOamDma(buf)
	case 0x4000 <= addr && addr < 0x4013, addr == 0x4015, addr == 0x4017:
		b.apu.WriteMem(addr, data)
	case addr == 0x4016:
		b.controller1.Write(data)
		b.controller2.Write(data)
	case addr <= 0x4018 && addr < 0x4020:
		// Disabled
	default:
		b.mapper.WriteMem(addr, data)
	}
}

func (b *Bus) UpdateInput() {
	b.controller1.UpdateInput()
	b.controller2.UpdateInput()
}
