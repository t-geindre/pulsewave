package effect

import (
	"math"
	"synth/audio"
)

type LowPassFilter struct {
	src   audio.Source
	sr    float64
	cutHz float64
	Q     float64 // 0.1..~10 ; 0.707 = Butterworth

	// Mono coefficients
	b0, b1, b2 float64
	a1, a2     float64

	// L/R states
	x1 [2]float64
	x2 [2]float64
	y1 [2]float64
	y2 [2]float64
}

func NewLowPassFilter(sampleRate float64, src audio.Source) *LowPassFilter {
	lp := &LowPassFilter{
		src:   src,
		sr:    sampleRate,
		cutHz: 1000.0,
		Q:     0.707, // Butterworth
	}
	lp.recalc()
	return lp
}

func (l *LowPassFilter) SetCutoffHz(hz float64) {
	if hz < 20 {
		hz = 20
	}
	max := l.sr * 0.49
	if hz > max {
		hz = max
	}
	l.cutHz = hz
	l.recalc()
}

func (l *LowPassFilter) SetQ(q float64) {
	if q < 0.1 {
		q = 0.1
	}
	l.Q = q
	l.recalc()
}

func (l *LowPassFilter) Reset() {
	l.x1 = [2]float64{}
	l.x2 = [2]float64{}
	l.y1 = [2]float64{}
	l.y2 = [2]float64{}
	l.src.Reset()
}

func (l *LowPassFilter) NextValue() (float64, float64) {
	xL, xR := l.src.NextValue()

	yL := l.processSample(0, xL)
	yR := l.processSample(1, xR)

	if yL > -1e-20 && yL < 1e-20 {
		yL = 0
	}
	if yR > -1e-20 && yR < 1e-20 {
		yR = 0
	}

	return yL, yR
}

func (l *LowPassFilter) processSample(ch int, x float64) float64 {
	// biquad direct form I (RBJ)
	y := l.b0*x + l.b1*l.x1[ch] + l.b2*l.x2[ch] - l.a1*l.y1[ch] - l.a2*l.y2[ch]

	// shift states
	l.x2[ch] = l.x1[ch]
	l.x1[ch] = x
	l.y2[ch] = l.y1[ch]
	l.y1[ch] = y

	return y
}

func (l *LowPassFilter) recalc() {
	omega := 2 * math.Pi * l.cutHz / l.sr
	sin := math.Sin(omega)
	cos := math.Cos(omega)
	alpha := sin / (2 * l.Q)

	// Low-pass RBJ
	b0 := (1 - cos) / 2
	b1 := 1 - cos
	b2 := (1 - cos) / 2
	a0 := 1 + alpha
	a1 := -2 * cos
	a2 := 1 - alpha

	// normalization a0=1
	l.b0 = b0 / a0
	l.b1 = b1 / a0
	l.b2 = b2 / a0
	l.a1 = a1 / a0
	l.a2 = a2 / a0
}

func (l *LowPassFilter) IsActive() bool {
	return l.src.IsActive()
}

func (l *LowPassFilter) NoteOn(freq, velocity float64) {
	l.src.NoteOn(freq, velocity)
}

func (l *LowPassFilter) NoteOff() {
	l.src.NoteOff()
}
