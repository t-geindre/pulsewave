package oscillator

import "math"

type Exponential struct {
	sr         float64
	freq       float64
	phase      float64 // [0..1)
	phaseShift float64 // [0..1)
	step       float64 // freq/sr
	shape      float64 // -10..+10 ; 0 => lin
	dcComp     bool
}

func NewExponential(sampleRate float64) *Exponential {
	return &Exponential{sr: sampleRate}
}

func (e *Exponential) SetFreq(freq float64) {
	if freq <= 0 {
		freq = 1
	}
	nyq := e.sr * 0.5
	if freq > nyq*0.999 {
		freq = nyq * 0.999
	}
	e.freq = freq
	e.step = e.freq / e.sr
}

func (e *Exponential) SetPhaseShift(p float64) {
	e.phaseShift = frac01(p)
}

func (e *Exponential) Reset() { e.phase = 0 }

func (e *Exponential) SetShape(s float64) { e.shape = s }
func (e *Exponential) SetDCComp(b bool)   { e.dcComp = b }

func (e *Exponential) mapExpo(t float64) float64 {
	if e.shape == 0 {
		return t
	}
	return math.Expm1(e.shape*t) / math.Expm1(e.shape)
}

func (e *Exponential) NextValue() (float64, float64) {
	t := frac01(e.phase + e.phaseShift)
	dt := e.step

	u := e.mapExpo(t)
	y := 2.0*u - 1.0

	y -= polyBLEP(t, dt)

	if e.dcComp && e.shape != 0 {
		k := e.shape
		meanU := (math.Exp(k) - 1.0 - k) / (k * (math.Exp(k) - 1.0))
		y -= 2.0*meanU - 1.0
	}

	e.phase += dt
	if e.phase >= 1.0 {
		e.phase -= 1.0
	}

	return y, y
}

func (e *Exponential) IsActive() bool {
	return true
}

func (e *Exponential) NoteOn(freq, velocity float64) {
	if freq > 0 {
		e.SetFreq(freq)
	}
	e.Reset()
}

func (e *Exponential) NoteOff(float64) {
}
