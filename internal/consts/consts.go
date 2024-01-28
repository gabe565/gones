package consts

const (
	// PrgRomAddr is the memory address that PRG begins.
	PrgRomAddr = 0x600

	// PrgChunkSize is the size of the smallest PRG
	PrgChunkSize = 0x4000

	// ChrChunkSize is the size of the smallest CHR
	ChrChunkSize = 0x2000

	CpuFrequency = 1789773.0

	AudioSampleRate = 44100.0

	FrameRate = 60.0

	RealFrameRate = 60.0988118623484

	FrameRateDifference = FrameRate / RealFrameRate
)
