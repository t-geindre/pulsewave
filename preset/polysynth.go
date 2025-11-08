package preset

import (
	"synth/dsp"
	"synth/midi"
	"synth/msg"
)

type Polysynth struct {
	dsp.Node
	voice               *dsp.PolyVoice
	pitch               dsp.Param
	pOutQueue, pInQueue *msg.Queue
	parameters          map[uint8]dsp.Param
}

func NewPolysynth(SampleRate float64, pInQueue, pOutQueue *msg.Queue) *Polysynth {
	// Parameters map
	preset := NewPreset()
	constZero := dsp.NewConstParam(0)

	// Global pitch bend
	pitchBend := dsp.NewParam(0)

	// Shape registry (uniq for all voices)
	reg := dsp.NewShapeRegistry()
	reg.Add(dsp.ShapeTableWave, dsp.NewSineWavetable(1024))
	reg.Add(dsp.ShapeSquare)
	reg.Add(dsp.ShapeSaw)
	reg.Add(dsp.ShapeTriangle)
	reg.Add(dsp.ShapeNoise)

	// Voice factory
	voiceFact := func() *dsp.Voice {
		// Base frequency param (uniq per voice)
		freq := dsp.NewSmoothedParam(SampleRate, 440, .001)
		pitchMod := dsp.NewParam(0)
		pitch := dsp.NewTunerParam(dsp.NewTunerParam(freq, pitchBend), pitchMod)

		pitchLfo := dsp.NewRegOscillator(SampleRate, reg, preset[PitchLfoShape], preset[PitchLfoFreq], preset[PitchLfoPhase], nil)
		pitchAdsr := dsp.NewADSR(SampleRate, preset[PitchAdsrAttack], preset[PitchAdsrDecay], preset[PitchAdsrSustain], preset[PitchAdsrRelease])

		*pitchMod.ModInputs() = append(*pitchMod.ModInputs(),
			dsp.NewModInput(pitchLfo, NewParamSkipper(preset[PitchLfoAmount], constZero, preset[PitchLfoOnOff]), nil),
			dsp.NewModInput(pitchAdsr, NewParamSkipper(preset[PitchAdsrAmount], constZero, preset[PitchAdsrOnOff]), nil),
		)

		// Oscillator factory todo implement phase
		oscFact := func(ph, dt dsp.Param) dsp.Node {
			// Mixer, registry
			mixer := dsp.NewMixer(dsp.NewParam(1), false)
			ft := dsp.NewTunerParam(pitch, dt)

			// 0
			mixer.Add(dsp.NewInput(
				dsp.NewRegOscillator(SampleRate, reg, preset[Osc0Shape], dsp.NewTunerParam(ft, preset[Osc0Detune]), ph, preset[Osc0Pw]),
				preset[Osc0Gain],
				dsp.NewParam(0),
			))

			// 1
			mixer.Add(dsp.NewInput(
				dsp.NewRegOscillator(SampleRate, reg, preset[Osc1Shape], dsp.NewTunerParam(ft, preset[Osc1Detune]), ph, preset[Osc1Pw]),
				preset[Osc1Gain],
				dsp.NewParam(0),
			))

			// 2
			mixer.Add(dsp.NewInput(
				dsp.NewRegOscillator(SampleRate, reg, preset[Osc2Shape], dsp.NewTunerParam(ft, preset[Osc2Detune]), ph, preset[Osc2Pw]),
				preset[Osc2Gain],
				dsp.NewParam(0),
			))

			return mixer
		}

		// Unison
		unison := dsp.NewUnison(dsp.UnisonOpts{
			SampleRate:   SampleRate,
			NumVoices:    preset[UnisonVoices],
			Factory:      oscFact,
			PanSpread:    preset[UnisonPanSpread],
			PhaseSpread:  preset[UnisonPhaseSpread],
			DetuneSpread: preset[UnisonDetuneSpread],
			CurveGamma:   preset[UnisonCurveGamma],
		})
		unisonSkip := NewNodeSkipper(
			unison,
			oscFact(dsp.NewParam(0), dsp.NewParam(0)), // Unique voice
			preset[UnisonOnOff],
		)

		// LPF
		cutoffLfo := dsp.NewRegOscillator(SampleRate, reg, preset[LpfLfoShape], preset[LpfLfoFreq], preset[LpfLfoPhase], nil)
		cutoffAdsr := dsp.NewADSR(SampleRate, preset[LpfAdsrAttack], preset[LpfAdsrDecay], preset[LpfAdsrSustain], preset[LpfAdsrRelease])

		cutoff := dsp.NewParam(0)
		*cutoff.ModInputs() = append(*cutoff.ModInputs(),
			dsp.NewModInput(preset[LPFCutoff], dsp.NewParam(1), nil),
			dsp.NewModInput(cutoffLfo, NewParamSkipper(preset[LpfLfoAmount], constZero, preset[LpfLfoOnOff]), nil),
			dsp.NewModInput(cutoffAdsr, NewParamSkipper(preset[LpfAdsrAmount], constZero, preset[LpfAdsrOnOff]), nil),
		)

		lpf := dsp.NewLowPassSVF(SampleRate, unisonSkip, cutoff, preset[LPFResonance])
		lpfSkip := NewNodeSkipper(lpf, unisonSkip, preset[LPFOnOff])

		// Amplitude envelope
		gainAdsr := dsp.NewADSR(SampleRate, preset[AmpEnvAttack], preset[AmpEnvDecay], preset[AmpEnvSustain], preset[AmpEnvRelease])

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
	poly := dsp.NewPolyVoice(8, voiceFact)

	// Delay with skipper
	delay := dsp.NewFeedbackDelay(SampleRate, 2.0, poly, preset[FBDelayParam], preset[FBFeedBack], preset[FBMix], preset[FBTone])
	delaySkip := NewNodeSkipper(delay, poly, preset[FBOnOff])

	return &Polysynth{
		Node:       delaySkip,
		voice:      poly,
		pitch:      pitchBend,
		pInQueue:   pInQueue,
		pOutQueue:  pOutQueue,
		parameters: preset,
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

func (p *Polysynth) Process(b *dsp.Block) {
	p.pInQueue.Drain(10, p.HandleMessage)
	p.Node.Process(b)
}

func (p *Polysynth) HandleMessage(m msg.Message) {
	switch m.Kind {
	case ParamPullAllKind:
		p.PublishParameters()
	case ParamUpdateKind:
		if param, ok := p.parameters[m.Key]; ok {
			param.SetBase(m.ValF)
		}
	case midi.NoteOnKind:
		// Todo handle vel properly with LUT (precalculated curve)
		p.voice.NoteOn(int(m.Key), float32(m.Val8)/127)
	case midi.NoteOffKind:
		p.voice.NoteOff(int(m.Key))
	case midi.PitchBendKind:
		rel := float32(0)
		if m.Val16 >= 128 || m.Val16 <= -128 {
			rel = float32(m.Val16) / 8192.0 * 2.0 // 2 semitones range
		}
		p.pitch.SetBase(rel)
	}
}

func (p *Polysynth) PublishParameters() {
	for key, param := range p.parameters {
		p.pOutQueue.TryWrite(msg.Message{
			Source: AudioSource,
			Kind:   ParamUpdateKind,
			Key:    key,
			ValF:   param.GetBase(),
		})
	}
}
