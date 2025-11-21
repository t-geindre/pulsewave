package dsp

const (
	NoiseWhite float32 = iota
	NoiseGaussian
	NoisePink
	NoiseBrown
	NoiseBlue
)

type Noise struct {
	rng uint32

	noiseType Param

	lastWhite float32
	brown     float32

	pinkRows    [16]float32
	pinkAccum   float32
	pinkCounter uint32
}

func NewNoise(noiseType Param) *Noise {
	return &Noise{
		rng:       0x9E3779B9,
		noiseType: noiseType,
	}
}

func (n *Noise) Reset(soft bool) {
	if !soft {
		n.lastWhite = 0
		n.brown = 0
		n.pinkAccum = 0
		n.pinkCounter = 0

		for i := range n.pinkRows {
			n.pinkRows[i] = 0
		}

		n.rng = 0x9E3779B9
	}
}

func (n *Noise) Process(block *Block) {
	switch n.noiseType.Resolve(block.Cycle)[0] {
	case NoiseWhite:
		n.processWhite(block)
	case NoiseGaussian:
		n.processGaussian(block)
	case NoisePink:
		n.processPink(block)
	case NoiseBrown:
		n.processBrown(block)
	case NoiseBlue:
		n.processBlue(block)
	default:
		n.processWhite(block)
	}
}

// White : uniform [-1,1]
func (n *Noise) processWhite(block *Block) {
	for i := 0; i < BlockSize; i++ {
		x := n.uniform()
		block.L[i] = x
		block.R[i] = x
	}
}

// Gaussian : approx via uniforms sum
func (n *Noise) processGaussian(block *Block) {
	for i := 0; i < BlockSize; i++ {
		var s float32
		for k := 0; k < 6; k++ {
			s += n.uniform() // [-1,1]
		}
		// s ∈ [-6,6], normalize ~[-1,1] and keep ~N(0, 1/3)
		x := s * (1.0 / 6.0)
		block.L[i] = x
		block.R[i] = x
	}
}

// Pink : Voss–McCartney
func (n *Noise) processPink(block *Block) {
	const numRows = 16
	const norm = 1.0 / numRows

	for i := 0; i < BlockSize; i++ {
		n.pinkCounter++
		c := n.pinkCounter
		acc := n.pinkAccum

		for r := 0; r < numRows; r++ {
			mask := uint32(1) << r
			if c&mask != 0 {
				acc -= n.pinkRows[r]
				v := n.uniform()
				n.pinkRows[r] = v
				acc += v
			}
		}
		n.pinkAccum = acc

		x := acc * norm // [-1,1]
		block.L[i] = x
		block.R[i] = x
	}
}

// Brown / Red : random walk white with leak
func (n *Noise) processBrown(block *Block) {
	const leak = 0.995  // slight leak to avoid drift
	const step = 0.02   // step size
	const clamp = 0.999 // clamp [-1,1]

	b := n.brown
	for i := 0; i < BlockSize; i++ {
		w := n.uniform() // [-1,1]
		b = leak*b + step*w
		if b > clamp {
			b = clamp
		} else if b < -clamp {
			b = -clamp
		}
		block.L[i] = b
		block.R[i] = b
	}
	n.brown = b
}

// Blue : differentiated white
func (n *Noise) processBlue(block *Block) {
	prev := n.lastWhite
	for i := 0; i < BlockSize; i++ {
		w := n.uniform()
		x := w - prev
		prev = w
		block.L[i] = x
		block.R[i] = x
	}
	n.lastWhite = prev
}

func (n *Noise) xorShift32() uint32 {
	x := n.rng
	x ^= x << 13
	x ^= x >> 17
	x ^= x << 5
	n.rng = x
	return x
}

func (n *Noise) uniform() float32 {
	const inv32 = 1.0 / 4294967296.0 // 1 / 2^32
	u := float32(n.xorShift32()) * inv32
	return 2*u - 1
}
