package dsp

import (
	"testing"
)

// Simple const param, avoid overhead
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

// SINE WAVETABLE
func BenchmarkOscillator_Resolve(b *testing.B) {
	const sr = 44100.0

	for _, test := range []struct {
		name  string
		shape OscShape
		table *Wavetable
	}{
		{
			name:  "Sine math.sin",
			shape: ShapeSine,
		},
		{
			name:  "WT 1024",
			shape: ShapeTableWave,
			table: NewSineWavetable(1024),
		},
		{
			name:  "WT 512",
			shape: ShapeTableWave,
			table: NewSineWavetable(512),
		},
		{
			name:  "WT 256",
			shape: ShapeTableWave,
			table: NewSineWavetable(256),
		},
		{
			name:  "WT 128",
			shape: ShapeTableWave,
			table: NewSineWavetable(128),
		},
		{
			name:  "WT 64",
			shape: ShapeTableWave,
			table: NewSineWavetable(64),
		},
		{
			name:  "Square",
			shape: ShapeSquare,
		},
		{
			name:  "Triangle",
			shape: ShapeTriangle,
		},
		{
			name:  "Noise",
			shape: ShapeNoise,
		},
		{
			name:  "Saw",
			shape: ShapeSaw,
		},
	} {
		b.Run(test.name, func(b *testing.B) {
			reg := NewShapeRegistry()
			sid := reg.Add(test.shape, test.table)
			osc := NewRegOscillator(sr, reg, sid, ConstParam(440), nil, nil)

			s, m := uint64(0), uint64(b.N) // avoid cast in loop

			b.ResetTimer()

			for n := s; n < m; n++ {
				osc.Resolve(n)
			}
		})
	}
}
