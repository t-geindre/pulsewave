package oscillator

import "math"

type Tuned struct {
	Oscillator
	octave    int
	semitones float64
	cents     float64
}

func NewTuned(osc Oscillator) *Tuned {
	return &Tuned{Oscillator: osc}
}

func (t *Tuned) SetFreq(freq float64) {
	t.Oscillator.SetFreq(
		freq *
			math.Pow(2, float64(t.octave)) *
			math.Pow(2, t.semitones/12.0) *
			math.Pow(2, t.cents/1200.0),
	)
}

func (t *Tuned) SetOctaveOffset(o int) { t.octave = o }
func (t *Tuned) NudgeOctave(delta int) { t.octave += delta }

func (t *Tuned) SetTransposeSemitones(semi float64) { t.semitones = semi }
func (t *Tuned) NudgeSemitones(delta float64)       { t.semitones += delta }

func (t *Tuned) SetDetuneCents(c float64) { t.cents = c }
func (t *Tuned) NudgeCents(delta float64) { t.cents += delta }
