package dsp

type Noise struct {
	rng       uint32
	buf       [BlockSize]float32
	stampedAt uint64
}

func NewNoise() *Noise {
	return &Noise{
		rng: 0x9E3779B9, // arbitrary non-zero seed
	}
}

func (n *Noise) Process(block *Block) {
	for i := 0; i < BlockSize; i++ {
		x := n.xorShift32()
		u := float32(x) * (1.0 / 4294967296.0)
		block.L[i] = 2*u - 1
		block.R[i] = block.L[i]
	}
}

func (n *Noise) Reset(bool) {
	// no-op
}

func (n *Noise) xorShift32() uint32 {
	x := n.rng
	x ^= x << 13
	x ^= x >> 17
	x ^= x << 5
	n.rng = x

	return x
}
