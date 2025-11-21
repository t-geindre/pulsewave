package preset

import (
	"synth/dsp"
	"synth/msg"
)

type Polysynth struct {
	dsp.Node
	voice      *dsp.PolyVoice
	pitch      dsp.Param
	parameters map[uint8]dsp.Param
	messenger  *msg.Messenger
}

func NewPolysynth(SampleRate float64) *Polysynth {
	// Parameters map
	preset := NewPreset()
	constZero := dsp.NewConstParam(0)

	// Global pitch bend
	pitchBend := dsp.NewSmoothedParam(SampleRate, 0, dsp.NewConstParam(.01))

	// Shape registry (uniq for all voices)
	reg := dsp.NewShapeRegistry()
	reg.Add(dsp.ShapeTableWave, dsp.NewSineWavetable(1024))
	reg.Add(dsp.ShapeSquare)
	reg.Add(dsp.ShapeSaw)
	reg.Add(dsp.ShapeTriangle)

	// Voice factory / 3 osc
	voiceFact := func() *dsp.Voice {
		// Base frequency param (uniq per voice)
		freq := dsp.NewSmoothedParam(SampleRate, 440, preset.Params[VoicesPitchGlide])
		pitchMod := dsp.NewParam(0)
		pitch := dsp.NewTunerParam(dsp.NewTunerParam(freq, pitchBend), pitchMod)

		pitchLfo := dsp.NewRegOscillator(SampleRate, reg, preset.Params[PitchLfoShape], preset.Params[PitchLfoFreq], preset.Params[PitchLfoPhase], nil)
		pitchAdsr := dsp.NewADSR(SampleRate, preset.Params[PitchAdsrAttack], preset.Params[PitchAdsrDecay], preset.Params[PitchAdsrSustain], preset.Params[PitchAdsrRelease])

		*pitchMod.ModInputs() = append(*pitchMod.ModInputs(),
			dsp.NewModInput(pitchLfo, NewParamSkipper(preset.Params[PitchLfoAmount], constZero, preset.Params[PitchLfoOnOff]), nil),
			dsp.NewModInput(pitchAdsr, NewParamSkipper(preset.Params[PitchAdsrAmount], constZero, preset.Params[PitchAdsrOnOff]), nil),
		)

		// Oscillator factory
		oscFact := func(ph, dt dsp.Param) dsp.Node {
			// Mixer, registry
			mixer := dsp.NewMixer(nil, false)
			ft := dsp.NewTunerParam(pitch, dt)

			// 0
			ph0 := dsp.NewParam(0)
			*ph0.ModInputs() = append(*ph0.ModInputs(),
				dsp.NewModInput(ph, dsp.NewConstParam(1), nil),
				dsp.NewModInput(preset.Params[Osc0Phase], dsp.NewConstParam(1), nil),
			)
			mixer.Add(dsp.NewInput(
				dsp.NewRegOscillator(SampleRate, reg, preset.Params[Osc0Shape], dsp.NewTunerParam(ft, preset.Params[Osc0Detune]), ph0, preset.Params[Osc0Pw]),
				preset.Params[Osc0Gain],
				dsp.NewParam(0),
			))

			// 1
			ph1 := dsp.NewParam(0)
			*ph1.ModInputs() = append(*ph1.ModInputs(),
				dsp.NewModInput(ph, dsp.NewConstParam(1), nil),
				dsp.NewModInput(preset.Params[Osc1Phase], dsp.NewConstParam(1), nil),
			)
			mixer.Add(dsp.NewInput(
				dsp.NewRegOscillator(SampleRate, reg, preset.Params[Osc1Shape], dsp.NewTunerParam(ft, preset.Params[Osc1Detune]), ph1, preset.Params[Osc1Pw]),
				preset.Params[Osc1Gain],
				dsp.NewParam(0),
			))

			// 2
			ph2 := dsp.NewParam(0)
			*ph2.ModInputs() = append(*ph2.ModInputs(),
				dsp.NewModInput(ph, dsp.NewConstParam(1), nil),
				dsp.NewModInput(preset.Params[Osc2Phase], dsp.NewConstParam(1), nil),
			)
			mixer.Add(dsp.NewInput(
				dsp.NewRegOscillator(SampleRate, reg, preset.Params[Osc2Shape], dsp.NewTunerParam(ft, preset.Params[Osc2Detune]), ph2, preset.Params[Osc2Pw]),
				preset.Params[Osc2Gain],
				dsp.NewParam(0),
			))

			return mixer
		}

		// Unison
		unison := dsp.NewUnison(dsp.UnisonOpts{
			SampleRate:   SampleRate,
			NumVoices:    preset.Params[UnisonVoices],
			Factory:      oscFact,
			PanSpread:    preset.Params[UnisonPanSpread],
			PhaseSpread:  preset.Params[UnisonPhaseSpread],
			DetuneSpread: preset.Params[UnisonDetuneSpread],
			CurveGamma:   preset.Params[UnisonCurveGamma],
		})
		unisonSkip := NewNodeSkipper(
			unison,
			oscFact(dsp.NewParam(0), dsp.NewParam(0)), // Unique voice
			preset.Params[UnisonOnOff],
		)

		// Global mixer
		globalMix := dsp.NewMixer(nil, false)
		globalMix.Add(dsp.NewInput(unisonSkip, nil, nil))

		// Noise oscillator
		noiseOsc := dsp.NewNoise(preset.Params[NoiseType])
		globalMix.Add(dsp.NewInput(noiseOsc, preset.Params[NoiseGain], nil))

		// Sub oscillator
		subOsc := dsp.NewRegOscillator(SampleRate, reg, preset.Params[SubOscShape], dsp.NewTunerParam(pitch, preset.Params[SubOscTranspose]), nil, nil)
		globalMix.Add(dsp.NewInput(subOsc, preset.Params[SubOscGain], nil))

		// LPF
		cutoffLfo := dsp.NewRegOscillator(SampleRate, reg, preset.Params[LpfLfoShape], preset.Params[LpfLfoFreq], preset.Params[LpfLfoPhase], nil)
		cutoffAdsr := dsp.NewADSR(SampleRate, preset.Params[LpfAdsrAttack], preset.Params[LpfAdsrDecay], preset.Params[LpfAdsrSustain], preset.Params[LpfAdsrRelease])

		cutoff := dsp.NewParam(0)
		*cutoff.ModInputs() = append(*cutoff.ModInputs(),
			dsp.NewModInput(preset.Params[LPFCutoff], dsp.NewParam(1), nil),
			dsp.NewModInput(cutoffLfo, NewParamSkipper(preset.Params[LpfLfoAmount], constZero, preset.Params[LpfLfoOnOff]), nil),
			dsp.NewModInput(cutoffAdsr, NewParamSkipper(preset.Params[LpfAdsrAmount], constZero, preset.Params[LpfAdsrOnOff]), nil),
		)

		lpf := dsp.NewLowPassSVF(SampleRate, globalMix, cutoff, preset.Params[LPFResonance])
		lpfSkip := NewNodeSkipper(lpf, globalMix, preset.Params[LPFOnOff])

		// Amplitude envelope
		gainAdsr := dsp.NewADSR(SampleRate, preset.Params[AmpEnvAttack], preset.Params[AmpEnvDecay], preset.Params[AmpEnvSustain], preset.Params[AmpEnvRelease])

		gain := dsp.NewParam(0)
		*gain.ModInputs() = append(*gain.ModInputs(),
			dsp.NewModInput(gainAdsr, dsp.NewConstParam(1.0), nil),
		)

		vca := dsp.NewVca(lpfSkip, gain)

		// Voice
		voice := dsp.NewVoice(vca, freq, gainAdsr, cutoffLfo, cutoffAdsr, pitchLfo, pitchAdsr)

		return voice
	}

	// Polyphonic voice
	poly := dsp.NewPolyVoice(16, preset.Params[VoicesActive], preset.Params[VoicesStealMode], voiceFact)

	// Delay with skipper
	delay := dsp.NewFeedbackDelay(SampleRate, 2.0, poly, preset.Params[FBDelayParam], preset.Params[FBFeedBack], preset.Params[FBMix], preset.Params[FBTone])
	delaySkip := NewNodeSkipper(delay, poly, preset.Params[FBOnOff])

	return &Polysynth{
		Node:       delaySkip,
		voice:      poly,
		pitch:      pitchBend,
		parameters: preset.Params,
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

func (p *Polysynth) SetParam(key uint8, val float32) {
	if param, ok := p.parameters[key]; ok {
		param.SetBase(val)
	}
}

func (p *Polysynth) LoadPreset(preset *Preset) {
	for key, param := range preset.Params {
		p.SetParam(key, param.GetBase())
	}
}

func (p *Polysynth) HydratePreset(preset *Preset) *Preset {
	for key, param := range p.parameters {
		preset.Params[key] = dsp.NewParam(param.GetBase())
	}

	return preset
}

func (p *Polysynth) AllNotesOff() {
	p.voice.AllNotesOff()
}
