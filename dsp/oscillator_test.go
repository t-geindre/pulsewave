package dsp

import "testing"

type oscillatorTestCase struct {
	name  string
	shape OscShape
	table *Wavetable
	osc   *Oscillator
}

func getOscTestCases() []*oscillatorTestCase {
	const sr = 44100.0

	cases := []*oscillatorTestCase{
		{"Sine math.sin", ShapeSine, nil, nil},
		{"WT 1024", ShapeTableWave, NewSineWavetable(1024), nil},
		{"WT 512", ShapeTableWave, NewSineWavetable(512), nil},
		{"WT 256", ShapeTableWave, NewSineWavetable(256), nil},
		{"WT 128", ShapeTableWave, NewSineWavetable(128), nil},
		{"WT 64", ShapeTableWave, NewSineWavetable(64), nil},
		{"Square", ShapeSquare, nil, nil},
		{"Triangle", ShapeTriangle, nil, nil},
		{"Noise", ShapeNoise, nil, nil},
		{"Saw", ShapeSaw, nil, nil},
	}

	for _, c := range cases {
		reg := NewShapeRegistry()
		sid := reg.Add(c.shape, c.table)
		c.osc = NewRegOscillator(sr, reg, NewConstParam(sid), NewConstParam(440), nil, nil)
	}

	return cases
}

func BenchmarkOscillator_Resolve(b *testing.B) {
	const sr = 44100.0

	for _, test := range getOscTestCases() {
		b.Run(test.name, func(b *testing.B) {
			s, m := uint64(0), uint64(b.N) // avoid cast in loop

			b.ResetTimer()

			for n := s; n < m; n++ {
				test.osc.Resolve(n)
			}
		})
	}
}

func TestOscillator_ResolveNoAlloc(t *testing.T) {
	for _, test := range getOscTestCases() {
		t.Run(test.name, func(t *testing.T) {
			s, m := uint64(0), uint64(1000)

			allocs := testing.AllocsPerRun(100, func() {
				for n := s; n < m; n++ {
					test.osc.Resolve(n)
				}
			})

			if allocs != 0 {
				t.Errorf("expected 0 allocs, got %v", allocs)
			}
		})
	}
}
