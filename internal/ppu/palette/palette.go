package palette

import (
	"image/color"
)

//go:generate go run ./generate/main.go

var Default = [0x40]color.RGBA{
	0x00: {0x75, 0x75, 0x75, 0xFF},
	0x01: {0x24, 0x18, 0x8E, 0xFF},
	0x02: {0x00, 0x00, 0xAA, 0xFF},
	0x03: {0x45, 0x00, 0x9E, 0xFF},
	0x04: {0x8E, 0x00, 0x75, 0xFF},
	0x05: {0xAA, 0x00, 0x10, 0xFF},
	0x06: {0xA6, 0x00, 0x00, 0xFF},
	0x07: {0x7D, 0x08, 0x00, 0xFF},
	0x08: {0x41, 0x2C, 0x00, 0xFF},
	0x09: {0x00, 0x45, 0x00, 0xFF},
	0x0A: {0x00, 0x51, 0x00, 0xFF},
	0x0B: {0x00, 0x3C, 0x14, 0xFF},
	0x0C: {0x18, 0x3C, 0x5D, 0xFF},
	0x0D: {0x00, 0x00, 0x00, 0xFF},
	0x0E: {0x00, 0x00, 0x00, 0xFF},
	0x0F: {0x00, 0x00, 0x00, 0xFF},
	0x10: {0xBE, 0xBE, 0xBE, 0xFF},
	0x11: {0x00, 0x71, 0xEF, 0xFF},
	0x12: {0x20, 0x38, 0xEF, 0xFF},
	0x13: {0x82, 0x00, 0xF3, 0xFF},
	0x14: {0xBE, 0x00, 0xBE, 0xFF},
	0x15: {0xE7, 0x00, 0x59, 0xFF},
	0x16: {0xDB, 0x28, 0x00, 0xFF},
	0x17: {0xCB, 0x4D, 0x0C, 0xFF},
	0x18: {0x8A, 0x71, 0x00, 0xFF},
	0x19: {0x00, 0x96, 0x00, 0xFF},
	0x1A: {0x00, 0xAA, 0x00, 0xFF},
	0x1B: {0x00, 0x92, 0x38, 0xFF},
	0x1C: {0x00, 0x82, 0x8A, 0xFF},
	0x1D: {0x00, 0x00, 0x00, 0xFF},
	0x1E: {0x00, 0x00, 0x00, 0xFF},
	0x1F: {0x00, 0x00, 0x00, 0xFF},
	0x20: {0xFF, 0xFF, 0xFF, 0xFF},
	0x21: {0x3C, 0xBE, 0xFF, 0xFF},
	0x22: {0x5D, 0x96, 0xFF, 0xFF},
	0x23: {0xCF, 0x8A, 0xFF, 0xFF},
	0x24: {0xF7, 0x79, 0xFF, 0xFF},
	0x25: {0xFF, 0x75, 0xB6, 0xFF},
	0x26: {0xFF, 0x75, 0x61, 0xFF},
	0x27: {0xFF, 0x9A, 0x38, 0xFF},
	0x28: {0xF3, 0xBE, 0x3C, 0xFF},
	0x29: {0x82, 0xD3, 0x10, 0xFF},
	0x2A: {0x4D, 0xDF, 0x49, 0xFF},
	0x2B: {0x59, 0xFB, 0x9A, 0xFF},
	0x2C: {0x00, 0xEB, 0xDB, 0xFF},
	0x2D: {0x79, 0x79, 0x79, 0xFF},
	0x2E: {0x00, 0x00, 0x00, 0xFF},
	0x2F: {0x00, 0x00, 0x00, 0xFF},
	0x30: {0xFF, 0xFF, 0xFF, 0xFF},
	0x31: {0xAA, 0xE7, 0xFF, 0xFF},
	0x32: {0xC7, 0xD7, 0xFF, 0xFF},
	0x33: {0xD7, 0xCB, 0xFF, 0xFF},
	0x34: {0xFF, 0xC7, 0xFF, 0xFF},
	0x35: {0xFF, 0xC7, 0xDB, 0xFF},
	0x36: {0xFF, 0xBE, 0xB2, 0xFF},
	0x37: {0xFF, 0xDB, 0xAA, 0xFF},
	0x38: {0xFF, 0xE7, 0xA2, 0xFF},
	0x39: {0xE3, 0xFF, 0xA2, 0xFF},
	0x3A: {0xAA, 0xF3, 0xBE, 0xFF},
	0x3B: {0xB2, 0xFF, 0xCF, 0xFF},
	0x3C: {0x9E, 0xFF, 0xF3, 0xFF},
	0x3D: {0xC7, 0xC7, 0xC7, 0xFF},
	0x3E: {0x00, 0x00, 0x00, 0xFF},
	0x3F: {0x00, 0x00, 0x00, 0xFF},
}
