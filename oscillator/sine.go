package oscillator

import "math"

type Sine struct {
	phase      float64 // radians
	phaseShift float64 // radians
	step       float64 // radians/sample
	sr         float64
}

func NewSine(sampleRate float64) *Sine {
	return &Sine{sr: sampleRate}
}

func (s *Sine) NextValue() (L, R float64) {
	v := math.Sin(s.phase + s.phaseShift) // applique le shift à la lecture
	s.phase += s.step
	if s.phase >= 2*math.Pi {
		s.phase -= 2 * math.Pi
	}
	return v, v
}

func (s *Sine) IsActive() bool {
	return true
}

func (s *Sine) Reset() {
	s.phase = 0
}

func (s *Sine) NoteOn(freq, _ float64) {
	if freq > 0 {
		s.SetFreq(freq)
	}
	s.Reset()
}

func (s *Sine) NoteOff(float64) {
}

func (s *Sine) SetFreq(freq float64) {
	if freq <= 0 {
		freq = 1
	}
	s.step = 2 * math.Pi * freq / s.sr
}

func (s *Sine) SetPhaseShift(cycles float64) {
	c := math.Mod(cycles, 1.0)
	if c < 0 {
		c += 1.0
	}
	s.phaseShift = c * 2 * math.Pi
}
