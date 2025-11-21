package dsp

import "testing"

type mixerTestCase struct {
	name        string
	withMaster  bool
	withPanGain bool
	mixer       *Mixer
}

func getMixerTestCases() []*mixerTestCase {
	const sr = 44100.0

	cases := []*mixerTestCase{
		{"NoMaster_NoPanGain", false, false, nil},
		{"NoMaster_PanGain", false, true, nil},
		{"Master_NoPanGain", true, false, nil},
		{"Master_PanGain", true, true, nil},
	}

	fillMixer := func(m *Mixer, g, p bool) {
		var gp, pp Param
		if g {
			gp = NewConstParam(0.1)
		}
		if p {
			pp = NewConstParam(0)
		}

		for i := 0; i < 16; i++ {
			osc := NewNoise(NewConstParam(NoiseWhite))
			in := NewInput(osc, gp, pp)
			m.Add(in)
		}
	}

	for _, c := range cases {
		var mg Param
		if c.withMaster {
			mg = NewConstParam(0.5)
		}
		c.mixer = NewMixer(mg, false)
		fillMixer(c.mixer, c.withPanGain, c.withPanGain)
	}

	return cases
}

func BenchmarkMixer_Process(b *testing.B) {
	for _, test := range getMixerTestCases() {
		b.Run(test.name, func(b *testing.B) {
			var block Block
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				test.mixer.Process(&block)
				block.Cycle++
			}
		})
	}
}

func TestMixer_ProcessNoAlloc(t *testing.T) {
	for _, test := range getMixerTestCases() {
		t.Run(test.name, func(t *testing.T) {
			var block Block

			allocs := testing.AllocsPerRun(100, func() {
				test.mixer.Process(&block)
			})

			if allocs != 0 {
				t.Fatalf("expected 0 allocs, got %.1f", allocs)
			}
		})
	}
}
