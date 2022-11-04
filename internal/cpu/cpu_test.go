package cpu

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_0xa9_lda_immediate_load_data(t *testing.T) {
	cpu := New()
	cpu.loadAndRun([]uint8{0xa9, 0x05, 0x00})

	assert.EqualValues(t, 0x05, cpu.RegisterA)
	assert.EqualValues(t, 0, cpu.Acc&0b0000_0010)
	assert.EqualValues(t, 0, cpu.Acc&0b1000_0010)
}

func Test_0xa9_lda_zero_flag(t *testing.T) {
	cpu := New()
	cpu.loadAndRun([]uint8{0xa9, 0x00, 0x00})

	assert.EqualValues(t, 0b10, cpu.Acc&0b0000_0010)
}

func Test_0xaa_tax_move_a_to_x(t *testing.T) {
	cpu := New()
	cpu.RegisterA = 10
	cpu.loadAndRun([]uint8{0xa9, 0x0a, 0xaa, 0x00})

	assert.EqualValues(t, 10, cpu.RegisterX)
}

func Test_5_operations(t *testing.T) {
	cpu := New()
	cpu.loadAndRun([]uint8{0xa9, 0xc0, 0xaa, 0xe8, 0x00})

	assert.EqualValues(t, 0xc1, cpu.RegisterX)
}

func Test_inx_overflow(t *testing.T) {
	cpu := New()
	cpu.loadAndRun([]uint8{0xa9, 0xff, 0xaa, 0xe8, 0xe8, 0x00})

	assert.EqualValues(t, 1, cpu.RegisterX)
}
