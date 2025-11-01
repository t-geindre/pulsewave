package dsp

import "testing"

func BenchmarkMixer_Process(b *testing.B) {
	const sr = 44100.0

	fillMixer := func(m *Mixer, g, p bool) {
		var gp, pp Param
		if g {
			gp = NewConstParam(0.1)
		}
		if p {
			pp = NewConstParam(0)
		}

		for i := 0; i < 16; i++ {
			osc := NewOscillator(sr, ShapeNoise, NewConstParam(440), nil, nil) // Noise = fastest
			in := NewInput(osc, gp, pp)
			m.Add(in)
		}
	}

	for _, test := range []struct {
		name        string
		withMaster  bool
		withPanGain bool
	}{
		{"NoMaster_NoPanGain", false, false},
		{"NoMaster_PanGain", false, true},
		{"Master_NoPanGain", true, false},
		{"Master_PanGain", true, true},
	} {
		b.Run(test.name, func(b *testing.B) {
			mixer := NewMixer(nil, false)
			if test.withMaster {
				mixer.MasterGain = NewConstParam(0.5)
			}
			fillMixer(mixer, test.withPanGain, test.withPanGain)

			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				var block Block
				mixer.Process(&block)
			}
		})
	}
}
