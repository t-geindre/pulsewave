package oscillator

import "math"

type Triangle struct {
	sr         float64
	freq       float64
	phase      float64 // [0..1)
	phaseShift float64 // [0..1)
	step       float64 // freq/sr
}

func NewTriangle(sampleRate float64) *Triangle {
	return &Triangle{sr: sampleRate}
}

func (t *Triangle) SetFreq(freq float64) {
	if freq <= 0 {
		freq = 1
	}
	nyq := t.sr * 0.5
	if freq > nyq*0.999 {
		freq = nyq * 0.999
	}
	t.freq = freq
	t.step = t.freq / t.sr
}

func (t *Triangle) NextValue() (float64, float64) {
	ph := frac01(t.phase + t.phaseShift)

	y := 1.0 - 4.0*math.Abs(ph-0.5)

	t.phase += t.step
	if t.phase >= 1.0 {
		t.phase -= 1.0
	}

	return y, y
}

func (t *Triangle) Reset()                  { t.phase = 0 }
func (t *Triangle) SetPhaseShift(p float64) { t.phaseShift = frac01(p) }

func (t *Triangle) IsActive() bool {
	return true
}

func (t *Triangle) NoteOn(freq, velocity float64) {
	if freq > 0 {
		t.SetFreq(freq)
	}
	t.Reset()
}

func (t *Triangle) NoteOff(float64) {
}
