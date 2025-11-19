package dsp

import "math"

type SmoothedParam struct {
	alpha float32
	sr    float64
	tc    Param

	base   float32
	inputs []ParamModInput

	last float32

	buf       [BlockSize]float32
	stampedAt uint64
}

// NewSmoothedParam tc: ex: 0.005 cutoff, 0.001 amp/pitch
func NewSmoothedParam(sr float64, base float32, tc Param) *SmoothedParam {
	s := &SmoothedParam{
		base: base,
		last: base,
		sr:   sr,
		tc:   nil,
	}

	// Const, enable fast path
	if cstTc, ok := tc.(*ConstParam); ok {
		t := cstTc.Resolve(0)[0]
		s.alpha = float32(1.0 - math.Exp(-1.0/(float64(t)*sr)))
		return s
	}

	s.tc = tc
	s.alpha = 1

	return s
}

func (s *SmoothedParam) SetBase(v float32)           { s.base = v }
func (s *SmoothedParam) GetBase() float32            { return s.base }
func (s *SmoothedParam) ModInputs() *[]ParamModInput { return &s.inputs }

func (s *SmoothedParam) Resolve(cycle uint64) []float32 {
	if s.stampedAt == cycle {
		return s.buf[:]
	}

	// Base
	for i := 0; i < BlockSize; i++ {
		s.buf[i] = s.base
	}

	// Modulation
	for _, mi := range s.inputs {
		src := mi.Src().Resolve(cycle)
		amount := mi.Amount().Resolve(cycle)
		mapf := mi.Map()
		if mapf == nil {
			for i := 0; i < BlockSize; i++ {
				s.buf[i] += amount[i] * src[i]
			}
		} else {
			for i := 0; i < BlockSize; i++ {
				s.buf[i] += amount[i] * mapf(src[i])
			}
		}
	}

	// Tc resolve
	if s.tc != nil {
		tcBuf := s.tc.Resolve(cycle) // seconds
		t := float64(tcBuf[0])
		s.alpha = fastSmoothAlpha(t, s.sr)
	}

	alpha := s.alpha
	cur := s.last

	// Smoothing
	if alpha >= 1 {
		for i := 0; i < BlockSize; i++ {
			cur = s.buf[i]
			s.buf[i] = cur
		}
	} else if alpha <= 0 {
		for i := 0; i < BlockSize; i++ {
			s.buf[i] = cur
		}
	} else {
		for i := 0; i < BlockSize; i++ {
			cur += alpha * (s.buf[i] - cur)
			s.buf[i] = cur
		}
	}

	s.last = cur
	s.stampedAt = cycle
	return s.buf[:]
}
