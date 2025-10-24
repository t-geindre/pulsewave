package effect

import (
	"math"
	"synth/audio"
)

type HighPassFilter struct {
	src   audio.Source
	sr    float64
	cutHz float64
	Q     float64 // 0.1..~10 ; 0.707 = Butterworth

	// Coefficients communs L/R
	b0, b1, b2 float64
	a1, a2     float64

	// États par canal: 0 = L, 1 = R
	x1 [2]float64
	x2 [2]float64
	y1 [2]float64
	y2 [2]float64
}

func NewHighPassFilter(sampleRate float64, src audio.Source) *HighPassFilter {
	hp := &HighPassFilter{
		src:   src,
		sr:    sampleRate,
		cutHz: 200.0,
		Q:     0.707, // Butterworth
	}
	hp.recalc()
	return hp
}

func (h *HighPassFilter) SetCutoffHz(hz float64) {
	if hz < 10 {
		hz = 10
	}
	max := h.sr * 0.49
	if hz > max {
		hz = max
	}
	h.cutHz = hz
	h.recalc()
}

func (h *HighPassFilter) SetQ(q float64) {
	if q < 0.1 {
		q = 0.1
	}
	h.Q = q
	h.recalc()
}

func (h *HighPassFilter) Reset() {
	h.x1 = [2]float64{}
	h.x2 = [2]float64{}
	h.y1 = [2]float64{}
	h.y2 = [2]float64{}
	h.src.Reset()
}

func (h *HighPassFilter) NextValue() (float64, float64) {
	xL, xR := h.src.NextValue()

	yL := h.processSample(0, xL)
	yR := h.processSample(1, xR)

	// flush-to-zero anti-denormals
	if yL > -1e-20 && yL < 1e-20 {
		yL = 0
	}
	if yR > -1e-20 && yR < 1e-20 {
		yR = 0
	}

	return yL, yR
}

func (h *HighPassFilter) processSample(ch int, x float64) float64 {
	// Biquad Direct Form I
	y := h.b0*x + h.b1*h.x1[ch] + h.b2*h.x2[ch] - h.a1*h.y1[ch] - h.a2*h.y2[ch]

	// shift states
	h.x2[ch] = h.x1[ch]
	h.x1[ch] = x
	h.y2[ch] = h.y1[ch]
	h.y1[ch] = y

	return y
}

func (h *HighPassFilter) recalc() {
	omega := 2 * math.Pi * h.cutHz / h.sr
	sin := math.Sin(omega)
	cos := math.Cos(omega)
	alpha := sin / (2 * h.Q)

	// High-pass RBJ
	b0 := (1 + cos) / 2
	b1 := -(1 + cos)
	b2 := (1 + cos) / 2
	a0 := 1 + alpha
	a1 := -2 * cos
	a2 := 1 - alpha

	// Normalisation a0=1
	h.b0 = b0 / a0
	h.b1 = b1 / a0
	h.b2 = b2 / a0
	h.a1 = a1 / a0
	h.a2 = a2 / a0
}

func (h *HighPassFilter) IsActive() bool {
	return h.src.IsActive()
}

func (h *HighPassFilter) NoteOn(freq, velocity float64) {
	h.src.NoteOn(freq, velocity)
}

func (h *HighPassFilter) NoteOff() {
	h.src.NoteOff()
}
