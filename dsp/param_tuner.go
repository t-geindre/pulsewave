package dsp

import (
	"math"
	"synth/audio"
)

type TunerParam struct {
	Param
	st        Param
	buf       [audio.BlockSize]float32
	stampedAt uint64
}

// NewTunerParam creates a new TunerParam
// Octave: ±12 st,
// Fifth: +7 st,
// Ton: +2 st,
// Cent:  +0.01 st.
func NewTunerParam(hz, semiTon Param) *TunerParam {
	return &TunerParam{Param: hz, st: semiTon}
}

func (p *TunerParam) Resolve(cycle uint64) []float32 {
	if p.stampedAt == cycle {
		return p.buf[:]
	}
	base := p.Param.Resolve(cycle)
	semi := p.st.Resolve(cycle)
	for i := 0; i < audio.BlockSize; i++ {
		ratio := float32(math.Pow(2, float64(semi[i])/12.0))
		p.buf[i] = base[i] * ratio
	}
	p.stampedAt = cycle
	return p.buf[:]
}
