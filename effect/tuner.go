package effect

import (
	"math"
	"synth/audio"
)

type Tuner struct {
	audio.Source
	octave    int
	semitones float64
	cents     float64
}

func NewTuner(src audio.Source) *Tuner {
	return &Tuner{Source: src}
}

func (t *Tuner) NoteOn(freq, vel float64) {
	t.Source.NoteOn(
		freq*
			math.Pow(2, float64(t.octave))*
			math.Pow(2, t.semitones/12.0)*
			math.Pow(2, t.cents/1200.0),
		vel,
	)
}

func (t *Tuner) SetOctaveOffset(o int) { t.octave = o }
func (t *Tuner) NudgeOctave(delta int) { t.octave += delta }

func (t *Tuner) SetTransposeSemitones(semi float64) { t.semitones = semi }
func (t *Tuner) NudgeSemitones(delta float64)       { t.semitones += delta }

func (t *Tuner) SetDetuneCents(c float64) { t.cents = c }
func (t *Tuner) NudgeCents(delta float64) { t.cents += delta }
