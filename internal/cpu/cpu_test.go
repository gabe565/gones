package cpu

import (
	"github.com/gabe565/gones/internal/bus"
	"github.com/gabe565/gones/internal/cartridge"
	"github.com/gabe565/gones/internal/ppu"
	"github.com/stretchr/testify/assert"
	"testing"
)

func stubCpu(program []byte) *CPU {
	cart := cartridge.FromBytes(program)
	mapper := cartridge.NewMapper2(cart)
	ppu := ppu.New(cart, mapper)
	bus := bus.New(mapper, ppu)
	cpu := New(bus)
	cpu.Reset()
	return cpu
}

func Test_0xa9_lda_immediate_load_data(t *testing.T) {
	cpu := stubCpu([]byte{0xA9, 0x05, 0x00})
	for {
		if _, err := cpu.Step(); !assert.NoError(t, err) {
			return
		}
		if cpu.Status.Has(Break) {
			break
		}
	}

	assert.EqualValues(t, 0x05, cpu.Accumulator)
	assert.EqualValues(t, 0, cpu.Status&0b0000_0010)
	assert.EqualValues(t, 0, cpu.Status&0b1000_0010)
}

func Test_0xa9_lda_zero_flag(t *testing.T) {
	cpu := stubCpu([]byte{0xA9, 0x00, 0x00})
	for {
		if _, err := cpu.Step(); !assert.NoError(t, err) {
			return
		}
		if cpu.Status.Has(Break) {
			break
		}
	}

	assert.EqualValues(t, 0b10, cpu.Status&0b0000_0010)
}

func Test_0xaa_tax_move_a_to_x(t *testing.T) {
	cpu := stubCpu([]byte{0xA9, 0x0A, 0xAA, 0x00})
	for {
		if _, err := cpu.Step(); !assert.NoError(t, err) {
			return
		}
		if cpu.Status.Has(Break) {
			break
		}
	}

	assert.EqualValues(t, 10, cpu.RegisterX)
}

func Test_5_operations(t *testing.T) {
	cpu := stubCpu([]byte{0xA9, 0xC0, 0xAA, 0xE8, 0x00})
	for {
		if _, err := cpu.Step(); !assert.NoError(t, err) {
			return
		}
		if cpu.Status.Has(Break) {
			break
		}
	}

	assert.EqualValues(t, 0xC1, cpu.RegisterX)
}

func Test_inx_overflow(t *testing.T) {
	cpu := stubCpu([]byte{0xA9, 0xFF, 0xAA, 0xE8, 0xE8, 0x00})
	for {
		if _, err := cpu.Step(); !assert.NoError(t, err) {
			return
		}
		if cpu.Status.Has(Break) {
			break
		}
	}

	assert.EqualValues(t, 1, cpu.RegisterX)
}

func Test_lda_from_memory(t *testing.T) {
	cpu := stubCpu([]byte{0xA5, 0x10, 0x00})
	cpu.MemWrite(0x10, 0x55)
	for {
		if _, err := cpu.Step(); !assert.NoError(t, err) {
			return
		}
		if cpu.Status.Has(Break) {
			break
		}
	}

	assert.EqualValues(t, 0x55, cpu.Accumulator)
}
