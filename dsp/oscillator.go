package dsp

import (
	"math"
	"synth/audio"
)

type Oscillator struct {
	shapeRegistry *ShapeRegistry
	shapeIndex    int

	freq       Param
	width      Param // ShapeSquare only
	phaseShift Param

	phase float64
	sr    float64

	rng uint32 // ShapeNoise only

	buf       [audio.BlockSize]float32
	stampedAt uint64
}

func NewOscillator(
	sampleRate float64,
	shape OscShape,
	freq Param,
	phaseShift Param,
	sqrWidth Param,
) *Oscillator {
	reg := NewShapeRegistry()
	reg.Set(0, shape)
	return NewRegOscillator(sampleRate, reg, 0, freq, phaseShift, sqrWidth)
}

// NewRegOscillator creates a new Oscillator instance.
// phaseShift: in [0..1] cycles,
// sqrWidth: in [0..1] duty cycle for ShapeSquare wave.
func NewRegOscillator(
	sampleRate float64,
	shapeRegistry *ShapeRegistry,
	shapeIndex int,
	freq Param,
	phaseShift Param,
	sqrWidth Param,
) *Oscillator {
	return &Oscillator{
		sr:            sampleRate,
		shapeRegistry: shapeRegistry,
		shapeIndex:    shapeIndex,
		freq:          freq,
		phaseShift:    phaseShift,
		width:         sqrWidth,
		rng:           0x9E3779B9, // arbitrary non-zero seed
	}
}

func (s *Oscillator) Reset() {
	s.phase = 0
}

func (s *Oscillator) Process(block *audio.Block) {
	v := s.Resolve(block.Cycle)
	for i := 0; i < audio.BlockSize; i++ {
		block.L[i] = v[i]
		block.R[i] = v[i]
	}
}

func (s *Oscillator) Resolve(cycle uint64) []float32 {
	if s.stampedAt == cycle {
		return s.buf[:]
	}

	shape := s.shapeRegistry.Get(s.shapeIndex)

	if shape == ShapeNoise {
		for i := 0; i < audio.BlockSize; i++ {
			x := s.xorShift32()
			u := float32(x) * (1.0 / 4294967296.0)
			s.buf[i] = 2*u - 1
		}
		s.stampedAt = cycle
		return s.buf[:]
	}

	fb := s.freq.Resolve(cycle)

	var wb []float32
	if shape == ShapeSquare && s.width != nil {
		wb = s.width.Resolve(cycle)
	}

	var phb []float32
	if s.phaseShift != nil {
		phb = s.phaseShift.Resolve(cycle) // 0..1 tours
	}

	const twoPi = 2 * math.Pi
	const invTwoPi = 1.0 / twoPi
	k := twoPi / s.sr

	for i := 0; i < audio.BlockSize; i++ {
		p := s.phase * invTwoPi
		p -= math.Floor(p)

		var shift float64
		if phb != nil {
			shift = float64(phb[i])
			p += shift
			p -= math.Floor(p)
		}

		switch shape {
		case ShapeSine:
			if shift != 0 {
				s.buf[i] = float32(math.Sin(s.phase + twoPi*shift))
			} else {
				s.buf[i] = float32(math.Sin(s.phase))
			}

		case ShapeSaw:
			y := float32(2*p - 1)
			dt := math.Abs(float64(fb[i])) / s.sr
			if dt > 0.5 {
				dt = 0.5
			}
			y -= polyBLEP(p, dt)
			s.buf[i] = y

		case ShapeTriangle:
			tri := 1.0 - 4.0*math.Abs(p-0.5)
			s.buf[i] = float32(tri)

		case ShapeSquare:
			duty := 0.5
			if wb != nil {
				d := float64(wb[i])
				if d < 0.01 {
					d = 0.01
				} else if d > 0.99 {
					d = 0.99
				}
				duty = d
			}
			dt := math.Abs(float64(fb[i])) / s.sr
			if dt > 0.5 {
				dt = 0.5
			}

			y1 := float32(2*p - 1)
			y1 -= polyBLEP(p, dt)

			pd := p + duty
			pd -= math.Floor(pd)
			y2 := float32(2*pd - 1)
			y2 -= polyBLEP(pd, dt)

			s.buf[i] = 0.5 * (y1 - y2)
		}

		s.phase += k * float64(fb[i])
		if s.phase >= twoPi {
			s.phase -= twoPi
		} else if s.phase < 0 {
			s.phase += twoPi
		}
	}

	s.stampedAt = cycle
	return s.buf[:]
}

func (s *Oscillator) xorShift32() uint32 {
	x := s.rng
	x ^= x << 13
	x ^= x >> 17
	x ^= x << 5
	s.rng = x
	return x
}
