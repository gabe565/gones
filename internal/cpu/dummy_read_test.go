package cpu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockBus struct {
	Memory  [0x10000]byte
	ReadLog []uint16
}

func (m *MockBus) ReadMem(addr uint16) byte {
	m.ReadLog = append(m.ReadLog, addr)
	return m.Memory[addr]
}

func (m *MockBus) ReadMemSafe(addr uint16) byte {
	return m.Memory[addr]
}

func (m *MockBus) WriteMem(addr uint16, data byte) {
	m.Memory[addr] = data
}

func (m *MockBus) WriteMem16(addr uint16, data uint16) {
	m.Memory[addr] = byte(data)
	m.Memory[addr+1] = byte(data >> 8)
}

func (m *MockBus) ReadMem16(addr uint16) uint16 {
	lo := uint16(m.ReadMem(addr))
	hi := uint16(m.ReadMem(addr + 1))
	return hi<<8 | lo
}

func stubCPUMockBus(program []byte) (*CPU, *MockBus) {
	bus := &MockBus{}
	// Load program at 0x8000 (Mapper 2 style but simplified)
	// We just write to memory directly.
	// Reset vector at FFFC.

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

func Test_DummyRead_AbsoluteX(t *testing.T) {
	// LDA $20F2, X (where X = $10)
	// Opcode: BD F2 20
	// Base: 20F2. X: 10.
	// Effective: 2102.
	// Page Crossed: Yes (20 -> 21).

	// Cycles:
	// 1. Read Opcode (PC)
	// 2. Read Low (PC+1) -> F2
	// 3. Read High (PC+2) -> 20. (Address built: 20F2)
	// 4. Dummy Read: (BaseHi << 8) | ((BaseLo + X) & FF)
	//    BaseHi = 20. BaseLo = F2. X = 10.
	//    BaseLo + X = 102. Low byte = 02.
	//    Dummy Addr = 2002.
	// 5. Read Effective: 2102.

	// Program: BD F2 20
	cpu, bus := stubCPUMockBus([]byte{0xBD, 0xF2, 0x20})
	cpu.RegisterX = 0x10

	bus.ReadLog = []uint16{} // Clear log from fetch/setup

	cpu.Step()

	// Expected Reads:
	// 1. 8000 (Opcode)
	// 2. 8001 (Low)
	// 3. 8002 (High)
	// 4. 2002 (Dummy Read - Page Crossed)
	// 5. 2102 (Effective Read)

	expectedReads := []uint16{0x8000, 0x8001, 0x8002, 0x2002, 0x2102}
	assert.Equal(t, expectedReads, bus.ReadLog, "Reads should match including dummy read")
}

func Test_DummyRead_AbsoluteX_NoCross(t *testing.T) {
	// LDA $2000, X (X=0)
	// BD 00 20
	// No page cross.
	// Reads: Op, Lo, Hi, Effective. No dummy.

	cpu, bus := stubCPUMockBus([]byte{0xBD, 0x00, 0x20})
	cpu.RegisterX = 0x00

	bus.ReadLog = []uint16{}

	cpu.Step()

	expectedReads := []uint16{0x8000, 0x8001, 0x8002, 0x2000}
	assert.Equal(t, expectedReads, bus.ReadLog, "Reads should match (no dummy)")
}
