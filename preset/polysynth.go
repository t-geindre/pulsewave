package preset

import (
	"synth/dsp"
	"synth/msg"
	"time"
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

	// Main pitch bend param
	pitchBend := dsp.NewParam(0)

	// Shape registry (uniq for all voices)
	reg := dsp.NewShapeRegistry()
	osc0Shp := reg.Add(dsp.ShapeSaw)
	osc1Shp := reg.Add(dsp.ShapeTriangle)
	osc2Shp := reg.Add(dsp.ShapeTableWave, dsp.NewSineWavetable(1024))

	// Unison parameters (all voices share)
	parameters[UnisonPanSpread] = dsp.NewParam(1.0)
	parameters[UnisonPhaseSpread] = dsp.NewParam(.1)
	parameters[UnisonDetuneSpread] = dsp.NewParam(12.0)
	parameters[UnisonCurveGamma] = dsp.NewParam(1.5)
	parameters[UnisonVoices] = dsp.NewParam(8)

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
				dsp.NewRegOscillator(SampleRate, reg, osc0Shp, ft, ph, nil),
				dsp.NewParam(.33),
				dsp.NewParam(0),
			))

			// 1
			mixer.Add(dsp.NewInput(
				dsp.NewRegOscillator(SampleRate, reg, osc1Shp, dsp.NewTunerParam(ft, dsp.NewParam(-12)), ph, nil),
				dsp.NewParam(.33),
				dsp.NewParam(0),
			))

			// 2
			mixer.Add(dsp.NewInput(
				dsp.NewRegOscillator(SampleRate, reg, osc2Shp, dsp.NewTunerParam(ft, dsp.NewParam(+24)), ph, nil),
				dsp.NewParam(0.15),
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
	parameters[FBDelayParam] = dsp.NewParam(0.35)
	parameters[FBFeedBack] = dsp.NewParam(0.3)
	parameters[FBMix] = dsp.NewParam(0.2)
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
