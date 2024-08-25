package test

import (
	"bytes"
	"io"
	"regexp"

	"github.com/gabe565/gones/internal/console"
)

type status int16

const (
	statusPreRun  status = -1
	statusSuccess status = 0
	statusRunning status = 0x80
	statusReset   status = 0x81
)

func newBlarggTest(r io.ReadSeeker) (*consoleTest, error) {
	return newConsoleTest(r, blarggCallback)
}

func blarggCallback(c *consoleTest) error {
	status := getBlarggStatus(c)
	switch status {
	case statusPreRun, statusRunning:
		return nil
	case statusReset:
		c.resetIn = 6
		c.console.Bus.WriteMem(0x6000, byte(statusRunning))
	default:
		return console.ErrExit
	}
	return nil
}

func getBlarggStatus(c *consoleTest) status {
	status := status(c.console.Bus.ReadMem(0x6000))
	if status == 0 {
		var marker [3]byte
		for i := range uint16(3) {
			marker[i] = c.console.Bus.ReadMem(0x6001 + i)
		}
		if marker != [3]byte{222, 176, 97} {
			return statusPreRun
		}
	}
	return status
}

func getBlarggMessage(c *consoleTest) string {
	var message []byte
	var i uint16
	for {
		data := c.console.Bus.ReadMem(0x6004 + i)
		if data == 0 {
			break
		}
		message = append(message, data)
		i++
	}
	return string(bytes.TrimSpace(message))
}

type ppuMessageError string

func (p ppuMessageError) Error() string {
	return string(p)
}

func newBlarggPPUMessageCallback() func(*consoleTest) error {
	var started bool
	re := regexp.MustCompile("  +")

	return func(c *consoleTest) error {
		if !started {
			if ready := c.console.Bus.ReadMem(0x7F1); ready != 0 {
				started = true
			}
			return nil
		}

		if ready := c.console.Bus.ReadMem(0x7F1); ready == 0 {
			var i int
			for {
				if c.console.PPU.VRAM[i] == 0 {
					break
				}
				i++
			}
			vram := c.console.PPU.VRAM[:i]
			vram = re.ReplaceAll(vram, []byte("\n"))
			vram = bytes.TrimSpace(vram)
			return ppuMessageError(vram)
		}

		return nil
	}
}
