package preset

import "synth/dsp"

type Preset struct {
	Params map[uint8]dsp.Param
	Name   string
}

func NewPreset() *Preset {
	p := &Preset{
		Params: make(map[uint8]dsp.Param),
		Name:   "01 Default",
	}
	p.setDefaults()

	return p
}

func NewPresetFromProto(pb *ProtoPreset) *Preset {
	p := NewPreset()
	p.Name = pb.Name
	for _, e := range pb.Params {
		p.Params[uint8(e.Id)] = dsp.NewParam(e.Value)
	}
	return p
}

func (p *Preset) ToProto() *ProtoPreset {
	msg := &ProtoPreset{}
	msg.Name = p.Name
	for id, param := range p.Params {
		msg.Params = append(msg.Params, &ProtoParamEntry{
			Id:    uint32(id),
			Value: param.GetBase(),
		})
	}
	return msg
}

func (p *Preset) setDefaults() {
	p.Params[Osc0Shape] = dsp.NewParam(0)
	p.Params[Osc1Shape] = dsp.NewParam(0)
	p.Params[Osc2Shape] = dsp.NewParam(0)

	p.Params[Osc0Gain] = dsp.NewParam(0.33)
	p.Params[Osc1Gain] = dsp.NewParam(0)
	p.Params[Osc2Gain] = dsp.NewParam(0)

	p.Params[Osc0Detune] = dsp.NewParam(0)
	p.Params[Osc1Detune] = dsp.NewParam(0)
	p.Params[Osc2Detune] = dsp.NewParam(0)

	p.Params[Osc0Phase] = dsp.NewParam(0)
	p.Params[Osc1Phase] = dsp.NewParam(0)
	p.Params[Osc2Phase] = dsp.NewParam(0)

	p.Params[Osc0Pw] = dsp.NewParam(0.5)
	p.Params[Osc1Pw] = dsp.NewParam(0.5)
	p.Params[Osc2Pw] = dsp.NewParam(0.5)

	// Pitch mod
	p.Params[PitchLfoOnOff] = dsp.NewParam(0)
	p.Params[PitchLfoAmount] = dsp.NewParam(0)
	p.Params[PitchLfoShape] = dsp.NewParam(0)
	p.Params[PitchLfoFreq] = dsp.NewParam(0)
	p.Params[PitchLfoPhase] = dsp.NewParam(0)

	p.Params[PitchAdsrOnOff] = dsp.NewParam(0)
	p.Params[PitchAdsrAmount] = dsp.NewParam(0)
	p.Params[PitchAdsrAttack] = dsp.NewParam(0)
	p.Params[PitchAdsrDecay] = dsp.NewParam(0)
	p.Params[PitchAdsrSustain] = dsp.NewParam(0)
	p.Params[PitchAdsrRelease] = dsp.NewParam(0)

	// Amplitude envelope
	p.Params[AmpEnvAttack] = dsp.NewParam(10.0 / 1000)
	p.Params[AmpEnvDecay] = dsp.NewParam(10.0 / 1000)
	p.Params[AmpEnvSustain] = dsp.NewParam(.9)
	p.Params[AmpEnvRelease] = dsp.NewParam(10.0 / 1000)

	// Unison (all voices share)
	p.Params[UnisonOnOff] = dsp.NewParam(0)
	p.Params[UnisonPanSpread] = dsp.NewParam(1)
	p.Params[UnisonPhaseSpread] = dsp.NewParam(.1)
	p.Params[UnisonDetuneSpread] = dsp.NewParam(12)
	p.Params[UnisonCurveGamma] = dsp.NewParam(1.5)
	p.Params[UnisonVoices] = dsp.NewParam(4)

	// Feedback Delay
	p.Params[FBOnOff] = dsp.NewParam(0)
	p.Params[FBDelayParam] = dsp.NewParam(350.0 / 1000)
	p.Params[FBFeedBack] = dsp.NewParam(.3)
	p.Params[FBMix] = dsp.NewParam(.3)
	p.Params[FBTone] = dsp.NewParam(5000)

	// Low pass filter
	p.Params[LPFOnOff] = dsp.NewParam(0)
	p.Params[LPFCutoff] = dsp.NewParam(3000)
	p.Params[LPFResonance] = dsp.NewParam(1)

	p.Params[LpfLfoOnOff] = dsp.NewParam(0)
	p.Params[LpfLfoAmount] = dsp.NewParam(800)
	p.Params[LpfLfoShape] = dsp.NewParam(0)
	p.Params[LpfLfoFreq] = dsp.NewParam(.3)
	p.Params[LpfLfoPhase] = dsp.NewParam(0)

	p.Params[LpfAdsrOnOff] = dsp.NewParam(0)
	p.Params[LpfAdsrAmount] = dsp.NewParam(0)
	p.Params[LpfAdsrAttack] = dsp.NewParam(0)
	p.Params[LpfAdsrDecay] = dsp.NewParam(0)
	p.Params[LpfAdsrSustain] = dsp.NewParam(0)
	p.Params[LpfAdsrRelease] = dsp.NewParam(0)

	// Voices
	p.Params[VoicesStealMode] = dsp.NewParam(dsp.PolyStealOldest)
	p.Params[VoicesActive] = dsp.NewParam(8)
	p.Params[VoicesPitchGlide] = dsp.NewParam(0.000)

	// Noise oscillator
	p.Params[NoiseGain] = dsp.NewParam(0)

	// Sub oscillator
	p.Params[SubOscGain] = dsp.NewParam(0)
	p.Params[SubOscShape] = dsp.NewParam(0)
	p.Params[SubOscTranspose] = dsp.NewParam(0)
}
