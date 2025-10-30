package msg

import "sync/atomic"

// Queue ring buffer single-producer / single-consumer
type Queue struct {
	// padding, avoid false sharing
	_pad0 [64]byte

	buf  []Message
	mask uint64 // = uint64(len(buf)-1)

	_pad1 [64]byte
	// head producer index
	head atomic.Uint64

	_pad2 [64]byte
	// tail consumer index
	tail atomic.Uint64

	_pad3 [64]byte
}

// NewQueue cap must be a power of two >=2.
func NewQueue(cap int) *Queue {
	if cap < 2 {
		cap = 2
	}

	n := 1
	for n < cap {
		n <<= 1
	}
	q := &Queue{
		buf:  make([]Message, n),
		mask: uint64(n - 1),
	}
	return q
}

func (q *Queue) Cap() int { return len(q.buf) }

func (q *Queue) Len() int {
	h := q.head.Load()
	t := q.tail.Load()
	return int(h - t)
}

// TryWrite writes a message, returning false if the queue is full.
func (q *Queue) TryWrite(m Message) bool {
	h := q.head.Load()
	t := q.tail.Load()
	if (h - t) >= uint64(len(q.buf)) {
		return false
	}
	q.buf[h&q.mask] = m
	q.head.Store(h + 1)
	return true
}

// WriteOverwrite writes a message, overwriting the oldest message if the queue is full.
func (q *Queue) WriteOverwrite(m Message) {
	h := q.head.Load()
	t := q.tail.Load()
	if (h - t) >= uint64(len(q.buf)) {
		q.tail.Store(t + 1)
	}
	q.buf[h&q.mask] = m
	q.head.Store(h + 1)
}

// TryRead reads a message, returning false if the queue is empty.
func (q *Queue) TryRead(out *Message) bool {
	t := q.tail.Load()
	h := q.head.Load()
	if t == h {
		return false // vide
	}
	*out = q.buf[t&q.mask]
	q.tail.Store(t + 1)
	return true
}

// Drain reads up to max messages, calling fn for each message.
func (q *Queue) Drain(max int, fn func(Message)) int {
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
