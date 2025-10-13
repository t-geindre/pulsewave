package oscillator

type Square struct {
	sr         float64
	freq       float64
	phase      float64 // [0..1)
	phaseShift float64 // [0..1)
	step       float64 // freq/sr
	pw         float64 // 0..1
}

func NewSquare(sampleRate, freq float64) *Square {
	s := &Square{sr: sampleRate, pw: 0.5}
	s.SetFreq(freq)
	return s
}

func (s *Square) SetFreq(freq float64) {
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

func (s *Square) NextSample() float64 {
	t := frac01(s.phase + s.phaseShift)
	dt := s.step

	y := -1.0
	if t < s.pw {
		y = 1.0
	}

	y += polyBLEP(t, dt)              // front t=0
	y -= polyBLEP(frac01(t-s.pw), dt) // front t=pw

	s.phase += dt
	if s.phase >= 1.0 {
		s.phase -= 1.0
	}
	return y
}

func (s *Square) SetPulseWidth(pw float64) {
	eps := 1.0 / s.sr
	if pw <= eps {
		pw = eps
	} else if pw >= 1.0-eps {
		pw = 1.0 - eps
	}
	s.pw = pw
}

func (s *Square) ResetPhase()             { s.phase = 0 }
func (s *Square) SetPhaseShift(p float64) { s.phaseShift = frac01(p) }
