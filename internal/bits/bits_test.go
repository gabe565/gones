package bits

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBits_Set(t *testing.T) {
	b := Bits(0)
	b.Set(0b10)
	assert.EqualValues(t, 0b10, b)
}

func TestBits_Clear(t *testing.T) {
	b := Bits(0b110)
	b.Clear(0b10)
	assert.EqualValues(t, 0b100, b)
}

func TestBits_Toggle(t *testing.T) {
	b := Bits(0b10)
	b.Toggle(0b1)
	assert.EqualValues(t, 0b11, b)
	b.Toggle(0b10)
	assert.EqualValues(t, 0b1, b)
}

func TestBits_Has(t *testing.T) {
	b := Bits(0b10)
	assert.Equal(t, true, b.Has(0b10))
	assert.Equal(t, false, b.Has(0b1))
}
