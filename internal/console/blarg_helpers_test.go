package console

import (
	"bytes"
	"io"
	"regexp"
)

type Status int16

const (
	StatusPrerun  Status = -1
	StatusSuccess Status = 0
	StatusRunning Status = 0x80
	StatusReset   Status = 0x81
)

var BlarggTestMarker = [...]byte{222, 176, 97}

func NewBlarggTest(r io.ReadSeeker) (*ConsoleTest, error) {
	return NewConsoleTest(r, BlarggCallback)
}

func BlarggCallback(b *ConsoleTest) error {
	status := GetBlarggStatus(b)
	switch status {
	case StatusPrerun, StatusRunning:
		return nil
	case StatusReset:
		b.ResetIn = 6
		b.Console.Bus.WriteMem(0x6000, byte(StatusRunning))
	default:
		return ErrExit
	}
	return nil
}

func GetBlarggStatus(b *ConsoleTest) Status {
	status := Status(b.Console.Bus.ReadMem(0x6000))
	if status == 0 {
		var marker [3]byte
		for i := uint16(0); i < 3; i += 1 {
			marker[i] = b.Console.Bus.ReadMem(0x6001 + i)
		}
		if marker != BlarggTestMarker {
			return StatusPrerun
		}
	}
	return status
}

func GetBlarggMessage(b *ConsoleTest) string {
	var message []byte
	for i := uint16(0); ; i += 1 {
		data := b.Console.Bus.ReadMem(0x6004 + i)
		if data == 0 {
			break
		}
		message = append(message, data)
	}
	return string(bytes.TrimSpace(message))
}

type PpuMessage struct {
	Message string
}

func (p PpuMessage) Error() string {
	return p.Message
}

func NewBlargPpuMessageCallback() func(*ConsoleTest) error {
	var started bool
	re := regexp.MustCompile("  +")

	return func(b *ConsoleTest) error {
		if !started {
			if ready := b.Console.Bus.ReadMem(0x7F1); ready != 0 {
				started = true
			}
			return nil
		}

		if ready := b.Console.Bus.ReadMem(0x7F1); ready == 0 {
			var i int
			for {
				if b.Console.PPU.Vram[i] == 0 {
					break
				}
				i += 1
			}
			vram := b.Console.PPU.Vram[:i]
			vram = re.ReplaceAll(vram, []byte("\n"))
			vram = bytes.TrimSpace(vram)
			return PpuMessage{Message: string(vram)}
		}

		return nil
	}
}
