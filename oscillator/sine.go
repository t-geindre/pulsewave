package oscillator

import "math"

type Sine struct {
	phase float64
	step  float64
	sr    float64
}

func NewSine(sampleRate, freq float64) *Sine {
	s := &Sine{
		phase: 0,
		sr:    sampleRate,
	}

	s.SetFreq(freq)

	return s
}

func (s *Sine) SetFreq(freq float64) {
	s.step = 2 * math.Pi * freq / s.sr
}

func (s *Sine) NextSample() float64 {
	v := math.Sin(s.phase)
	s.phase += s.step

	if s.phase >= 2*math.Pi {
		s.phase -= 2 * math.Pi
	}

	return v
}
