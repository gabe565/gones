package consts

const (
	// PRGROMAddr is the memory address that PRG begins.
	PRGROMAddr = 0x600

	// PRGChunkSize is the size of the smallest PRG
	PRGChunkSize = 0x4000

	// CHRChunkSize is the size of the smallest CHR
	CHRChunkSize = 0x2000

	CPUFrequency = 1789773.0

	AudioSampleRate = 44100.0

	FrameRate = 60.0

	RealFrameRate = 60.0988118623484

	FrameRateDifference = FrameRate / RealFrameRate
)
