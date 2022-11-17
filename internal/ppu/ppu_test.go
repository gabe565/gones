package ppu

import (
	"github.com/gabe565/gones/internal/cartridge"
	"github.com/gabe565/gones/internal/ppu/registers"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPPU_VramWrite(t *testing.T) {
	var ppu PPU
	ppu.cartridge = &cartridge.Cartridge{}
	ppu.WriteAddr(0x23)
	ppu.WriteAddr(0x05)
	ppu.Write(0x66)
	assert.EqualValues(t, 0x66, ppu.Vram[0x305])
}

func TestPPU_VramRead(t *testing.T) {
	var ppu PPU
	ppu.cartridge = &cartridge.Cartridge{}
	ppu.Vram[0x305] = 0x66
	ppu.WriteAddr(0x23)
	ppu.WriteAddr(0x05)
	ppu.Read() // Buffer
	assert.EqualValues(t, 0x2306, ppu.Addr.Get())
	assert.EqualValues(t, 0x66, ppu.Read())
}

func TestPPU_VramRead_CrossPage(t *testing.T) {
	var ppu PPU
	ppu.cartridge = &cartridge.Cartridge{}
	ppu.Vram[0x1FF] = 0x66
	ppu.Vram[0x200] = 0x77
	ppu.WriteAddr(0x21)
	ppu.WriteAddr(0xFF)
	ppu.Read() // Buffer
	assert.EqualValues(t, 0x66, ppu.Read())
	assert.EqualValues(t, 0x77, ppu.Read())
}

func TestPPU_VramRead_Step32(t *testing.T) {
	var ppu PPU
	ppu.cartridge = &cartridge.Cartridge{}
	ppu.WriteCtrl(0b100)
	ppu.Vram[0x1FF] = 0x66
	ppu.Vram[0x1FF+32] = 0x77
	ppu.Vram[0x1FF+64] = 0x88
	ppu.WriteAddr(0x21)
	ppu.WriteAddr(0xFF)
	ppu.Read() // Buffer
	assert.EqualValues(t, 0x66, ppu.Read())
	assert.EqualValues(t, 0x77, ppu.Read())
	assert.EqualValues(t, 0x88, ppu.Read())
}

func TestPPU_HorizontalMirror(t *testing.T) {
	var ppu PPU
	ppu.cartridge = &cartridge.Cartridge{}
	ppu.WriteAddr(0x24)
	ppu.WriteAddr(0x05)
	ppu.Write(0x66) // A

	ppu.WriteAddr(0x28)
	ppu.WriteAddr(0x05)
	ppu.Write(0x77) // B

	ppu.WriteAddr(0x20)
	ppu.WriteAddr(0x05)
	ppu.Read() // Buffer

	assert.EqualValues(t, 0x66, ppu.Read()) // A

	ppu.WriteAddr(0x2C)
	ppu.WriteAddr(0x05)

	ppu.Read()
	assert.EqualValues(t, 0x77, ppu.Read()) // B
}

func TestPPU_VerticalMirror(t *testing.T) {
	var ppu PPU
	ppu.cartridge = &cartridge.Cartridge{}
	ppu.cartridge.Mirror = cartridge.Vertical
	ppu.WriteAddr(0x20)
	ppu.WriteAddr(0x05)
	ppu.Write(0x66) // A

	ppu.WriteAddr(0x2C)
	ppu.WriteAddr(0x05)
	ppu.Write(0x77) // B

	ppu.WriteAddr(0x28)
	ppu.WriteAddr(0x05)
	ppu.Read() // Buffer

	assert.EqualValues(t, 0x66, ppu.Read()) // A

	ppu.WriteAddr(0x24)
	ppu.WriteAddr(0x05)

	ppu.Read()
	assert.EqualValues(t, 0x77, ppu.Read()) // B
}

func TestPPU_StatusResetsLatch(t *testing.T) {
	var ppu PPU
	ppu.cartridge = &cartridge.Cartridge{}
	ppu.cartridge.Chr = make([]byte, 2048)
	ppu.mapper = cartridge.NewMapper2(ppu.cartridge)
	ppu.Vram[0x305] = 0x66
	ppu.WriteAddr(0x21)
	ppu.WriteAddr(0x23)
	ppu.WriteAddr(0x05)
	ppu.Read() // Buffer
	assert.NotEqualValues(t, 0x66, ppu.Read())

	ppu.ReadStatus()
	ppu.WriteAddr(0x23)
	ppu.WriteAddr(0x05)
	ppu.Read()
	assert.EqualValues(t, 0x66, ppu.Read())
}

func TestPPU_VramMirror(t *testing.T) {
	var ppu PPU
	ppu.cartridge = &cartridge.Cartridge{}
	ppu.Vram[0x305] = 0x66

	ppu.WriteAddr(0x63)
	ppu.WriteAddr(0x05)

	ppu.Read() // Buffer
	assert.EqualValues(t, 0x66, ppu.Read())
}

func TestPPU_StatusResetsVblank(t *testing.T) {
	var ppu PPU
	ppu.cartridge = &cartridge.Cartridge{}
	ppu.cartridge.Chr = make([]byte, 2048)
	ppu.Status.Insert(registers.Vblank)
	assert.EqualValues(t, 1, ppu.ReadStatus()>>7)
	assert.EqualValues(t, 0, ppu.ReadStatus()>>7)
}

func TestCPU_OamReadWrite(t *testing.T) {
	var ppu PPU
	ppu.cartridge = &cartridge.Cartridge{}
	ppu.WriteOamAddr(0x10)
	ppu.WriteOam(0x66)
	ppu.WriteOam(0x77)
	ppu.WriteOamAddr(0x10)
	assert.EqualValues(t, 0x66, ppu.ReadOam())
	ppu.WriteOamAddr(0x11)
	assert.EqualValues(t, 0x77, ppu.ReadOam())
}

func TestCPU_OamDma(t *testing.T) {
	var ppu PPU
	ppu.cartridge = &cartridge.Cartridge{}

	var data [256]byte
	for k := range data {
		data[k] = 0x66
	}
	data[0] = 0x77
	data[255] = 0x88

	ppu.WriteOamAddr(0x10)
	ppu.WriteOamDma(data)

	ppu.WriteOamAddr(0xF) // Wrap around
	assert.EqualValues(t, 0x88, ppu.ReadOam())

	ppu.WriteOamAddr(0x10)
	assert.EqualValues(t, 0x77, ppu.ReadOam())

	ppu.WriteOamAddr(0x11)
	assert.EqualValues(t, 0x66, ppu.ReadOam())
}
