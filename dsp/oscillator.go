package dsp

import "math"

type Oscillator struct {
	shapeRegistry *ShapeRegistry
	shapeIndex    Param

	shapeIdx  float32
	shape     OscShape
	wavetable *Wavetable

	freq       Param
	width      Param // ShapeSquare only
	phaseShift Param

	phase float64
	sr    float64
	invSr float64 // 1/sr

	buf       [BlockSize]float32
	stampedAt uint64
}

func NewRegOscillator(
	sampleRate float64,
	shapeRegistry *ShapeRegistry,
	shapeIndex Param,
	freq Param,
	phaseShift Param,
	sqrWidth Param,
) *Oscillator {
	return &Oscillator{
		sr:            sampleRate,
		invSr:         1.0 / sampleRate,
		shapeRegistry: shapeRegistry,
		shapeIndex:    shapeIndex,
		freq:          freq,
		phaseShift:    phaseShift,
		width:         sqrWidth,
		shapeIdx:      -1,
	}
}

func (s *Oscillator) Reset(soft bool) {
	if soft {
		return
	}
	s.phase = 0
}

func (s *Oscillator) Process(block *Block) {
	v := s.Resolve(block.Cycle)
	for i := 0; i < BlockSize; i++ {
		block.L[i] = v[i]
		block.R[i] = v[i]
	}
}

func (s *Oscillator) Resolve(cycle uint64) []float32 {
	if s.stampedAt == cycle {
		return s.buf[:]
	}

	shapeBuf := s.shapeIndex.Resolve(cycle)
	sIdx := shapeBuf[0]
	if sIdx != s.shapeIdx {
		shape, wavetable := s.shapeRegistry.Get(sIdx)
		s.shape = shape
		s.wavetable = wavetable
		s.shapeIdx = sIdx
	}

	fb := s.freq.Resolve(cycle)

	var wb []float32
	if s.shape == ShapeSquare && s.width != nil {
		wb = s.width.Resolve(cycle)
	}

	var phb []float32
	if s.phaseShift != nil {
		phb = s.phaseShift.Resolve(cycle) // 0..1 cycle
	}

	switch s.shape {
	case ShapeSaw:
		s.processSaw(fb, phb)
	case ShapeTriangle:
		s.processTriangle(fb, phb)
	case ShapeSquare:
		s.processSquare(fb, wb, phb)
	case ShapeTableWave:
		s.processTable(fb, phb)
	}

	s.stampedAt = cycle
	return s.buf[:]
}

func (s *Oscillator) processSaw(fb, phb []float32) {
	const twoPi = 2 * math.Pi
	const invTwoPi = 1.0 / twoPi
	k := twoPi * s.invSr

	for i := 0; i < BlockSize; i++ {
		p := s.phase * invTwoPi
		if phb != nil {
			p += float64(phb[i])
		}

		if p >= 1 {
			p -= 1
		} else if p < 0 {
			p += 1
		}

		f := float64(fb[i])
		dt := math.Abs(f) * s.invSr
		if dt > 0.5 {
			dt = 0.5
		}

		y := float32(2*p - 1)
		y -= polyBLEP(p, dt)
		s.buf[i] = y

		s.phase += k * f
		if s.phase >= twoPi {
			s.phase -= twoPi
		} else if s.phase < 0 {
			s.phase += twoPi
		}
	}
}

func (s *Oscillator) processTriangle(fb, phb []float32) {
	const twoPi = 2 * math.Pi
	const invTwoPi = 1.0 / twoPi
	k := twoPi * s.invSr

	for i := 0; i < BlockSize; i++ {
		p := s.phase * invTwoPi
		if phb != nil {
			p += float64(phb[i])
		}

		if p >= 1 {
			p -= 1
		} else if p < 0 {
			p += 1
		}

		tri := 1.0 - 4.0*math.Abs(p-0.5)
		s.buf[i] = float32(tri)

		f := float64(fb[i])
		s.phase += k * f
		if s.phase >= twoPi {
			s.phase -= twoPi
		} else if s.phase < 0 {
			s.phase += twoPi
		}
	}
}

func (s *Oscillator) processSquare(fb, wb, phb []float32) {
	const twoPi = 2 * math.Pi
	const invTwoPi = 1.0 / twoPi
	k := twoPi * s.invSr

	for i := 0; i < BlockSize; i++ {
		p := s.phase * invTwoPi
		if phb != nil {
			p += float64(phb[i])
		}

		if p >= 1 {
			p -= 1
		} else if p < 0 {
			p += 1
		}

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

		f := float64(fb[i])
		dt := math.Abs(f) * s.invSr
		if dt > 0.5 {
			dt = 0.5
		}

		y1 := float32(2*p - 1)
		y1 -= polyBLEP(p, dt)

		pd := p + duty
		if pd >= 1 {
			pd -= 1
		}
		y2 := float32(2*pd - 1)
		y2 -= polyBLEP(pd, dt)

		s.buf[i] = 0.5 * (y1 - y2)

		s.phase += k * f
		if s.phase >= twoPi {
			s.phase -= twoPi
		} else if s.phase < 0 {
			s.phase += twoPi
		}
	}
}

func (s *Oscillator) processTable(fb, phb []float32) {
	const twoPi = 2 * math.Pi
	const invTwoPi = 1.0 / twoPi
	k := twoPi * s.invSr

	if s.wavetable == nil || s.wavetable.Size == 0 {
		for i := 0; i < BlockSize; i++ {
			s.buf[i] = 0
		}
		return
	}

	size := s.wavetable.Size
	lastIdx := size - 1
	sizeF := float64(size)

	for i := 0; i < BlockSize; i++ {
		p := s.phase * invTwoPi
		if phb != nil {
			p += float64(phb[i])
		}

		if p >= 1 {
			p -= 1
		} else if p < 0 {
			p += 1
		}

		pos := p * float64(lastIdx)
		if pos >= sizeF {
			pos -= sizeF
		}

		idx := int(pos)
		next := idx + 1
		if next >= size {
			next = 0
		}

		f := float32(pos - float64(idx))
		v1 := s.wavetable.Table[idx]
		v2 := s.wavetable.Table[next]
		s.buf[i] = v1 + f*(v2-v1)

		freq := float64(fb[i])
		s.phase += k * freq
		if s.phase >= twoPi {
			s.phase -= twoPi
		} else if s.phase < 0 {
			s.phase += twoPi
		}
	}
}
