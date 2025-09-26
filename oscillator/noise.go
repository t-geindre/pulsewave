package oscillator

import "math/rand"

type Noise struct {
}

func NewNoise() *Noise {
	return &Noise{}
}

func (n Noise) NextSample() float64 {
	return rand.Float64()*2 - 1
}

func (n Noise) SetFreq(float64) {
}
