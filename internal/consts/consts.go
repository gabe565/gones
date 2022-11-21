package consts

const (
	// PrgRomAddr is the memory address that PRG begins.
	PrgRomAddr = 0x600

	// PrgChunkSize is the size of the smallest PRG
	PrgChunkSize = 0x4000

	// ChrChunkSize is the size of the smallest CHR
	ChrChunkSize = 0x2000

	// ResetAddr is the memory address for the Reset Interrupt Vector.
	ResetAddr = 0xFFFC

	CpuFrequency = 1789772

	AudioSampleRate = 48000
)
