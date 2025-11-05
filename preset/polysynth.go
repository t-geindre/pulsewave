package preset

import (
	"synth/dsp"
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
	parameters := make(map[uint8]dsp.Param)

	// Oscillator parameters
	parameters[Osc0Shape] = dsp.NewParam(2) // Saw
	parameters[Osc1Shape] = dsp.NewParam(0)
	parameters[Osc2Shape] = dsp.NewParam(0)

	parameters[Osc0Gain] = dsp.NewParam(1.0)
	parameters[Osc1Gain] = dsp.NewParam(0)
	parameters[Osc2Gain] = dsp.NewParam(0)

	parameters[Osc0Detune] = dsp.NewParam(0)
	parameters[Osc1Detune] = dsp.NewParam(0)
	parameters[Osc2Detune] = dsp.NewParam(0)

	// Amplitude envelope parameters
	parameters[AmpEnvAttack] = dsp.NewParam(0.01)
	parameters[AmpEnvDecay] = dsp.NewParam(0.1)
	parameters[AmpEnvSustain] = dsp.NewParam(0.9)
	parameters[AmpEnvRelease] = dsp.NewParam(0.01)

	// Unison parameters (all voices share)
	parameters[UnisonPanSpread] = dsp.NewParam(0)
	parameters[UnisonPhaseSpread] = dsp.NewParam(0)
	parameters[UnisonDetuneSpread] = dsp.NewParam(0)
	parameters[UnisonCurveGamma] = dsp.NewParam(1)
	parameters[UnisonVoices] = dsp.NewParam(1)

	// Low pass filter parameters
	parameters[LPFCutoff] = dsp.NewSmoothedParam(SampleRate, 8000, 0.005) // Modulated needs smoothing
	parameters[LPFResonance] = dsp.NewParam(1)

	// Main pitch bend param
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
		pitch := dsp.NewTunerParam(freq, pitchBend)

		// Oscillator factory
		oscFact := func(ph, dt dsp.Param) dsp.Node {
			// Mixer, registry
			mixer := dsp.NewMixer(dsp.NewParam(1), false)
			ft := dsp.NewTunerParam(pitch, dt)

			// 0
			mixer.Add(dsp.NewInput(
				dsp.NewRegOscillator(SampleRate, reg, parameters[Osc0Shape], dsp.NewTunerParam(ft, parameters[Osc0Detune]), ph, nil),
				parameters[Osc0Gain],
				dsp.NewParam(0),
			))

			// 1
			mixer.Add(dsp.NewInput(
				dsp.NewRegOscillator(SampleRate, reg, parameters[Osc1Shape], dsp.NewTunerParam(ft, parameters[Osc1Detune]), ph, nil),
				parameters[Osc1Gain],
				dsp.NewParam(0),
			))

			// 2
			mixer.Add(dsp.NewInput(
				dsp.NewRegOscillator(SampleRate, reg, parameters[Osc2Shape], dsp.NewTunerParam(ft, parameters[Osc2Detune]), ph, nil),
				parameters[Osc2Gain],
				dsp.NewParam(0),
			))
			return mixer
		}

		// Unison
		unison := dsp.NewUnison(dsp.UnisonOpts{
			SampleRate:   SampleRate,
			NumVoices:    parameters[UnisonVoices],
			Factory:      oscFact,
			PanSpread:    parameters[UnisonPanSpread],
			PhaseSpread:  parameters[UnisonPhaseSpread],
			DetuneSpread: parameters[UnisonDetuneSpread],
			CurveGamma:   parameters[UnisonCurveGamma],
		})

		// LPF
		lpf := dsp.NewLowPassSVF(SampleRate, unison, parameters[LPFCutoff], parameters[LPFResonance])

		//ctModRateAdsr := dsp.NewADSR(SampleRate, time.Millisecond*5, time.Millisecond*50, 0, time.Millisecond*100)
		//*cutoff.ModInputs() = append(*cutoff.ModInputs(), dsp.NewModInput(ctModRateAdsr, 4000, nil))

		ctModRateOsc := dsp.NewOscillator(SampleRate, dsp.ShapeSine, dsp.NewParam(.5), dsp.NewParam(1), nil)
		//*cutoff.ModInputs() = append(*cutoff.ModInputs(), dsp.NewModInput(ctModRateOsc, 300, nil))

		// Voice
		adsr := dsp.NewADSR(
			SampleRate,
			parameters[AmpEnvAttack],
			parameters[AmpEnvDecay],
			parameters[AmpEnvSustain],
			parameters[AmpEnvRelease],
		)
		voice := dsp.NewVoice(lpf, freq, adsr, ctModRateOsc) //, ctModRateAdsr)

		return voice
	}

	// Polyphonic voice
	poly := dsp.NewPolyVoice(8, voiceFact)

	// Delay
	parameters[FBDelayParam] = dsp.NewParam(0.35)
	parameters[FBFeedBack] = dsp.NewParam(0)
	parameters[FBMix] = dsp.NewParam(0)
	parameters[FBTone] = dsp.NewParam(2000)

	delay := dsp.NewFeedbackDelay(
		SampleRate,
		2.0,
		poly,
		parameters[FBDelayParam], // delay time in seconds
		parameters[FBFeedBack],   // feedback amount (0-1)
		parameters[FBMix],        // mix
		parameters[FBTone],       // tone
	)

	return &Polysynth{
		Node:       delay,
		voice:      poly,
		pitch:      pitchBend,
		pInQueue:   pInQueue,
		pOutQueue:  pOutQueue,
		parameters: parameters,
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

func (p *Polysynth) HandleMessage(m msg.Message) {
	switch m.Kind {
	case ParamPullAllKind:
		p.PublishParameters()
	case ParamUpdateKind:
		if param, ok := p.parameters[m.Key]; ok {
			param.SetBase(m.ValF)
		}
	}
}
