package consts

import "time"

const (
	// PRGROMAddr is the memory address that PRG begins.
	PRGROMAddr = 0x600

	// PRGChunkSize is the size of the smallest PRG
	PRGChunkSize = 0x4000

	// CHRChunkSize is the size of the smallest CHR
	CHRChunkSize = 0x2000

	CPUFrequency = 1789773

	AudioSampleRate     = 44100
	AudioBufferSize     = time.Second / 20
	AudioBytesPerSample = 4 * 2
	AudioBufferBytes    = int(AudioBufferSize * AudioBytesPerSample * AudioSampleRate / time.Second)

	TargetFrameRate   = 60
	HardwareFrameRate = 60.0988118623484
	FrameRateDiff     = TargetFrameRate / HardwareFrameRate

	Width  = 256
	Height = 240

	PPUOAMSize     = 256
	PPUSpriteLimit = 8
)
