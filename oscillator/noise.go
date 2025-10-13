package oscillator

import "math/rand"

type Noise struct {
}

func NewNoise() *Noise {
	return &Noise{}
}

func (n Noise) NextSample() float64 {
	return rand.Float64()*2 - 1 // Range -1 to 1
}

func (n Noise) SetFreq(float64) {
}

func (n Noise) ResetPhase() {
}

func (n Noise) SetPhaseShift(float64) {
}
