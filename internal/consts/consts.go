package consts

const (
	// PRGROMAddr is the memory address that PRG begins.
	PRGROMAddr = 0x600

	// PRGChunkSize is the size of the smallest PRG
	PRGChunkSize = 0x4000

	// CHRChunkSize is the size of the smallest CHR
	CHRChunkSize = 0x2000

	CPUFrequency = 1789773

	AudioSampleRate = 44100

	FrameRate = 60

	RealFrameRate = 60.0988118623484

	FrameRateDifference = FrameRate / RealFrameRate

	Width  = 256
	Height = 240
)
