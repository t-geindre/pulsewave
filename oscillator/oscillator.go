package oscillator

type Oscillator interface {
	NextSample() float64
	SetFreq(freq float64)
}
