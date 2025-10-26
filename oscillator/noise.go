package oscillator

import "math/rand"

type Noise struct {
}

func NewNoise() *Noise {
	return &Noise{}
}

func (n *Noise) NextValue() (L, R float64) {
	return rand.Float64()*2 - 1, rand.Float64()*2 - 1
}

func (n *Noise) IsActive() bool {
	return true
}

func (n *Noise) Reset() {
}

func (n *Noise) NoteOn(_, _ float64) {
}

func (n *Noise) NoteOff(float64) {
}

func (n *Noise) SetFreq(float64) {
}

func (n *Noise) SetPhaseShift(float64) {
}
