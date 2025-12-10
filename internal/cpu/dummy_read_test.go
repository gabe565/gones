package cpu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
	// LDA $2000, X (X=0) -> BD 00 20. No page cross. No dummy (Load).
	// STA $2000, X (X=0) -> 9D 00 20. No page cross. Dummy read required (Store).

	t.Run("LDA No Cross", func(t *testing.T) {
		cpu, bus := stubCPUMockBus([]byte{0xBD, 0x00, 0x20})
		cpu.RegisterX = 0x00
		bus.ReadLog = []uint16{}
		cpu.Step()
		expectedReads := []uint16{0x8000, 0x8001, 0x8002, 0x2000}
		assert.Equal(t, expectedReads, bus.ReadLog)
	})

	t.Run("STA No Cross (Dummy Read)", func(t *testing.T) {
		// STA $2000, X
		// 1. Fetch Op (8000)
		// 2. Fetch Lo (8001)
		// 3. Fetch Hi (8002)
		// 4. Dummy Read (2000)
		// 5. Write (2000) - Not logged in ReadLog

		cpu, bus := stubCPUMockBus([]byte{0x9D, 0x00, 0x20})
		cpu.RegisterX = 0x00
		bus.ReadLog = []uint16{}
		cpu.Step()

		expectedReads := []uint16{0x8000, 0x8001, 0x8002, 0x2000}
		assert.Equal(t, expectedReads, bus.ReadLog)
	})
}
