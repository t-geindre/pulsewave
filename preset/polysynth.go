package preset

import (
	"fmt"
	"synth/dsp"
	"synth/msg"
)

type Polysynth struct {
	dsp.Node
	voice     *dsp.PolyVoice
	pitch     dsp.Param
	messenger *msg.Messenger
	modSlots  map[int]*ModSlot

	modulators map[uint8]dsp.ParamModulator
	parameters map[uint8]dsp.Param

	voiceModulators []map[uint8]dsp.ParamModulator // per voice
	voiceParams     []map[uint8]dsp.Param          // per voice
}

const MaxVoices = 16

func NewPolysynth(SampleRate float64) *Polysynth {
	// Parameters map
	preset := NewPreset()

	// Shape registry (uniq for all oscillators)
	reg := dsp.NewShapeRegistry()
	reg.Add(dsp.ShapeTableWave, dsp.NewSineWavetable(1024))
	reg.Add(dsp.ShapeSquare)
	reg.Add(dsp.ShapeSaw)
	reg.Add(dsp.ShapeTriangle)

	// Global pitch bend
	pitchBend := dsp.NewSmoothedParam(SampleRate, 0, dsp.NewConstParam(.01))

	// Per voice parameters
	voiceModulators := make([]map[uint8]dsp.ParamModulator, 0)
	voiceParams := make([]map[uint8]dsp.Param, 0)

	// Voice factory / 3 osc
	voiceFact := func() *dsp.Voice {
		// Voice modulators
		modulators := make(map[uint8]dsp.ParamModulator)
		modulators[ModSrcLfo0] = dsp.NewRegOscillator(SampleRate, reg, preset.Params[Lfo0Shape], preset.Params[Lfo0rate], preset.Params[Lfo0Phase], nil)
		modulators[ModSrcLfo1] = dsp.NewRegOscillator(SampleRate, reg, preset.Params[Lfo1Shape], preset.Params[Lfo1rate], preset.Params[Lfo1Phase], nil)
		modulators[ModSrcLfo2] = dsp.NewRegOscillator(SampleRate, reg, preset.Params[Lfo2Shape], preset.Params[Lfo2rate], preset.Params[Lfo2Phase], nil)
		modulators[ModSrcAdsr0] = dsp.NewADSR(SampleRate, preset.Params[Adsr0Attack], preset.Params[Adsr0Decay], preset.Params[Adsr0Sustain], preset.Params[Adsr0Release])
		modulators[ModSrcAdsr1] = dsp.NewADSR(SampleRate, preset.Params[Adsr1Attack], preset.Params[Adsr1Decay], preset.Params[Adsr1Sustain], preset.Params[Adsr1Release])
		modulators[ModSrcAdsr2] = dsp.NewADSR(SampleRate, preset.Params[Adsr2Attack], preset.Params[Adsr2Decay], preset.Params[Adsr2Sustain], preset.Params[Adsr2Release])
		voiceModulators = append(voiceModulators, modulators)

		// Voice params
		params := createLocalParametersMap(preset.Params,
			// Oscillators
			Osc0Phase, Osc1Phase, Osc2Phase,
			Osc0Detune, Osc1Detune, Osc2Detune,
			Osc0Pw, Osc1Pw, Osc2Pw,
			Osc0Gain, Osc1Gain, Osc2Gain,
			// Voices
			VoicesPitch,
			// LPF
			LPFCutoff, LPFResonance,
		)
		voiceParams = append(voiceParams, params)

		// Base frequency param (uniq per voice)
		freq := dsp.NewSmoothedParam(SampleRate, 440, preset.Params[VoicesPitchGlide])
		pitchMod := params[VoicesPitch]
		pitch := dsp.NewTunerParam(dsp.NewTunerParam(freq, pitchBend), pitchMod)

		// Oscillator factory
		oscFact := func(ph, dt dsp.Param) dsp.Node {
			// Mixer, registry
			mixer := dsp.NewMixer(nil, false)
			ft := dsp.NewTunerParam(pitch, dt)

			// 0
			ph0 := dsp.NewParam(0)
			*ph0.ModInputs() = append(*ph0.ModInputs(),
				dsp.NewModInput(ph, dsp.NewConstParam(1), nil),
				dsp.NewModInput(params[Osc0Phase], dsp.NewConstParam(1), nil),
			)
			mixer.Add(dsp.NewInput(
				dsp.NewRegOscillator(SampleRate, reg, preset.Params[Osc0Shape], dsp.NewTunerParam(ft, params[Osc0Detune]), ph0, params[Osc0Pw]),
				params[Osc0Gain],
				dsp.NewParam(0),
			))

			// 1
			ph1 := dsp.NewParam(0)
			*ph1.ModInputs() = append(*ph1.ModInputs(),
				dsp.NewModInput(ph, dsp.NewConstParam(1), nil),
				dsp.NewModInput(params[Osc1Phase], dsp.NewConstParam(1), nil),
			)
			mixer.Add(dsp.NewInput(
				dsp.NewRegOscillator(SampleRate, reg, preset.Params[Osc1Shape], dsp.NewTunerParam(ft, params[Osc1Detune]), ph1, params[Osc1Pw]),
				params[Osc1Gain],
				dsp.NewParam(0),
			))

			// 2
			ph2 := dsp.NewParam(0)
			*ph2.ModInputs() = append(*ph2.ModInputs(),
				dsp.NewModInput(ph, dsp.NewConstParam(1), nil),
				dsp.NewModInput(params[Osc2Phase], dsp.NewConstParam(1), nil),
			)
			mixer.Add(dsp.NewInput(
				dsp.NewRegOscillator(SampleRate, reg, preset.Params[Osc2Shape], dsp.NewTunerParam(ft, params[Osc2Detune]), ph2, params[Osc2Pw]),
				params[Osc2Gain],
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
		lpf := dsp.NewLowPassSVF(SampleRate, globalMix, params[LPFCutoff], params[LPFResonance])
		lpfSkip := NewNodeSkipper(lpf, globalMix, preset.Params[LPFOnOff])

		// Amplitude envelope
		gain := dsp.NewParam(0)
		*gain.ModInputs() = append(*gain.ModInputs(),
			dsp.NewModInput(modulators[ModSrcAdsr0], preset.Params[VoicesGain], nil),
		)

		vca := dsp.NewVca(lpfSkip, gain)

		// Voice
		voice := dsp.NewVoice(vca, freq,
			modulators[ModSrcAdsr0], // First one drives the voice
			modulators[ModSrcAdsr1],
			modulators[ModSrcAdsr2],
			modulators[ModSrcLfo0],
			modulators[ModSrcLfo1],
			modulators[ModSrcLfo2],
		)

		return voice
	}

	// Polyphonic voice
	poly := dsp.NewPolyVoice(MaxVoices, preset.Params[VoicesActive], preset.Params[VoicesStealMode], voiceFact)

	// Delay with skipper
	delay := dsp.NewFeedbackDelay(SampleRate, 2.0, poly, preset.Params[FBDelayParam], preset.Params[FBFeedBack], preset.Params[FBMix], preset.Params[FBTone])
	delaySkip := NewNodeSkipper(delay, poly, preset.Params[FBOnOff])

	// Global modulators
	modulators := make(map[uint8]dsp.ParamModulator)
	modulators[ModSrcLfo0] = dsp.NewRegOscillator(SampleRate, reg, preset.Params[Lfo0Shape], preset.Params[Lfo0rate], preset.Params[Lfo0Phase], nil)
	modulators[ModSrcLfo1] = dsp.NewRegOscillator(SampleRate, reg, preset.Params[Lfo1Shape], preset.Params[Lfo1rate], preset.Params[Lfo1Phase], nil)
	modulators[ModSrcLfo2] = dsp.NewRegOscillator(SampleRate, reg, preset.Params[Lfo2Shape], preset.Params[Lfo2rate], preset.Params[Lfo2Phase], nil)

	// Modulation slots
	modSlots := make(map[int]*ModSlot)
	for i, slot := range preset.ModSlots {
		modSlots[i] = &ModSlot{
			Source:         slot.Source,
			Destination:    slot.Destination,
			Amount:         slot.Amount,
			Shape:          slot.Shape,
			GlobalModInput: dsp.NewModInput(modulators[slot.Source], dsp.NewParam(slot.Amount), nil),
		}
		for j := 0; j < MaxVoices; j++ {
			modSlots[i].PerVoiceModInput = append(modSlots[i].PerVoiceModInput,
				dsp.NewModInput(voiceModulators[j][slot.Source], dsp.NewParam(slot.Amount), nil),
			)
		}
	}

	return &Polysynth{
		Node:            delaySkip,
		voice:           poly,
		pitch:           pitchBend,
		modSlots:        modSlots,
		modulators:      modulators,
		parameters:      preset.Params,
		voiceModulators: voiceModulators,
		voiceParams:     voiceParams,
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

func (p *Polysynth) UpdateModMatrix(slot int, key uint8, val float32) {
	switch key {
	case ModParamSrc:
		p.UpdateModSource(slot, uint8(val))
	case ModParamDst:
		p.UpdateModDestination(slot, uint8(val))
	case ModParamAmt:
		p.UpdateModAmount(slot, val)
	case ModParamShp:
		fmt.Printf("Update mod slot %d shape to %d\n", slot, int(val))
	default:
		panic("unknown param")
	}
}

func (p *Polysynth) UpdateModSource(s int, src uint8) {
	slot := p.modSlots[s]

	if slot.Source == src {
		return
	}
	slot.Source = src

	slot.GlobalModInput.SetSrc(p.modulators[slot.Source])
	for v, mi := range slot.PerVoiceModInput {
		mi.SetSrc(p.voiceModulators[v][slot.Source])
	}
}

func (p *Polysynth) UpdateModDestination(s int, dst uint8) {
	slot := p.modSlots[s]
	if slot.Destination == dst {
		return
	}

	// Detach previous destination
	if slot.Destination != ParamNone {
		// Per voice
		for v, mi := range slot.PerVoiceModInput {
			param := p.voiceParams[v][slot.Destination]
			if param != nil {
				param.RemoveModInput(mi)
			}
		}
		// Global
		param := p.parameters[slot.Destination]
		if param != nil {
			param.RemoveModInput(slot.GlobalModInput)
		}
	}

	// Attach new destination
	if dst != ParamNone {
		// Per voice
		found := false
		for v, mi := range slot.PerVoiceModInput {
			param := p.voiceParams[v][dst]
			if param != nil {
				found = true
				param.AddModInput(mi)
			}
		}
		// Global only if not found per voice
		if !found {
			param := p.parameters[dst]
			if param != nil {
				param.AddModInput(slot.GlobalModInput)
			}
		}
	}

	slot.Destination = dst
}

func (p *Polysynth) UpdateModAmount(s int, amt float32) {
	p.modSlots[s].Amount = amt
	p.modSlots[s].GlobalModInput.Amount().SetBase(amt)
	for _, mi := range p.modSlots[s].PerVoiceModInput {
		mi.Amount().SetBase(amt)
	}
}

func (p *Polysynth) LoadPreset(preset *Preset) {
	for key, param := range preset.Params {
		p.SetParam(key, param.GetBase())
	}

	for i, slot := range preset.ModSlots {
		p.UpdateModSource(i, slot.Source)
		p.UpdateModDestination(i, slot.Destination)
		p.UpdateModAmount(i, slot.Amount)
	}
}

func (p *Polysynth) HydratePreset(preset *Preset) *Preset {
	for key, param := range p.parameters {
		preset.Params[key] = dsp.NewParam(param.GetBase())
	}

	for i, slot := range p.modSlots {
		preset.ModSlots[i] = &ModSlot{
			Source:      slot.Source,
			Destination: slot.Destination,
			Amount:      slot.Amount,
			Shape:       slot.Shape,
		}
	}

	return preset
}

func (p *Polysynth) AllNotesOff() {
	p.voice.AllNotesOff()
}

func createLocalParametersMap(src map[uint8]dsp.Param, keys ...uint8) map[uint8]dsp.Param {
	local := make(map[uint8]dsp.Param)

	for _, k := range keys {
		local[k] = dsp.NewParam(0)
		*local[k].ModInputs() = append(*local[k].ModInputs(),
			dsp.NewModInput(src[k], dsp.NewConstParam(1), nil),
		)
	}

	return local
}
