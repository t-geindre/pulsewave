package effect

import (
	"math"
	"synth/audio"
	"time"
)

type FeedbackDelay struct {
	src audio.Source
	sr  float64

	buf       []float64
	writePos  int
	delaySamp int

	feedback float64 // 0..1
	wet      float64 // 0..1
	dry      float64 // 0..1

	damp float64 // 0..1
	z    float64
}

func NewFeedbackDelay(sampleRate float64, src audio.Source) *FeedbackDelay {
	fd := &FeedbackDelay{
		src:      src,
		sr:       sampleRate,
		feedback: 0.4,
		wet:      0.3,
		dry:      1.0 - 0.3,
		damp:     0.0,
	}
	fd.SetDelay(400)
	return fd
}

func (f *FeedbackDelay) ensureBuf(min int) {
	if len(f.buf) >= min {
		return
	}

	n := 1
	for n < min {
		n <<= 1
	}
	f.buf = make([]float64, n)
	f.writePos = 0
}

func (f *FeedbackDelay) SetDelay(t time.Duration) {
	ms := float64(t.Microseconds()) / 1000.0
	if ms < 0 {
		ms = 0
	}
	samp := int(math.Round(ms * f.sr / 1000.0))
	if samp < 1 {
		samp = 1
	}
	f.delaySamp = samp
	f.ensureBuf(f.delaySamp + 1)
}

func (f *FeedbackDelay) SetDelaySamples(n int) {
	if n < 1 {
		n = 1
	}
	f.delaySamp = n
	f.ensureBuf(f.delaySamp + 1)
}

func (f *FeedbackDelay) SetFeedback(g float64) {
	if g < 0 {
		g = 0
	}
	if g > 0.98 {
		g = 0.98
	}
	f.feedback = g
}
func (f *FeedbackDelay) SetMix(wet float64) {
	if wet < 0 {
		wet = 0
	}
	if wet > 1 {
		wet = 1
	}
	f.wet = wet
	f.dry = 1 - wet
}
func (f *FeedbackDelay) SetDamping(amount float64) {
	if amount < 0 {
		amount = 0
	}
	if amount > 1 {
		amount = 1
	}
	f.damp = amount
}

func (f *FeedbackDelay) NextSample() float64 {
	in := f.src.NextSample()

	readPos := f.writePos - f.delaySamp
	if readPos < 0 {
		readPos += len(f.buf)
	}
	del := f.buf[readPos]

	out := f.dry*in + f.wet*del

	fb := del
	if f.damp > 0 {
		a := f.damp * 0.5
		f.z += a * (del - f.z)
		fb = f.z
	}

	f.buf[f.writePos] = in + f.feedback*fb

	f.writePos++
	if f.writePos >= len(f.buf) {
		f.writePos = 0
	}

	return out
}
