package oscillator

import "math"

type Sine struct {
	phase      float64 // radians
	phaseShift float64 // radians
	step       float64 // radians/sample
	sr         float64
}

func NewSine(sampleRate, freq float64) *Sine {
	s := &Sine{sr: sampleRate}
	s.SetFreq(freq)
	return s
}

func (s *Sine) SetFreq(freq float64) {
	if freq <= 0 {
		freq = 1
	}
	s.step = 2 * math.Pi * freq / s.sr
}

func (s *Sine) NextSample() float64 {
	v := math.Sin(s.phase + s.phaseShift) // applique le shift à la lecture
	s.phase += s.step
	if s.phase >= 2*math.Pi {
		s.phase -= 2 * math.Pi
	}
	return v
}

func (s *Sine) ResetPhase() { s.phase = 0 } // reset propre (le shift est lu)
func (s *Sine) SetPhaseShift(cycles float64) {
	// cycles -> [0..1) -> radians
	c := math.Mod(cycles, 1.0)
	if c < 0 {
		c += 1.0
	}
	s.phaseShift = c * 2 * math.Pi
}
