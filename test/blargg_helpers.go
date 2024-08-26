package test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"regexp"

	"github.com/gabe565/gones/internal/console"
	"github.com/gabe565/gones/internal/consts"
)

type status int16

const (
	statusPreRun  status = -1
	statusSuccess status = 0
	statusRunning status = 0x80
	statusReset   status = 0x81
)

type msgType uint8

const (
	msgTypeSRAM msgType = iota
	msgTypePPUVRAM
)

var errInvalidMessageType = errors.New("invalid message type")

func newBlarggTest(r io.ReadSeeker, t msgType) (*consoleTest, error) {
	switch t {
	case msgTypeSRAM:
		return newConsoleTest(r, blarggSRAMMsgCallback)
	case msgTypePPUVRAM:
		return newConsoleTest(r, blarggPPUMsgCallback())
	}
	return nil, fmt.Errorf("%w: %d", errInvalidMessageType, t)
}

func blarggSRAMMsgCallback(c *consoleTest) error {
	status := getBlarggStatus(c)
	switch status {
	case statusPreRun, statusRunning:
		return nil
	case statusReset:
		if c.resetIn == 0 {
			c.resetIn = consts.CPUFrequency / 10
		}
	default:
		return console.ErrExit
	}
	return nil
}

func getBlarggStatus(c *consoleTest) status {
	status := status(c.console.Bus.ReadMem(0x6000))
	if status == 0 {
		for i, b := range [3]byte{0xDE, 0xB0, 0x61} {
			if got := c.console.Bus.ReadMem(0x6001 + uint16(i)); got != b {
				return statusPreRun
			}
		}
	}
	return status
}

func getBlarggMessage(c *consoleTest, t msgType) string {
	var msg []byte
	switch t {
	case msgTypeSRAM:
		msg = c.console.Cartridge.SRAM[4:]
	case msgTypePPUVRAM:
		msg = c.console.PPU.VRAM[:]
	}
	msg, _, found := bytes.Cut(msg, []byte{0})
	if !found {
		return ""
	}
	if t == msgTypePPUVRAM {
		msg = regexp.MustCompile("  +").ReplaceAll(msg, []byte("\n"))
	}
	return string(bytes.TrimSpace(msg))
}

func blarggPPUMsgCallback() func(*consoleTest) error {
	var started bool

	return func(c *consoleTest) error {
		if !started {
			if ready := c.console.Bus.ReadMem(0x7F1); ready != 0 {
				started = true
			}
			return nil
		}

		if ready := c.console.Bus.ReadMem(0x7F1); ready == 0 {
			return console.ErrExit
		}
		return nil
	}
}
