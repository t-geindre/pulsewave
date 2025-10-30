package dsp

type TunerParam struct {
	Param
	st        Param
	buf       [BlockSize]float32
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

	// Fast path: if semi is (nearly) constant across the block, compute once.
	const eps = 1e-9
	s0 := semi[0]
	isFlat := true
	for i := 1; i < BlockSize; i++ {
		if diff := semi[i] - s0; diff > eps || diff < -eps {
			isFlat = false
			break
		}
	}
	if isFlat {
		r := fastExpSemi(s0)
		for i := 0; i < BlockSize; i++ {
			p.buf[i] = base[i] * r
		}
		p.stampedAt = cycle
		return p.buf[:]
	}

	// General path
	for i := 0; i < BlockSize; i++ {
		p.buf[i] = base[i] * fastExpSemi(semi[i])
	}
	p.stampedAt = cycle
	return p.buf[:]
}
