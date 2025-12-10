package cpu

import (
	"testing"

	"gabe565.com/gones/internal/apu"
	"gabe565.com/gones/internal/bus"
	"gabe565.com/gones/internal/cartridge"
	"gabe565.com/gones/internal/config"
	"gabe565.com/gones/internal/ppu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type WriteLog struct {
	Addr uint16
	Data byte
}

type MockBus struct {
	Memory   [0x10000]byte
	ReadLog  []uint16
	WriteLog []WriteLog
}

func (m *MockBus) ReadMem(addr uint16) byte {
	m.ReadLog = append(m.ReadLog, addr)
	return m.Memory[addr]
}

func (m *MockBus) ReadMemSafe(addr uint16) byte {
	return m.Memory[addr]
}

func (m *MockBus) WriteMem(addr uint16, data byte) {
	m.WriteLog = append(m.WriteLog, WriteLog{addr, data})
	m.Memory[addr] = data
}

func (m *MockBus) WriteMem16(addr uint16, data uint16) {
	m.WriteMem(addr, byte(data))
	m.WriteMem(addr+1, byte(data>>8))
}

func (m *MockBus) ReadMem16(addr uint16) uint16 {
	lo := uint16(m.ReadMem(addr))
	hi := uint16(m.ReadMem(addr + 1))
	return hi<<8 | lo
}

func stubCPU(program []byte) *CPU {
	cart := cartridge.FromBytes(program)
	mapper := cartridge.NewMapper2(cart, false)
	ppu := ppu.New(config.NewDefault(), mapper)
	conf := config.NewDefault()
	apu := apu.New(conf)
	bus := bus.New(conf, mapper, ppu, apu)
	cpu := New(bus)
	apu.SetCPU(cpu)
	return cpu
}

func stubCPUMockBus(program []byte) (*CPU, *MockBus) {
	bus := &MockBus{}
	for i, b := range program {
		bus.Memory[0x8000+i] = b
	}
	bus.Memory[0xFFFC] = 0x00
	bus.Memory[0xFFFD] = 0x80

	cpu := &CPU{
		StackPointer:   0xFD,
		Status:         Status{InterruptDisable: true},
		bus:            bus,
		Cycles:         7,
		ProgramCounter: 0x8000,
	}
	return cpu, bus
}

func Test_0xa9_lda_immediate_load_data(t *testing.T) {
	t.Parallel()

	cpu := stubCPU([]byte{0xA9, 0x05, 0x00})
	for {
		cpu.Step()
		require.NoError(t, cpu.StepErr)
		if cpu.Status.Break {
			break
		}
	}

	assert.EqualValues(t, 0x05, cpu.Accumulator)
	assert.EqualValues(t, 0, cpu.Status.Get()&0b0000_0010)
	assert.EqualValues(t, 0, cpu.Status.Get()&0b1000_0010)
}

func Test_0xa9_lda_zero_flag(t *testing.T) {
	t.Parallel()

	cpu := stubCPU([]byte{0xA9, 0x00, 0x00})
	for {
		cpu.Step()
		require.NoError(t, cpu.StepErr)
		if cpu.Status.Break {
			break
		}
	}

	assert.EqualValues(t, 0b10, cpu.Status.Get()&0b0000_0010)
}

func Test_0xaa_tax_move_a_to_x(t *testing.T) {
	t.Parallel()

	cpu := stubCPU([]byte{0xA9, 0x0A, 0xAA, 0x00})
	for {
		cpu.Step()
		require.NoError(t, cpu.StepErr)
		if cpu.Status.Break {
			break
		}
	}

	assert.EqualValues(t, 10, cpu.RegisterX)
}

func Test_5_operations(t *testing.T) {
	t.Parallel()

	cpu := stubCPU([]byte{0xA9, 0xC0, 0xAA, 0xE8, 0x00})
	for {
		cpu.Step()
		require.NoError(t, cpu.StepErr)
		if cpu.Status.Break {
			break
		}
	}

	assert.EqualValues(t, 0xC1, cpu.RegisterX)
}

func Test_inx_overflow(t *testing.T) {
	t.Parallel()

	cpu := stubCPU([]byte{0xA9, 0xFF, 0xAA, 0xE8, 0xE8, 0x00})
	for {
		cpu.Step()
		require.NoError(t, cpu.StepErr)
		if cpu.Status.Break {
			break
		}
	}

	assert.EqualValues(t, 1, cpu.RegisterX)
}

func Test_lda_from_memory(t *testing.T) {
	t.Parallel()

	cpu := stubCPU([]byte{0xA5, 0x10, 0x00})
	cpu.WriteMem(0x10, 0x55)
	for {
		cpu.Step()
		require.NoError(t, cpu.StepErr)
		if cpu.Status.Break {
			break
		}
	}

	assert.EqualValues(t, 0x55, cpu.Accumulator)
}
