package preset

import (
	"synth/dsp"
	"time"
)

type Polysynth struct {
	dsp.Node
	voice *dsp.PolyVoice
	pitch dsp.Param
}

func NewPolysynth(SampleRate float64) *Polysynth {
	// Main pitch bend param
	pitchBend := dsp.NewParam(0)

	// Shape registry (uniq for all voices)
	reg := dsp.NewShapeRegistry()
	reg.Set(0, dsp.ShapeSaw)
	reg.Set(1, dsp.ShapeTriangle)
	reg.Set(2, dsp.ShapeTableWave, dsp.NewSineWavetable(1024))

	// Voice factory
	voiceFact := func() *dsp.Voice {
		// Base frequency param (uniq per voice)
		freq := dsp.NewSmoothedParam(SampleRate, 440, .001)
		pitch := dsp.NewTunerParam(freq, pitchBend)

		// Oscillator factory
		oscFact := func(ph, dt dsp.Param) dsp.Node {
			// Mixer, registry
			mixer := dsp.NewMixer(dsp.NewParam(1), false)
			ft := dsp.NewTunerParam(pitch, dt)

			// 0
			mixer.Add(dsp.NewInput(
				dsp.NewRegOscillator(SampleRate, reg, 0, ft, ph, nil),
				dsp.NewParam(.33),
				dsp.NewParam(0),
			))

			// 1
			mixer.Add(dsp.NewInput(
				dsp.NewRegOscillator(SampleRate, reg, 1, dsp.NewTunerParam(ft, dsp.NewParam(-12)), ph, nil),
				dsp.NewParam(.33),
				dsp.NewParam(0),
			))

			// 2
			mixer.Add(dsp.NewInput(
				dsp.NewRegOscillator(SampleRate, reg, 2, dsp.NewTunerParam(ft, dsp.NewParam(+24)), ph, nil),
				dsp.NewParam(0.15),
				dsp.NewParam(0),
			))
			return mixer
		}

		// Unison
		unison := dsp.NewUnison(dsp.UnisonOpts{
			SampleRate:   SampleRate,
			NumVoices:    8,
			Factory:      oscFact,
			PanSpread:    dsp.NewParam(1.0),
			PhaseSpread:  dsp.NewParam(.1),
			DetuneSpread: dsp.NewParam(12.0),
			CurveGamma:   dsp.NewParam(1.5),
		})

		// LPF
		cutoff := dsp.NewSmoothedParam(SampleRate, 800, 0.005)
		reson := dsp.NewParam(1)
		lpf := dsp.NewLowPassSVF(SampleRate, unison, cutoff, reson)

		ctModRateAdsr := dsp.NewADSR(SampleRate, time.Millisecond*5, time.Millisecond*50, 0, time.Millisecond*100)
		*cutoff.ModInputs() = append(*cutoff.ModInputs(), dsp.NewModInput(ctModRateAdsr, 4000, nil))

		ctModRateOsc := dsp.NewOscillator(SampleRate, dsp.ShapeSine, dsp.NewParam(.5), dsp.NewParam(1), nil)
		*cutoff.ModInputs() = append(*cutoff.ModInputs(), dsp.NewModInput(ctModRateOsc, 300, nil))

		// Voice
		adsr := dsp.NewADSR(SampleRate, time.Millisecond*10, time.Millisecond*800, .9, time.Millisecond*100)
		voice := dsp.NewVoice(lpf, freq, adsr, ctModRateAdsr, ctModRateOsc)

		return voice
	}

	// Polyphonic voice
	poly := dsp.NewPolyVoice(8, voiceFact)

	// Delay
	delay := dsp.NewFeedbackDelay(
		SampleRate,
		2.0,
		poly,
		dsp.NewParam(0.35), // delay time in seconds
		dsp.NewParam(0.3),  // feedback amount (0-1)
		dsp.NewParam(0.2),  // mix
		dsp.NewParam(2000), // mix amount (0-1)
	)

	return &Polysynth{
		Node:  delay,
		voice: poly,
		pitch: pitchBend,
	}
}

func (p *Polysynth) NoteOn(key int, vel float32) {
	p.voice.NoteOn(key, vel)
}

func (p *Polysynth) NoteOff(key int) {
	p.voice.NoteOff(key)
}

func (p *Polysynth) SetPitchBend(semiTones float32) {
	p.pitch.SetBase(semiTones)
}
