package ppu

import (
	"github.com/gabe565/gones/internal/cartridge"
	"github.com/stretchr/testify/assert"
	"testing"
)

func stubPpu() (*PPU, *cartridge.Cartridge) {
	var ppu PPU
	cart := &cartridge.Cartridge{}
	cart.Chr = make([]byte, 2048)
	ppu.mapper = cartridge.NewMapper2(cart)
	return &ppu, cart
}

func TestPPU_VramWrite(t *testing.T) {
	ppu, _ := stubPpu()
	ppu.WriteAddr(0x23)
	ppu.WriteAddr(0x05)
	ppu.WriteData(0x66)
	assert.EqualValues(t, 0x66, ppu.Vram[0x305])
}

func TestPPU_VramRead(t *testing.T) {
	ppu, _ := stubPpu()
	ppu.Vram[0x305] = 0x66
	ppu.WriteAddr(0x23)
	ppu.WriteAddr(0x05)
	ppu.ReadData() // Buffer
	assert.EqualValues(t, 0x2306, ppu.Addr.Get())
	assert.EqualValues(t, 0x66, ppu.ReadData())
}

func TestPPU_VramRead_CrossPage(t *testing.T) {
	ppu, _ := stubPpu()
	ppu.Vram[0x1FF] = 0x66
	ppu.Vram[0x200] = 0x77
	ppu.WriteAddr(0x21)
	ppu.WriteAddr(0xFF)
	ppu.ReadData() // Buffer
	assert.EqualValues(t, 0x66, ppu.ReadData())
	assert.EqualValues(t, 0x77, ppu.ReadData())
}

func TestPPU_VramRead_Step32(t *testing.T) {
	ppu, _ := stubPpu()
	ppu.WriteCtrl(0b100)
	ppu.Vram[0x1FF] = 0x66
	ppu.Vram[0x1FF+32] = 0x77
	ppu.Vram[0x1FF+64] = 0x88
	ppu.WriteAddr(0x21)
	ppu.WriteAddr(0xFF)
	ppu.ReadData() // Buffer
	assert.EqualValues(t, 0x66, ppu.ReadData())
	assert.EqualValues(t, 0x77, ppu.ReadData())
	assert.EqualValues(t, 0x88, ppu.ReadData())
}

func TestPPU_HorizontalMirror(t *testing.T) {
	ppu, _ := stubPpu()
	ppu.WriteAddr(0x24)
	ppu.WriteAddr(0x05)
	ppu.WriteData(0x66) // A

	ppu.WriteAddr(0x28)
	ppu.WriteAddr(0x05)
	ppu.WriteData(0x77) // B

	ppu.WriteAddr(0x20)
	ppu.WriteAddr(0x05)
	ppu.ReadData() // Buffer

	assert.EqualValues(t, 0x66, ppu.ReadData()) // A

	ppu.WriteAddr(0x2C)
	ppu.WriteAddr(0x05)

	ppu.ReadData()
	assert.EqualValues(t, 0x77, ppu.ReadData()) // B
}

func TestPPU_VerticalMirror(t *testing.T) {
	ppu, cart := stubPpu()
	cart.Mirror = cartridge.Vertical
	ppu.WriteAddr(0x20)
	ppu.WriteAddr(0x05)
	ppu.WriteData(0x66) // A

	ppu.WriteAddr(0x2C)
	ppu.WriteAddr(0x05)
	ppu.WriteData(0x77) // B

	ppu.WriteAddr(0x28)
	ppu.WriteAddr(0x05)
	ppu.ReadData() // Buffer

	assert.EqualValues(t, 0x66, ppu.ReadData()) // A

	ppu.WriteAddr(0x24)
	ppu.WriteAddr(0x05)

	ppu.ReadData()
	assert.EqualValues(t, 0x77, ppu.ReadData()) // B
}

func TestPPU_StatusResetsLatch(t *testing.T) {
	ppu, _ := stubPpu()
	ppu.Vram[0x305] = 0x66
	ppu.WriteAddr(0x21)
	ppu.WriteAddr(0x23)
	ppu.WriteAddr(0x05)
	ppu.ReadData() // Buffer
	assert.NotEqualValues(t, 0x66, ppu.ReadData())

	ppu.ReadStatus()
	ppu.WriteAddr(0x23)
	ppu.WriteAddr(0x05)
	ppu.ReadData()
	assert.EqualValues(t, 0x66, ppu.ReadData())
}

func TestPPU_VramMirror(t *testing.T) {
	ppu, _ := stubPpu()
	ppu.Vram[0x305] = 0x66

	ppu.WriteAddr(0x63)
	ppu.WriteAddr(0x05)

	ppu.ReadData() // Buffer
	assert.EqualValues(t, 0x66, ppu.ReadData())
}

func TestPPU_StatusResetsVblank(t *testing.T) {
	ppu, _ := stubPpu()
	ppu.Status.Vblank = true
	assert.EqualValues(t, 1, ppu.ReadStatus()>>7)
	assert.EqualValues(t, 0, ppu.ReadStatus()>>7)
}

func TestCPU_OamReadWrite(t *testing.T) {
	ppu, _ := stubPpu()
	ppu.WriteOamAddr(0x10)
	ppu.WriteOam(0x66)
	ppu.WriteOam(0x77)
	ppu.WriteOamAddr(0x10)
	assert.EqualValues(t, 0x66, ppu.ReadOam())
	ppu.WriteOamAddr(0x11)
	assert.EqualValues(t, 0x77, ppu.ReadOam())
}

func TestCPU_OamDma(t *testing.T) {
	ppu, _ := stubPpu()

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
