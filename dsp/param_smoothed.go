package dsp

import (
	"math"
)

type SmoothedParam struct {
	alpha float32

	base   float32
	inputs []ParamModInput

	last float32

	buf       [BlockSize]float32
	stampedAt uint64
}

// NewSmoothedParam tc: ex: 0.005 cutoff, 0.001 amp/pitch
func NewSmoothedParam(sr float64, base float32, tc float64) *SmoothedParam {
	return &SmoothedParam{
		base:  base,
		last:  base,
		alpha: float32(1.0 - math.Exp(-1.0/(tc*sr))),
	}
}

func (s *SmoothedParam) SetBase(v float32)           { s.base = v }
func (s *SmoothedParam) GetBase() float32            { return s.base }
func (s *SmoothedParam) ModInputs() *[]ParamModInput { return &s.inputs }

func (s *SmoothedParam) Resolve(cycle uint64) []float32 {
	if s.stampedAt == cycle {
		return s.buf[:]
	}

	for i := 0; i < BlockSize; i++ {
		s.buf[i] = s.base
	}

	for _, mi := range s.inputs {
		src := mi.Src().Resolve(cycle) // read-only
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

	cur := s.last
	for i := 0; i < BlockSize; i++ {
		cur += s.alpha * (s.buf[i] - cur)
		s.buf[i] = cur
	}
	s.last = cur

	s.stampedAt = cycle
	return s.buf[:]
}
