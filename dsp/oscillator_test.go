package dsp

import (
	"testing"
)

// Simple const param
type ConstParam float64

func (c ConstParam) Resolve(cycle uint64) []float32 {
	buf := make([]float32, BlockSize)
	for i := range buf {
		buf[i] = float32(c)
	}
	return buf
}

func (c ConstParam) SetBase(float32)             {}
func (c ConstParam) ModInputs() *[]ParamModInput { return nil }

// SINE
func BenchmarkOscillator_Sine(b *testing.B) {
	const sr = 44100.0
	osc := NewOscillator(sr, ShapeSine, ConstParam(440), nil, nil)

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		osc.Resolve(uint64(n))
	}
}

// SINE WAVETABLE
func BenchmarkOscillator_SineWaveTable(b *testing.B) {
	const sr = 44100.0
	reg := NewShapeRegistry()
	reg.Set(0, ShapeTableWave, NewSineWavetable(1024))
	osc := NewRegOscillator(sr, reg, 0, ConstParam(440), nil, nil)

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		osc.Resolve(uint64(n))
	}
}
