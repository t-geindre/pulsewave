package oscillator

type Saw struct {
	sr         float64 // sample rate
	freq       float64 // Hz
	phase      float64 // [0..1)
	step       float64 // freq/sr
	phaseShift float64 // [0..1)
}

func NewSaw(sampleRate float64) *Saw {
	return &Saw{sr: sampleRate}
}

func (s *Saw) SetFreq(freq float64) {
	if freq <= 0 {
		freq = 1
	}
	nyq := s.sr * 0.5
	if freq > nyq*0.999 {
		freq = nyq * 0.999
	}
	s.freq = freq
	s.step = s.freq / s.sr
}

func (s *Saw) NextValue() (float64, float64) {
	t := frac01(s.phase + s.phaseShift)
	dt := s.step

	y := 2.0*t - 1.0

	y -= polyBLEP(t, dt)

	s.phase += dt
	if s.phase >= 1.0 {
		s.phase -= 1.0
	}

	return y, y
}

func (s *Saw) Reset() { s.phase = 0 }

func (s *Saw) SetPhaseShift(p float64) { s.phaseShift = frac01(p) }

func (s *Saw) IsActive() bool {
	return true
}

func (s *Saw) NoteOn(freq, _ float64) {
	s.SetFreq(freq)
	s.Reset()
}

func (s *Saw) NoteOff() {
}
