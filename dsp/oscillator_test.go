package dsp

import (
	"testing"
)

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
			osc := NewRegOscillator(sr, reg, sid, NewConstParam(440), nil, nil)

			s, m := uint64(0), uint64(b.N) // avoid cast in loop

			b.ResetTimer()

			for n := s; n < m; n++ {
				osc.Resolve(n)
			}
		})
	}
}
