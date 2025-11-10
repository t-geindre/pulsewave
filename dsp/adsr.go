package dsp

import (
	"math"
)

type State int

const (
	EnvIdle State = iota
	EnvAttack
	EnvDecay
	EnvSustain
	EnvRelease
)

type ADSR struct {
	sr float64

	Atk, Dec, Sus, Rel Param

	state State
	value float32
	gate  bool

	aCoef, dCoef, rCoef float32

	buf       [BlockSize]float32
	stampedAt uint64

	needRecalc bool
}

func NewADSR(sr float64, atk, dec, sus, rel Param) *ADSR {
	a := &ADSR{
		sr:         sr,
		Atk:        atk,
		Dec:        dec,
		Sus:        sus,
		Rel:        rel,
		needRecalc: true,
		state:      EnvIdle,
	}

	return a
}

func (a *ADSR) NoteOn() {
	a.gate = true
	a.setState(EnvAttack)
}

func (a *ADSR) Reset() {
	a.NoteOff()
	a.NoteOn()
}

func (a *ADSR) NoteOff() {
	a.gate = false
	a.setState(EnvRelease)
}

func coefFromTime(t float32, sr float64) float32 {
	if t <= 0 {
		return 0
	}

	return float32(1 - math.Exp(-1.0/(float64(t)*sr)))
}

func (a *ADSR) setState(s State) {
	switch s {
	case EnvRelease:
		if a.state == EnvIdle {
			return
		}
	}

	a.state = s
	a.needRecalc = true

}

func (a *ADSR) recalc(cycle uint64) {
	if !a.needRecalc {
		return
	}
	a.needRecalc = false

	switch a.state {
	case EnvAttack:
		a.aCoef = coefFromTime(a.Atk.Resolve(cycle)[0], a.sr)
	case EnvDecay:
		a.dCoef = coefFromTime(a.Dec.Resolve(cycle)[0], a.sr)
	case EnvRelease:
		a.rCoef = coefFromTime(a.Rel.Resolve(cycle)[0], a.sr)
	}
}

func (a *ADSR) Resolve(cycle uint64) []float32 {
	if a.stampedAt == cycle {
		return a.buf[:]
	}

	a.recalc(cycle)

	for i := 0; i < BlockSize; i++ {
		switch a.state {
		case EnvIdle:
			a.value = 0
		case EnvAttack:
			a.value += (1 - a.value) * a.aCoef
			if a.value > 0.999 || a.aCoef == 0 {
				a.value = 1
				a.setState(EnvDecay)
			}
		case EnvDecay:
			target := a.Sus.Resolve(cycle)[0]
			a.value += (target - a.value) * a.dCoef
			if (a.dCoef == 0 && a.value == target) || (a.value-target)*(1) <= 1e-6 {
				a.setState(EnvSustain)
			}
		case EnvSustain:
			a.value = a.Sus.Resolve(cycle)[0]
			if !a.gate {
				a.setState(EnvRelease)
			}
		case EnvRelease:
			a.value += (0 - a.value) * a.rCoef
			if a.value < 0.001 || a.rCoef == 0 {
				a.value = 0
				a.setState(EnvIdle)
			}
		}
		a.buf[i] = a.value
	}
	a.stampedAt = cycle
	return a.buf[:]
}

func (a *ADSR) IsIdle() bool {
	return a.state == EnvIdle
}
