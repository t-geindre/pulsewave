package dsp

import (
	"synth/audio"
)

type Param interface {
	SetBase(value float32)
	Resolve(cycle uint64) []float32
	ModInputs() *[]ParamModInput
}

type ParamSimple struct {
	base      float32
	inputs    []ParamModInput
	buf       [audio.BlockSize]float32
	stampedAt uint64
}

func NewParam(base float32) *SmoothedParam {
	return &SmoothedParam{
		base: base,
	}
}

func (s *ParamSimple) SetBase(v float32)           { s.base = v }
func (s *ParamSimple) ModInputs() *[]ParamModInput { return &s.inputs }

func (s *ParamSimple) Resolve(cycle uint64) []float32 {
	if s.stampedAt == cycle {
		return s.buf[:]
	}

	for i := 0; i < audio.BlockSize; i++ {
		s.buf[i] = s.base
	}

	for _, mi := range s.inputs {
		src := mi.Src.Resolve(cycle)
		for i := 0; i < audio.BlockSize; i++ {
			s.buf[i] += mi.Amount * src[i]
		}
	}

	return s.buf[:]
}
