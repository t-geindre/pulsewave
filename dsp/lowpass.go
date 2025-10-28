package dsp

import (
	"math"
	"synth/audio"
)

type LowPassSVF struct {
	Src    Node
	Cutoff Param // Hz
	ResonQ Param // Q (1/sqrt(2) ≈ 0.707, clamp [0.3..20]
	sr     float64

	// Canal states
	ic1L, ic2L float64
	ic1R, ic2R float64
}

func NewLowPassSVF(sr float64, src Node, cutoff Param, q Param) *LowPassSVF {
	return &LowPassSVF{
		Src:    src,
		Cutoff: cutoff,
		ResonQ: q,
		sr:     sr,
	}
}

func (f *LowPassSVF) Process(b *audio.Block) {
	var in audio.Block
	in.Cycle = b.Cycle
	f.Src.Process(&in)

	cb := f.Cutoff.Resolve(b.Cycle)
	qb := f.ResonQ.Resolve(b.Cycle)

	nyq := 0.5 * f.sr
	minFc := 5.0
	maxFc := 0.49 * nyq

	for i := 0; i < audio.BlockSize; i++ {
		x := float64(in.L[i])

		fc := float64(cb[i])
		if fc < minFc {
			fc = minFc
		}
		if fc > maxFc {
			fc = maxFc
		}
		g := math.Tan(math.Pi * fc / f.sr)

		Q := float64(qb[i])
		if Q <= 0.3 {
			Q = 0.3
		}
		if Q > 20 {
			Q = 20
		}
		R := 1.0 / Q

		// ZDF SVF (Zavalishin)
		h := 1.0 / (1.0 + g*(g+R))
		hp := (x - R*f.ic1L - f.ic2L) * h
		bp := g*hp + f.ic1L
		lp := g*bp + f.ic2L
		// update states
		f.ic1L = g*hp + bp
		f.ic2L = g*bp + lp

		b.L[i] = float32(lp)

		xr := float64(in.R[i])

		hp = (xr - R*f.ic1R - f.ic2R) * h
		bp = g*hp + f.ic1R
		lp = g*bp + f.ic2R
		f.ic1R = g*hp + bp
		f.ic2R = g*bp + lp

		b.R[i] = float32(lp)
	}
}

func (f *LowPassSVF) Reset() {
	f.ic1L = 0
	f.ic2L = 0
	f.ic1R = 0
	f.ic2R = 0
	f.Src.Reset()
}
