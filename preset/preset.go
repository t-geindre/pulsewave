package preset

import "synth/dsp"

type Preset map[uint8]dsp.Param

func NewPreset() Preset {
	p := make(Preset)

	// Oscillator
	p[Osc0Shape] = dsp.NewParam(2) // Saw
	p[Osc1Shape] = dsp.NewParam(0) // Sine
	p[Osc2Shape] = dsp.NewParam(0) // Sine

	p[Osc0Gain] = dsp.NewParam(.33)
	p[Osc1Gain] = dsp.NewParam(.10)
	p[Osc2Gain] = dsp.NewParam(.10)

	p[Osc0Detune] = dsp.NewParam(0)
	p[Osc1Detune] = dsp.NewParam(12)
	p[Osc2Detune] = dsp.NewParam(-12)

	// Amplitude envelope
	p[AmpEnvAttack] = dsp.NewParam(0.02)
	p[AmpEnvDecay] = dsp.NewParam(0.02)
	p[AmpEnvSustain] = dsp.NewParam(0.9)
	p[AmpEnvRelease] = dsp.NewParam(0.02)

	// Unison (all voices share)
	p[UnisonOnOff] = dsp.NewParam(1)
	p[UnisonPanSpread] = dsp.NewParam(1)
	p[UnisonPhaseSpread] = dsp.NewParam(.1)
	p[UnisonDetuneSpread] = dsp.NewParam(10)
	p[UnisonCurveGamma] = dsp.NewParam(1.5)
	p[UnisonVoices] = dsp.NewParam(8)

	// Feedback Delay
	p[FBDelayParam] = dsp.NewParam(0.35)
	p[FBFeedBack] = dsp.NewParam(.3)
	p[FBMix] = dsp.NewParam(.2)
	p[FBTone] = dsp.NewParam(2000)
	p[FBOnOff] = dsp.NewParam(1)

	// Low pass filter
	p[LPFOnOff] = dsp.NewParam(1)
	p[LPFCutoff] = dsp.NewParam(2000)
	p[LPFResonance] = dsp.NewParam(1)

	p[LpfLfoOnOff] = dsp.NewParam(0)
	p[LpfLfoAmount] = dsp.NewParam(1500)
	p[LpfLfoShape] = dsp.NewParam(0) // Sine
	p[LpfLfoFreq] = dsp.NewParam(0.5)
	p[LpfLfoPhase] = dsp.NewParam(0)

	p[LpfAdsrOnOff] = dsp.NewParam(0)
	p[LpfAdsrAmount] = dsp.NewParam(4000)
	p[LpfAdsrAttack] = dsp.NewParam(0.01)
	p[LpfAdsrDecay] = dsp.NewParam(0.1)
	p[LpfAdsrSustain] = dsp.NewParam(1)
	p[LpfAdsrRelease] = dsp.NewParam(0.2)

	return p
}
