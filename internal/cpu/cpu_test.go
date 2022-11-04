package cpu

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_0xa9_lda_immediate_load_data(t *testing.T) {
	cpu := New()
	if err := cpu.loadAndRun([]uint8{0xA9, 0x05, 0x00}); err != nil {
		assert.NoErrorf(t, err, "loadAndRun")
	}

	assert.EqualValues(t, 0x05, cpu.RegisterA)
	assert.EqualValues(t, 0, cpu.Acc&0b0000_0010)
	assert.EqualValues(t, 0, cpu.Acc&0b1000_0010)
}

func Test_0xa9_lda_zero_flag(t *testing.T) {
	cpu := New()
	if err := cpu.loadAndRun([]uint8{0xA9, 0x00, 0x00}); err != nil {
		assert.NoErrorf(t, err, "loadAndRun")
	}

	assert.EqualValues(t, 0b10, cpu.Acc&0b0000_0010)
}

func Test_0xaa_tax_move_a_to_x(t *testing.T) {
	cpu := New()
	cpu.RegisterA = 10
	if err := cpu.loadAndRun([]uint8{0xA9, 0x0A, 0xAA, 0x00}); err != nil {
		assert.NoErrorf(t, err, "loadAndRun")
	}

	assert.EqualValues(t, 10, cpu.RegisterX)
}

func Test_5_operations(t *testing.T) {
	cpu := New()
	if err := cpu.loadAndRun([]uint8{0xA9, 0xC0, 0xAA, 0xE8, 0x00}); err != nil {
		assert.NoErrorf(t, err, "loadAndRun")
	}

	assert.EqualValues(t, 0xC1, cpu.RegisterX)
}

func Test_inx_overflow(t *testing.T) {
	cpu := New()
	if err := cpu.loadAndRun([]uint8{0xA9, 0xFF, 0xAA, 0xE8, 0xE8, 0x00}); err != nil {
		assert.NoErrorf(t, err, "loadAndRun")
	}

	assert.EqualValues(t, 1, cpu.RegisterX)
}

func Test_lda_from_memory(t *testing.T) {
	cpu := New()
	cpu.memWrite(0x10, 0x55)
	if err := cpu.loadAndRun([]uint8{0xA5, 0x10, 0x00}); err != nil {
		assert.NoErrorf(t, err, "loadAndRun")
	}

	assert.EqualValues(t, 0x55, cpu.RegisterA)
}
