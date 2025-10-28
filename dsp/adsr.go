package dsp

import (
	"math"
	"synth/audio"
	"time"
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

	Atk, Dec, Sus, Rel float64 // secondes, secondes, [0..1], secondes

	state State
	value float32 // niveau courant [0..1]
	gate  bool

	// coefficients par-sample (calculés par bloc)
	aCoef, dCoef, rCoef float32

	buf       [audio.BlockSize]float32
	stampedAt uint64
}

func NewADSR(sr float64, atk, dec time.Duration, sus float64, rel time.Duration) *ADSR {
	return &ADSR{
		sr:  sr,
		Atk: atk.Seconds(),
		Dec: dec.Seconds(),
		Sus: sus,
		Rel: rel.Seconds(),
	}
}

func (a *ADSR) NoteOn() {
	a.gate = true
	if a.state == EnvIdle || a.state == EnvRelease {
		a.state = EnvAttack
	}
}

func (a *ADSR) Reset() {
	a.state = EnvIdle
	a.NoteOn()
}

func (a *ADSR) NoteOff() {
	a.gate = false
	if a.state != EnvIdle {
		a.state = EnvRelease
	}
}

func coefFromTime(t, sr float64) float32 {
	if t <= 0 {
		return 0
	} // passe quasi-instantanée
	// One-pole approx: y += (target - y)*(1 - exp(-1/(t*sr)))
	return float32(1 - math.Exp(-1.0/(t*sr)))
}

func (a *ADSR) prepare() {
	a.aCoef = coefFromTime(a.Atk, a.sr)
	a.dCoef = coefFromTime(a.Dec, a.sr)
	a.rCoef = coefFromTime(a.Rel, a.sr)
}

func (a *ADSR) Resolve(cycle uint64) []float32 {
	if a.stampedAt == cycle {
		return a.buf[:]
	}
	a.prepare()

	for i := 0; i < audio.BlockSize; i++ {
		switch a.state {
		case EnvIdle:
			a.value = 0
		case EnvAttack:
			// vers 1.0
			a.value += (1 - a.value) * a.aCoef
			if a.value > 0.9999 || a.aCoef == 0 {
				a.value = 1
				a.state = EnvDecay
			}
		case EnvDecay:
			// vers EnvSustain
			target := float32(a.Sus)
			a.value += (target - a.value) * a.dCoef
			if (a.dCoef == 0 && a.value == target) || (a.value-target)*(1) <= 1e-6 {
				a.state = EnvSustain
			}
		case EnvSustain:
			a.value = float32(a.Sus)
			if !a.gate {
				a.state = EnvRelease
			}
		case EnvRelease:
			// vers 0.0
			a.value += (0 - a.value) * a.rCoef
			if a.value < 1e-5 {
				a.value = 0
				a.state = EnvIdle
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
