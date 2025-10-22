package effect

import (
	"math"
	"synth/audio"
	"time"
)

// Feedback is a feedback delay effect.
type Feedback struct {
	src audio.Source
	sr  float64

	buf       [2][]float64 // stereo buffer
	writePos  int
	delaySamp int

	feedback float64 // 0..1
	wet      float64 // 0..1
	dry      float64 // 0..1

	damp float64 // 0..1
	z    [2]float64

	energyBuf []float64
	energySum float64
	energyIdx int
}

func NewFeedback(sampleRate float64, src audio.Source) *Feedback {
	f := &Feedback{
		src:       src,
		sr:        sampleRate,
		feedback:  0.4,
		wet:       0.3,
		dry:       1.0 - 0.3,
		damp:      0.0,
		energyBuf: make([]float64, int(math.Max(64.0, sampleRate/50.0))),
	}

	f.SetDelay(400 * time.Millisecond)

	return f
}

func (f *Feedback) ensureBuf(min int) {
	for i, buf := range f.buf {
		if len(buf) >= min {
			continue
		}
		n := 1
		for n < min {
			n <<= 1
		}
		f.buf[i] = make([]float64, n)
		f.writePos = 0
	}
}

func (f *Feedback) SetDelay(t time.Duration) {
	f.delaySamp = int(math.Max(0, math.Round(float64(t)*f.sr/float64(time.Second))))
	f.ensureBuf(f.delaySamp + 1)
}

func (f *Feedback) SetFeedback(g float64) {
	if g < 0 {
		g = 0
	}
	if g > 0.98 {
		g = 0.98
	}
	f.feedback = g
}
func (f *Feedback) SetMix(wet float64) {
	if wet < 0 {
		wet = 0
	}
	if wet > 1 {
		wet = 1
	}
	f.wet = wet
	f.dry = 1 - wet
}
func (f *Feedback) SetDamping(amount float64) {
	if amount < 0 {
		amount = 0
	}
	if amount > 1 {
		amount = 1
	}
	f.damp = amount
}

func (f *Feedback) NextValue() (float64, float64) {
	inL, inR := f.src.NextValue()

	readPos := f.writePos - f.delaySamp
	if readPos < 0 {
		readPos += len(f.buf[0])
	}
	delL := f.buf[0][readPos]
	delR := f.buf[1][readPos]

	outL := f.dry*inL + f.wet*delL
	outR := f.dry*inR + f.wet*delR

	fbL := delL
	fbR := delR

	// apply damping
	if f.damp > 0 {
		a := f.damp * 0.5
		f.z[0] += a * (delL - f.z[0])
		f.z[1] += a * (delR - f.z[1])
		fbL = f.z[0]
		fbR = f.z[1]
	}

	xL := inL + f.feedback*fbL
	xR := inR + f.feedback*fbR

	// Energy ring buffer update
	old := f.energyBuf[f.energyIdx]
	e := xL*xL + xR*xR
	f.energyBuf[f.energyIdx] = e
	f.energySum += e - old
	f.energyIdx++
	if f.energyIdx >= len(f.energyBuf) {
		f.energyIdx = 0
	}

	// Delay buffer write
	f.buf[0][f.writePos] = xL
	f.buf[1][f.writePos] = xR

	f.writePos++
	if f.writePos >= len(f.buf[0]) {
		f.writePos = 0
	}

	return outL, outR
}

func (f *Feedback) IsActive() bool {
	if f.src.IsActive() {
		return true
	}

	eps := 1e-8 * float64(len(f.energyBuf))
	return f.energySum > eps
}

func (f *Feedback) Reset() {
	f.src.Reset()
	f.z = [2]float64{0, 0}

	for i := range f.buf {
		for j := range f.buf[i] {
			f.buf[i][j] = 0
		}
	}
	f.writePos = 0

	for i := range f.energyBuf {
		f.energyBuf[i] = 0
	}
	f.energySum = 0
	f.energyIdx = 0
}

func (f *Feedback) NoteOn(freq, velocity float64) {
	f.src.NoteOn(freq, velocity)
}

func (f *Feedback) NoteOff() {
	f.src.NoteOff()
}
