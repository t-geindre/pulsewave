package ui

import (
	"sync/atomic"
	"synth/dsp"
)

type AudioQueue struct {
	// padding, avoid false sharing
	_pad0 [64]byte

	buf  []dsp.Block
	mask uint64 // = uint64(len(buf)-1)

	_pad1 [64]byte
	// head producer index
	head atomic.Uint64

	_pad2 [64]byte
	// tail consumer index
	tail atomic.Uint64

	_pad3 [64]byte
}

func NewAudioQueue(cap int) *AudioQueue {
	if cap < 2 {
		cap = 2
	}

	n := 1
	for n < cap {
		n <<= 1
	}
	q := &AudioQueue{
		buf:  make([]dsp.Block, n),
		mask: uint64(n - 1),
	}
	return q
}

func (q *AudioQueue) Cap() int { return len(q.buf) }

func (q *AudioQueue) Len() int {
	h := q.head.Load()
	t := q.tail.Load()
	return int(h - t)
}

func (q *AudioQueue) TryWrite(m dsp.Block) bool {
	h := q.head.Load()
	t := q.tail.Load()
	if (h - t) >= uint64(len(q.buf)) {
		return false
	}
	q.buf[h&q.mask] = m
	q.head.Store(h + 1)
	return true
}

func (q *AudioQueue) WriteOverwrite(m dsp.Block) {
	h := q.head.Load()
	t := q.tail.Load()
	if (h - t) >= uint64(len(q.buf)) {
		q.tail.Store(t + 1)
	}
	q.buf[h&q.mask] = m
	q.head.Store(h + 1)
}

func (q *AudioQueue) TryRead(out *dsp.Block) bool {
	t := q.tail.Load()
	h := q.head.Load()
	if t == h {
		return false // vide
	}
	*out = q.buf[t&q.mask]
	q.tail.Store(t + 1)
	return true
}

func (q *AudioQueue) Drain(max int, fn func(dsp.Block)) int {
	t := q.tail.Load()
	h := q.head.Load()
	n := int(h - t)
	if n <= 0 {
		return 0
	}
	if max > 0 && n > max {
		n = max
	}
	for i := 0; i < n; i++ {
		msg := q.buf[(t+uint64(i))&q.mask]
		fn(msg)
	}
	q.tail.Store(t + uint64(n))
	return n
}
