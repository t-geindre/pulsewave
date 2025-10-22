package oscillator

import "synth/audio"

type Oscillator interface {
	audio.Source
	SetFreq(freq float64)
	SetPhaseShift(phase float64)
}
