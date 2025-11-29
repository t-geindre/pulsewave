package dsp

import "math"

type LowPassSVF struct {
	Src    Node
	Cutoff Param // Hz
	ResonQ Param // Q (≈ resonance)
	sr     float64

	ic1L, ic2L float64
	ic1R, ic2R float64

	tmp Block
}

func NewLowPassSVF(sr float64, src Node, cutoff Param, q Param) *LowPassSVF {
	return &LowPassSVF{
		Src:    src,
		Cutoff: cutoff,
		ResonQ: q,
		sr:     sr,
	}
}

// Process low-pass TPT SVF (Zavalishin) on tmp -> b.
func (f *LowPassSVF) Process(b *Block) {
	f.tmp.Cycle = b.Cycle
	f.Src.Process(&f.tmp)

	cb := f.Cutoff.Resolve(b.Cycle)
	qb := f.ResonQ.Resolve(b.Cycle)

	nyq := 0.5 * f.sr
	minFc := 5.0
	maxFc := 0.49 * nyq

	for i := 0; i < BlockSize; i++ {
		fc := float64(cb[i])
		if fc < minFc {
			fc = minFc
		}
		if fc > maxFc {
			fc = maxFc
		}

		Q := float64(qb[i])
		if Q < 0.3 {
			Q = 0.3
		}
		if Q > 20 {
			Q = 20
		}

		g := math.Tan(math.Pi * fc / f.sr)
		R := 1.0 / Q
		h := 1.0 / (1.0 + R*g + g*g)

		// Left channel
		xL := float64(f.tmp.L[i])

		v1L := (f.ic1L + g*(xL-f.ic2L)) * h
		v2L := f.ic2L + g*v1L

		lpL := v2L

		f.ic1L = 2*v1L - f.ic1L
		f.ic2L = 2*v2L - f.ic2L

		b.L[i] = float32(lpL)

		// Right channel
		xR := float64(f.tmp.R[i])

		v1R := (f.ic1R + g*(xR-f.ic2R)) * h
		v2R := f.ic2R + g*v1R

		lpR := v2R

		f.ic1R = 2*v1R - f.ic1R
		f.ic2R = 2*v2R - f.ic2R

		b.R[i] = float32(lpR)
	}
}

func (f *LowPassSVF) Reset(soft bool) {
	if !soft {
		f.ic1L, f.ic2L = 0, 0
		f.ic1R, f.ic2R = 0, 0
	}
	f.Src.Reset(soft)
}
