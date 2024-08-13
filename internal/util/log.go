package util

import (
	"fmt"
)

type HexAddr uint16

func (h HexAddr) String() string {
	return fmt.Sprintf("$%04X", uint8(h))
}

type HexVal uint8

func (h HexVal) String() string {
	return fmt.Sprintf("%02X", uint8(h))
}
