package dsp

import (
	"math"
	"synth/audio"
)

type SmoothedParam struct {
	sr           float64
	timeConstant float32 // s, ex: 0.005 cutoff, 0.001 amp/pitch

	base   float32
	inputs []ParamModInput

	last float32

	buf       [audio.BlockSize]float32
	stampedAt uint64
}

func NewSmoothedParam(sr float64, base, tc float32) *SmoothedParam {
	return &SmoothedParam{
		sr:           sr,
		base:         base,
		timeConstant: tc,
		last:         base,
	}
}

func (s *SmoothedParam) SetBase(v float32)           { s.base = v }
func (s *SmoothedParam) ModInputs() *[]ParamModInput { return &s.inputs }

func (s *SmoothedParam) Resolve(cycle uint64) []float32 {
	if s.stampedAt == cycle {
		return s.buf[:]
	}

	for i := 0; i < audio.BlockSize; i++ {
		s.buf[i] = s.base
	}

	for _, mi := range s.inputs {
		src := mi.Src.Resolve(cycle) // read-only
		if mi.Map == nil {
			for i := 0; i < audio.BlockSize; i++ {
				s.buf[i] += mi.Amount * src[i]
			}
		} else {
			for i := 0; i < audio.BlockSize; i++ {
				s.buf[i] += mi.Amount * mi.Map(src[i])
			}
		}
	}

	alpha := float32(1.0 - math.Exp(-1.0/(float64(s.timeConstant)*s.sr)))
	cur := s.last
	for i := 0; i < audio.BlockSize; i++ {
		cur += alpha * (s.buf[i] - cur)
		s.buf[i] = cur
	}
	s.last = cur

	s.stampedAt = cycle
	return s.buf[:]
}
