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
	p[Osc0Shape] = dsp.NewParam(0)
	p[Osc1Shape] = dsp.NewParam(0)
	p[Osc2Shape] = dsp.NewParam(0)

	p[Osc0Gain] = dsp.NewParam(0.33)
	p[Osc1Gain] = dsp.NewParam(0)
	p[Osc2Gain] = dsp.NewParam(0)

	p[Osc0Detune] = dsp.NewParam(0)
	p[Osc1Detune] = dsp.NewParam(0)
	p[Osc2Detune] = dsp.NewParam(0)

	p[Osc0Phase] = dsp.NewParam(0)
	p[Osc1Phase] = dsp.NewParam(0)
	p[Osc2Phase] = dsp.NewParam(0)

	p[Osc0Pw] = dsp.NewParam(0)
	p[Osc1Pw] = dsp.NewParam(0)
	p[Osc2Pw] = dsp.NewParam(0)

	// Pitch mod
	p[PitchLfoOnOff] = dsp.NewParam(0)
	p[PitchLfoAmount] = dsp.NewParam(0)
	p[PitchLfoShape] = dsp.NewParam(0)
	p[PitchLfoFreq] = dsp.NewParam(0)
	p[PitchLfoPhase] = dsp.NewParam(0)

	p[PitchAdsrOnOff] = dsp.NewParam(0)
	p[PitchAdsrAmount] = dsp.NewParam(0)
	p[PitchAdsrAttack] = dsp.NewParam(0)
	p[PitchAdsrDecay] = dsp.NewParam(0)
	p[PitchAdsrSustain] = dsp.NewParam(0)
	p[PitchAdsrRelease] = dsp.NewParam(0)

	// Amplitude envelope
	p[AmpEnvAttack] = dsp.NewParam(.01)
	p[AmpEnvDecay] = dsp.NewParam(.01)
	p[AmpEnvSustain] = dsp.NewParam(.9)
	p[AmpEnvRelease] = dsp.NewParam(.01)

	// Unison (all voices share)
	p[UnisonOnOff] = dsp.NewParam(0)
	p[UnisonPanSpread] = dsp.NewParam(0)
	p[UnisonPhaseSpread] = dsp.NewParam(0)
	p[UnisonDetuneSpread] = dsp.NewParam(0)
	p[UnisonCurveGamma] = dsp.NewParam(0)
	p[UnisonVoices] = dsp.NewParam(0)

	// Feedback Delay
	p[FBDelayParam] = dsp.NewParam(0)
	p[FBFeedBack] = dsp.NewParam(0)
	p[FBMix] = dsp.NewParam(0)
	p[FBTone] = dsp.NewParam(0)
	p[FBOnOff] = dsp.NewParam(0)

	// Low pass filter
	p[LPFOnOff] = dsp.NewParam(0)
	p[LPFCutoff] = dsp.NewParam(0)
	p[LPFResonance] = dsp.NewParam(0)

	p[LpfLfoOnOff] = dsp.NewParam(0)
	p[LpfLfoAmount] = dsp.NewParam(0)
	p[LpfLfoShape] = dsp.NewParam(0)
	p[LpfLfoFreq] = dsp.NewParam(0)
	p[LpfLfoPhase] = dsp.NewParam(0)

	p[LpfAdsrOnOff] = dsp.NewParam(0)
	p[LpfAdsrAmount] = dsp.NewParam(0)
	p[LpfAdsrAttack] = dsp.NewParam(0)
	p[LpfAdsrDecay] = dsp.NewParam(0)
	p[LpfAdsrSustain] = dsp.NewParam(0)
	p[LpfAdsrRelease] = dsp.NewParam(0)
}
