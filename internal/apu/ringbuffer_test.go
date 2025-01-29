package apu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ringBuffer(t *testing.T) {
	const bufSize = 16

	buf := newRingBuffer(16)
	p := make([]byte, bufSize)
	n := buf.Read(p)
	assert.Equal(t, 0, n)
	assert.Equal(t, p, make([]byte, bufSize))
	assert.Equal(t, bufSize, buf.free())
	assert.Equal(t, 0, buf.len())

	const message = "hello world"
	for range 3 {
		buf.Write([]byte(message))
		assert.Equal(t, bufSize-len(message), buf.free())
		assert.Equal(t, len(message), buf.len())

		n = buf.Read(p)
		assert.Equal(t, len(message), n)
		assert.Equal(t, []byte(message), p[:n])
		assert.Equal(t, bufSize, buf.free())
		assert.Equal(t, 0, buf.len())
	}

	buf.Write([]byte(message))
	p = make([]byte, 5)
	n = buf.Read(p)
	assert.Equal(t, bufSize-len(message)+5, buf.free())
	assert.Equal(t, len(message)-5, buf.len())
	assert.Equal(t, 5, n)
	assert.Equal(t, []byte("hello"), p[:n])
}
