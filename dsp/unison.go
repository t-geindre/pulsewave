// unison.go
package dsp

import (
	"math"
	"synth/audio"
)

type UnisonFactory func(phase Param, detuneSemi Param) Node

type uniSlot struct {
	voice  Node           // sous-graphe audio de la voix (ex: mixer d'osc)
	input  *Input         // entrée correspondante dans le Mixer interne (pour MAJ Pan/Gain)
	phase  *SmoothedParam // décalage de phase (0..1 tours), modulable à chaud
	detune *SmoothedParam // décalage en demi-tons, modulable à chaud
	gain   float32        // normalisation per-voice (ex: 1/sqrt(N))
	xPos   float64        // position spatiale abstraite [-1..1]
}

type Unison struct {
	*Mixer

	voices []*uniSlot

	PanSpread    Param // 0..1 (écartement stéréo)
	PhaseSpread  Param // 0..1 (tours 0..1)
	DetuneSpread Param // en cents (±), converti en demi-tons ( /100 )
	CurveGamma   Param // >0 ; 1 = linéaire ; >1 = plus de densité au centre

	factory UnisonFactory

	lastCycle     uint64
	pendingVoices int // >0 => rebuild au prochain Reset/NoteOn
	sr            float64
}

// UnisonOpts pour créer l'unison
type UnisonOpts struct {
	SampleRate float64
	NumVoices  int
	Factory    UnisonFactory

	// spreads & shape
	PanSpread    Param // 0..1
	PhaseSpread  Param // 0..1 (tours)
	DetuneSpread Param // en cents
	CurveGamma   Param // >0 (ex: Const(2.0))
}

func NewUnison(o UnisonOpts) *Unison {
	u := &Unison{
		Mixer:        NewMixer(NewParam(1.0), false), // master gain = 1, pas de soft-clip
		PanSpread:    o.PanSpread,
		PhaseSpread:  o.PhaseSpread,
		DetuneSpread: o.DetuneSpread,
		CurveGamma:   o.CurveGamma,
		factory:      o.Factory,
		sr:           o.SampleRate,
	}
	u.rebuild(o.NumVoices)
	return u
}

// SetVoices demande un nouveau nombre de voix (appliqué au prochain Reset/NoteOn).
func (u *Unison) SetVoices(n int) {
	u.pendingVoices = n
}

// Reset applique les changements en attente (rebuild) et reset le Mixer.
func (u *Unison) Reset() {
	if u.pendingVoices > 0 {
		u.rebuild(u.pendingVoices)
		u.pendingVoices = 0
	}
	u.Mixer.Reset()
}

func (u *Unison) Process(b *audio.Block) {
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
}
