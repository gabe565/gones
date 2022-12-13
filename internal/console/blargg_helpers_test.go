package console

import (
	"github.com/gabe565/gones/internal/memory"
)

type Status int16

const (
	StatusPrerun  Status = -1
	StatusSuccess Status = 0
	StatusRunning Status = 0x80
	StatusReset   Status = 0x81
)

var TestMarker = [...]byte{222, 176, 97}

func getBlarggStatus(bus memory.Read8) Status {
	status := Status(bus.ReadMem(0x6000))
	if status == 0 {
		var marker [3]byte
		for i := uint16(0); i < 3; i += 1 {
			marker[i] = bus.ReadMem(0x6001 + i)
		}
		if marker != TestMarker {
			return StatusPrerun
		}
	}
	return status
}

func getBlarggMessage(bus memory.Read8) string {
	var message []byte
	for i := uint16(0); ; i += 1 {
		data := bus.ReadMem(0x6004 + i)
		if data == 0 {
			break
		}
		message = append(message, data)
	}
	return string(message)
}

func runBlarggTest(console *Console) (Status, error) {
	var resetDelay uint8

	for {
		if resetDelay != 0 {
			resetDelay -= 1
			if resetDelay == 0 {
				console.Reset()
			}
		}

		if err := console.Step(); err != nil {
			return 0, err
		}

		status := getBlarggStatus(console.Bus)
		switch status {
		case StatusPrerun, StatusRunning:
			continue
		case StatusReset:
			resetDelay = 6
			console.Bus.WriteMem(0x6000, byte(StatusRunning))
		default:
			return status, nil
		}
	}
}
