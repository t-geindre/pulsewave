package oscillator

type Oscillator interface {
	NextSample() float64
	SetFreq(freq float64)
	ResetPhase()
	SetPhaseShift(phase float64)
}
