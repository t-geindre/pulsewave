package dsp

import "testing"

var noiseTestCases = []struct {
	name string
	kind float32
}{
	{"White", NoiseWhite},
	{"Gaussian", NoiseGaussian},
	{"Pink", NoisePink},
	{"Brown", NoiseBrown},
	{"Blue", NoiseBlue},
}

func BenchmarkNoise_Process(b *testing.B) {
	for _, test := range noiseTestCases {
		b.Run(test.name, func(b *testing.B) {
			var block Block
			b.ResetTimer()

			noise := NewNoise(NewConstParam(test.kind))

			for i := 0; i < b.N; i++ {
				noise.Process(&block)
				block.Cycle++
			}
		})
	}
}

func TestNoise_ResolveNoAlloc(t *testing.T) {
	for _, test := range noiseTestCases {
		t.Run(test.name, func(t *testing.T) {
			noise := NewNoise(NewConstParam(test.kind))
			allocs := testing.AllocsPerRun(100, func() {
				noise.Process(&Block{})
			})

			if allocs != 0 {
				t.Errorf("expected 0 allocs, got %v", allocs)
			}
		})
	}
}
