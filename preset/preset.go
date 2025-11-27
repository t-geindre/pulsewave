package preset

import "synth/dsp"

type Preset struct {
	Params   map[uint8]dsp.Param
	Name     string
	ModSlots map[int]*ModSlot
}

func NewPreset() *Preset {
	p := &Preset{Params: make(map[uint8]dsp.Param)}
	p.setDefaults()

	return p
}

func NewPresetFromProto(pb *ProtoPreset) *Preset {
	p := NewPreset()
	p.Name = pb.Name

	for _, e := range pb.Params {
		p.Params[uint8(e.Id)] = dsp.NewParam(e.Value)
	}

	for i := 0; i < ModSlots && i < len(pb.ModSlots); i++ {
		slot := pb.ModSlots[i]
		p.ModSlots[i] = &ModSlot{
			Source:      uint8(slot.Source),
			Destination: uint8(slot.Destination),
			Amount:      slot.Amount,
			Shape:       uint8(slot.Shape),
		}
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

	for i := 0; i < ModSlots; i++ {
		slot := p.ModSlots[i]
		msg.ModSlots = append(msg.ModSlots, &ProtoModSlot{
			Source:      uint32(slot.Source),
			Destination: uint32(slot.Destination),
			Amount:      slot.Amount,
			Shape:       uint32(slot.Shape),
		})
	}

	return msg
}

func (p *Preset) setDefaults() {
	p.Name = "01 Default"

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
	p.Params[VoicesGain] = dsp.NewParam(1.0)

	// Noise oscillator
	p.Params[NoiseGain] = dsp.NewParam(0)
	p.Params[NoiseType] = dsp.NewParam(dsp.NoiseWhite)

	// Sub oscillator
	p.Params[SubOscGain] = dsp.NewParam(0)
	p.Params[SubOscShape] = dsp.NewParam(0)
	p.Params[SubOscTranspose] = dsp.NewParam(0)

	// LFOs
	p.Params[Lfo0rate] = dsp.NewParam(0.33)
	p.Params[Lfo0Phase] = dsp.NewParam(0)
	p.Params[Lfo0Shape] = dsp.NewParam(0)

	p.Params[Lfo1rate] = dsp.NewParam(0.33)
	p.Params[Lfo1Phase] = dsp.NewParam(0)
	p.Params[Lfo1Shape] = dsp.NewParam(0)

	p.Params[Lfo2rate] = dsp.NewParam(0.33)
	p.Params[Lfo2Phase] = dsp.NewParam(0)
	p.Params[Lfo2Shape] = dsp.NewParam(0)

	// ADSRs
	p.Params[Adsr0Attack] = dsp.NewParam(10.0 / 1000)
	p.Params[Adsr0Decay] = dsp.NewParam(10.0 / 1000)
	p.Params[Adsr0Sustain] = dsp.NewParam(.9)
	p.Params[Adsr0Release] = dsp.NewParam(10.0 / 1000)

	p.Params[Adsr1Attack] = dsp.NewParam(10.0 / 1000)
	p.Params[Adsr1Decay] = dsp.NewParam(10.0 / 1000)
	p.Params[Adsr1Sustain] = dsp.NewParam(.9)
	p.Params[Adsr1Release] = dsp.NewParam(10.0 / 1000)

	p.Params[Adsr2Attack] = dsp.NewParam(10.0 / 1000)
	p.Params[Adsr2Decay] = dsp.NewParam(10.0 / 1000)
	p.Params[Adsr2Sustain] = dsp.NewParam(.9)
	p.Params[Adsr2Release] = dsp.NewParam(10.0 / 1000)

	// Reset modulation slots
	p.ModSlots = make(map[int]*ModSlot)
	for i := 0; i < ModSlots; i++ {
		p.ModSlots[i] = &ModSlot{
			Source:      ModSrcLfo0,
			Destination: ParamNone,
			Amount:      0,
			Shape:       ModShapeLinear,
		}
	}
}
