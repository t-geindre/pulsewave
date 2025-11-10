package dsp

import (
	"math"
)

type UnisonFactory func(phase Param, detuneSemi Param) Node

type uniSlot struct {
	voice  Node
	input  *Input
	phase  *SmoothedParam
	detune *SmoothedParam
	gain   float32
	xPos   float64
}

type Unison struct {
	*Mixer

	voices []*uniSlot

	PanSpread    Param // 0..1
	PhaseSpread  Param // 0..1 (cycles)
	DetuneSpread Param // cents
	CurveGamma   Param // >0 ; 1 = linear ; >1 = more center ; <1 = more edges
	NumVoices    Param

	factory UnisonFactory

	lastCycle     uint64
	pendingVoices int // >0 => rebuild next noteOn
	sr            float64
}

type UnisonOpts struct {
	SampleRate float64
	Factory    UnisonFactory

	NumVoices    Param
	PanSpread    Param // 0..1
	PhaseSpread  Param // 0..1 (cycles)
	DetuneSpread Param // cents
	CurveGamma   Param // >0 ; 1 = linear ; >1 = more center ; <1 = more edges
}

func NewUnison(o UnisonOpts) *Unison {
	u := &Unison{
		Mixer:        NewMixer(nil, false),
		PanSpread:    o.PanSpread,
		PhaseSpread:  o.PhaseSpread,
		DetuneSpread: o.DetuneSpread,
		CurveGamma:   o.CurveGamma,
		NumVoices:    o.NumVoices,
		factory:      o.Factory,
		sr:           o.SampleRate,
	}
	u.rebuild(int(o.NumVoices.Resolve(0)[0]))
	return u
}

func (u *Unison) SetVoices(n int) {
	u.pendingVoices = n
}

func (u *Unison) Reset(soft bool) {
	if u.pendingVoices > 0 {
		u.rebuild(u.pendingVoices)
		u.pendingVoices = 0
	}
	u.Mixer.Reset(soft)
}

func (u *Unison) Process(b *Block) {
	u.apply(b.Cycle)
	u.Mixer.Process(b)
}

func (u *Unison) rebuild(n int) {
	if n < 1 {
		n = 1
	}
	u.voices = make([]*uniSlot, n)
	u.Mixer.Clear()

	// Normalized gain
	norm := float32(1.0 / math.Sqrt(float64(n)))

	for i := 0; i < n; i++ {
		var x float64
		if n == 1 {
			x = 0
		} else {
			x = 2*float64(i)/float64(n-1) - 1
		}

		ph := NewSmoothedParam(u.sr, 0.5, 0.002) // center 0.5
		dt := NewSmoothedParam(u.sr, 0.0, 0.002) // semitones

		v := u.factory(ph, dt)

		in := NewInput(v, NewParam(norm), NewParam(0))
		u.Mixer.Add(in)

		u.voices[i] = &uniSlot{
			voice:  v,
			input:  in,
			phase:  ph,
			detune: dt,
			gain:   norm,
			xPos:   x,
		}
	}
}

func (u *Unison) apply(cycle uint64) {
	if u.lastCycle == cycle {
		return
	}
	u.lastCycle = cycle

	panSp := float32(0)
	if u.PanSpread != nil {
		panSp = clamp01(u.PanSpread.Resolve(cycle)[0])
	}

	phSp := float32(0)
	if u.PhaseSpread != nil {
		phSp = clamp01(u.PhaseSpread.Resolve(cycle)[0])
	}

	dtSemiMax := float32(0)
	if u.DetuneSpread != nil {
		dtSemiMax = u.DetuneSpread.Resolve(cycle)[0] / 100.0
	}

	gamma := float64(1)
	if u.CurveGamma != nil {
		g := float64(u.CurveGamma.Resolve(cycle)[0])
		if g > 0 {
			gamma = g
		}
	}

	for _, v := range u.voices {
		c := float32(centeredPower(v.xPos, gamma))

		pan := c * panSp
		v.input.Pan.SetBase(pan)

		v.phase.SetBase(0.5 + 0.5*c*phSp)

		v.detune.SetBase(c * dtSemiMax)
	}

	nVoice := int(u.NumVoices.Resolve(cycle)[0])
	if nVoice != len(u.voices) {
		u.rebuild(nVoice)
	}
}
