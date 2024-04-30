package util

import (
	"fmt"
)

func EncodeHexAddr(i uint16) string {
	return fmt.Sprintf("$%04X", i)
}

func EncodeHexVal(i uint8) string {
	return fmt.Sprintf("%02X", i)
}
