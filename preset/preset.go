package preset

import "synth/dsp"

type Preset map[uint8]dsp.Param

func NewPreset() Preset {
	p := make(Preset)
	p.setDefaults()

	return p
}

func NewFromProto(pb *ProtoPreset) Preset {
	p := make(Preset)
	p.setDefaults()
	for _, e := range pb.Params {
		p[uint8(e.Id)] = dsp.NewParam(e.Value)
	}
	return p
}

func (p Preset) ToProto() *ProtoPreset {
	msg := &ProtoPreset{}
	for id, param := range p {
		msg.Params = append(msg.Params, &ProtoParamEntry{
			Id:    uint32(id),
			Value: param.GetBase(),
		})
	}
	return msg
}

func (p Preset) setDefaults() {
	for _, id := range []uint8{
		// Oscillator todo phase + pulse width
		Osc0Shape,
		Osc1Shape,
		Osc2Shape,

		Osc0Gain,
		Osc1Gain,
		Osc2Gain,

		Osc0Detune,
		Osc1Detune,
		Osc2Detune,

		Osc0Phase,
		Osc1Phase,
		Osc2Phase,

		Osc0Pw,
		Osc1Pw,
		Osc2Pw,

		// Pitch mod
		PitchLfoOnOff,
		PitchLfoAmount,
		PitchLfoShape,
		PitchLfoFreq,
		PitchLfoPhase,

		PitchAdsrOnOff,
		PitchAdsrAmount,
		PitchAdsrAttack,
		PitchAdsrDecay,
		PitchAdsrSustain,
		PitchAdsrRelease,

		// Amplitude envelope
		AmpEnvAttack,
		AmpEnvDecay,
		AmpEnvSustain,
		AmpEnvRelease,

		// Unison (all voices share)
		UnisonOnOff,
		UnisonPanSpread,
		UnisonPhaseSpread,
		UnisonDetuneSpread,
		UnisonCurveGamma,
		UnisonVoices,

		// Feedback Delay
		FBDelayParam,
		FBFeedBack,
		FBMix,
		FBTone,
		FBOnOff,

		// Low pass filter
		LPFOnOff,
		LPFCutoff,
		LPFResonance,

		LpfLfoOnOff,
		LpfLfoAmount,
		LpfLfoShape,
		LpfLfoFreq,
		LpfLfoPhase,

		LpfAdsrOnOff,
		LpfAdsrAmount,
		LpfAdsrAttack,
		LpfAdsrDecay,
		LpfAdsrSustain,
		LpfAdsrRelease,
	} {
		p[id] = dsp.NewParam(0)
	}

	// Output some sounds by default
	p[Osc0Gain].SetBase(.33)
	p[AmpEnvAttack].SetBase(0.005)
	p[AmpEnvDecay].SetBase(0.005)
	p[AmpEnvSustain].SetBase(0.9)
	p[AmpEnvRelease].SetBase(0.005)

	// Set default oscillator PW
	p[Osc0Pw].SetBase(0.5)
	p[Osc1Pw].SetBase(0.5)
	p[Osc2Pw].SetBase(0.5)

	p[UnisonVoices].SetBase(1)
}
