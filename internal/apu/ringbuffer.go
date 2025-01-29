package apu

import (
	"sync"
)

type ringBuffer struct {
	buf  []byte
	size int
	r, w int
	full bool
	mu   sync.Mutex
}

func newRingBuffer(size int) *ringBuffer {
	return &ringBuffer{
		buf:  make([]byte, size),
		size: size,
	}
}

func (r *ringBuffer) Read(p []byte) int {
	r.mu.Lock()
	defer r.mu.Unlock()

	var n int
	switch {
	case r.w == r.r && !r.full:
		// Buffer empty
		return 0
	case r.w > r.r:
		// Writer ahead of reader
		if n = r.w - r.r; n > len(p) {
			n = len(p)
		}
		copy(p, r.buf[r.r:r.r+n])
	default:
		// Reader ahead of writer
		if n = r.size - r.r + r.w; n > len(p) {
			n = len(p)
		}

		if r.r+n <= r.size {
			// End of buffer has enough elements
			copy(p, r.buf[r.r:r.r+n])
		} else {
			// End of buffer does not have enough elements; read will wrap around
			c1 := r.size - r.r
			copy(p, r.buf[r.r:r.size])
			c2 := n - c1
			copy(p[c1:], r.buf[0:c2])
		}

		r.full = false
	}
	r.r = (r.r + n) % r.size
	return n
}

func (r *ringBuffer) Write(p []byte) {
	r.mu.Lock()
	defer r.mu.Unlock()

	n := len(p)
	switch {
	case r.free() < len(p):
		// Buffer too full; discard write
		return
	case r.w >= r.r:
		// Writer ahead of reader
		if c1 := r.size - r.w; c1 >= n {
			// Slice fits in the end of the buffer
			copy(r.buf[r.w:], p)
			r.w += n
		} else {
			// Slice does not fit in the end of the buffer; write will wrap around
			copy(r.buf[r.w:], p[:c1])
			c2 := n - c1
			copy(r.buf[0:], p[c1:])
			r.w = c2
		}
	default:
		// Reader is ahead of writer
		copy(r.buf[r.w:], p)
		r.w += n
	}

	if r.w == r.size {
		r.w = 0
	}
	r.full = r.w == r.r
}

func (r *ringBuffer) len() int {
	r.mu.Lock()
	defer r.mu.Unlock()

	switch {
	case r.w == r.r:
		if r.full {
			return r.size
		}
		return 0
	case r.w > r.r:
		return r.w - r.r
	default:
		return r.size - r.r + r.w
	}
}

func (r *ringBuffer) free() int {
	switch {
	case r.w == r.r:
		if r.full {
			return 0
		}
		return r.size
	case r.w < r.r:
		return r.r - r.w
	default:
		return r.size - r.w + r.r
	}
}

func (r *ringBuffer) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.r = 0
	r.w = 0
	r.full = false
}
