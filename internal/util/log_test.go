package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHexAddr_String(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		h    HexAddr
		want string
	}{
		{"0", HexAddr(0), "$0000"},
		{"10", HexAddr(10), "$000A"},
		{"100", HexAddr(100), "$0064"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, tt.h.String())
		})
	}
}

func TestHexVal_String(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		h    HexVal
		want string
	}{
		{"0", HexVal(0), "00"},
		{"10", HexVal(10), "0A"},
		{"100", HexVal(100), "64"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, tt.h.String())
		})
	}
}
