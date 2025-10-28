package main

import (
	"synth/audio"
	"synth/dsp"
	"time"
)

func main() {
	const SampleRate = 44100

	// Mixer
	gain := dsp.NewConstParam(1)
	mixer := dsp.NewMixer(gain, true)

	// Oscillators
	reg := dsp.NewShapeRegistry()
	reg.Set(0, dsp.Sine)
	reg.Set(1, dsp.Triangle)
	reg.Set(2, dsp.Saw)

	freq := dsp.NewSmoothedParam(SampleRate, 440, .001)
	for i := 0; i < 3; i++ {
		oscillator := dsp.NewOscillator(SampleRate, reg, i, freq, nil, nil)
		input := &dsp.Input{
			Src:  oscillator,
			Gain: dsp.NewConstParam(0.3),
			Pan:  dsp.NewConstParam(-0.5 + float32(i)*0.5),
		}
		mixer.Add(input)
	}

	// Freq mod
	fmRate := dsp.NewOscillator(SampleRate, reg, 0, dsp.NewConstParam(1), nil, nil)
	*freq.ModInputs() = append(*freq.ModInputs(), dsp.NewModInput(fmRate, 100.0, nil))

	// Player
	p := audio.NewPlayer(SampleRate, mixer)
	p.SetBufferSize(time.Millisecond * 20)

	for {
		time.Sleep(time.Second * 2)
	}
}
