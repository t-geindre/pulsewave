package dsp

import (
	"math"
)

type FeedbackDelay struct {
	Src Source

	Time     Param // secondes (0..max)
	Feedback Param // 0..~0.95 (clampé)
	Mix      Param // 0..1
	ToneHz   Param // cutoff LPF dans la boucle (ex. 1000 Hz). nil => pas de LPF

	sr       float64
	maxSamps int
	wpos     int
	bufL     []float32
	bufR     []float32

	// LPF state
	lpfL, lpfR float32

	tmp Block
}

func NewFeedbackDelay(sr float64, maxDelaySeconds float64, src Source,
	time Param, feedback Param, mix Param, toneHz Param,
) *FeedbackDelay {
	if maxDelaySeconds <= 0 {
		maxDelaySeconds = 2.0
	}
	maxSamps := int(math.Ceil(maxDelaySeconds * sr))

	// +2 margin interpolation
	bufL := make([]float32, maxSamps+2)
	bufR := make([]float32, maxSamps+2)

	return &FeedbackDelay{
		Src:      src,
		Time:     time,
		Feedback: feedback,
		Mix:      mix,
		ToneHz:   toneHz,
		sr:       sr,
		maxSamps: len(bufL),
		bufL:     bufL,
		bufR:     bufR,
	}
}

func readInterp(buf []float32, idx float64) float32 {
	n := float64(len(buf))
	idx = math.Mod(idx, n)
	if idx < 0 {
		idx += n
	}
	i0 := int(idx)
	i1 := i0 + 1
	if i1 >= len(buf) {
		i1 -= len(buf)
	}
	f := float32(idx - float64(i0))
	a := buf[i0]
	b := buf[i1]
	return a + (b-a)*f
}

func (d *FeedbackDelay) Process(b *Block) {
	d.tmp.Cycle = b.Cycle
	d.Src.Process(&d.tmp)

	tb := d.safeResolve(d.Time, b.Cycle)
	fb := d.safeResolve(d.Feedback, b.Cycle)
	mb := d.safeResolve(d.Mix, b.Cycle)

	var alpha float32
	if d.ToneHz != nil {
		tone := d.ToneHz.Resolve(b.Cycle)
		alpha = fastLpfCoef(float64(tone[0]), d.sr)
	} else {
		alpha = 0
	}

	N := float64(d.maxSamps)
	w := float64(d.wpos)

	for i := 0; i < BlockSize; i++ {
		xL := d.tmp.L[i]
		xR := d.tmp.R[i]

		delayS := float64(tb[i]) * d.sr
		if delayS < 1 {
			delayS = 1
		}
		if delayS > N-2 {
			delayS = N - 2
		}

		rIdx := w - delayS

		yL := readInterp(d.bufL, rIdx)
		yR := readInterp(d.bufR, rIdx)

		fyL, fyR := yL, yR
		if alpha > 0 {
			d.lpfL += alpha * (yL - d.lpfL)
			d.lpfR += alpha * (yR - d.lpfR)
			fyL, fyR = d.lpfL, d.lpfR
		}

		fbk := fb[i]
		if fbk < 0 {
			fbk = 0
		}
		if fbk > 0.97 {
			fbk = 0.97
		}

		wp := int(w)
		if wp >= len(d.bufL) {
			wp -= len(d.bufL)
		}
		d.bufL[wp] = xL + fbk*fyL
		d.bufR[wp] = xR + fbk*fyR

		// Dry/Wet mix
		mix := clamp01(mb[i])
		dry := 1 - mix
		outL := dry*xL + mix*yL
		outR := dry*xR + mix*yR

		b.L[i] = outL
		b.R[i] = outR

		w++
		if w >= N {
			w -= N
		}
	}

	d.wpos = int(w)
}

func (d *FeedbackDelay) safeResolve(p Param, cycle uint64) []float32 {
	if p == nil {
		return zeroBlock[:]
	}
	return p.Resolve(cycle)
}

func (d *FeedbackDelay) Reset(bool) {
	// delay comes after voicing, no reset needed
}

// todo move this elsewhere
var zeroBlock = func() [BlockSize]float32 {
	var z [BlockSize]float32
	for i := 0; i < BlockSize; i++ {
		z[i] = 0
	}
	return z
}()
