package tree

import (
	"synth/dsp"
	"synth/preset"
	"synth/settings"
)

func NewTree(presets []string) Node {
	tree := NewNode("",
		NewNode("Oscillators",
			NewOscillatorNode("Osc 01", preset.Osc0Shape, preset.Osc0Detune, preset.Osc0Gain, preset.Osc0Phase, preset.Osc0Pw),
			NewOscillatorNode("Osc 02", preset.Osc1Shape, preset.Osc1Detune, preset.Osc1Gain, preset.Osc1Phase, preset.Osc1Pw),
			NewOscillatorNode("Osc 03", preset.Osc2Shape, preset.Osc2Detune, preset.Osc2Gain, preset.Osc2Phase, preset.Osc2Pw),
			NewNode("Noise",
				NewSelectorNode("Type", preset.UpdateParameterKind, preset.NoiseType,
					NewSelectorOption("White", "", dsp.NoiseWhite),
					NewSelectorOption("Pink", "", dsp.NoisePink),
					NewSelectorOption("Brown", "", dsp.NoiseBrown),
					NewSelectorOption("Blue", "", dsp.NoiseBlue),
					NewSelectorOption("Gaussian", "", dsp.NoiseGaussian),
				),
				NewSliderNode("Gain", preset.UpdateParameterKind, preset.NoiseGain, 0, 1, .01, nil),
			),
			NewNode("Sub",
				NewWaveFormNode(preset.SubOscShape),
				NewSliderNode("Gain", preset.UpdateParameterKind, preset.SubOscGain, 0, 1, .01, nil),
				NewSliderNode("Transpose", preset.UpdateParameterKind, preset.SubOscTranspose, -48, 48, 12, formatOctave),
			),
		),
		NewNode("Modulation",
			NewModulationMatrixNode("Matrix"),
			// NewNode("Velocity"), // todo implement velocity modulation source
			NewLfoNode("LFO 01", preset.Lfo0Shape, preset.Lfo0rate, preset.Lfo0Phase),
			NewLfoNode("LFO 02", preset.Lfo1Shape, preset.Lfo1rate, preset.Lfo1Phase),
			NewLfoNode("LFO 03", preset.Lfo2Shape, preset.Lfo2rate, preset.Lfo2Phase),
			NewAdsrNode("ADSR 01", preset.Adsr0Attack, preset.Adsr0Decay, preset.Adsr0Sustain, preset.Adsr0Release),
			NewAdsrNode("ADSR 02", preset.Adsr1Attack, preset.Adsr1Decay, preset.Adsr1Sustain, preset.Adsr1Release),
			NewAdsrNode("ADSR 03", preset.Adsr2Attack, preset.Adsr2Decay, preset.Adsr2Sustain, preset.Adsr2Release),
		),
		NewNode("Effects",
			NewNode("Feedback delay",
				NewOnOffNode(preset.FBOnOff),
				NewSliderNode("Delay", preset.UpdateParameterKind, preset.FBDelayParam, 0, 2, .001, formatMillisecond),
				NewSliderNode("Feedback", preset.UpdateParameterKind, preset.FBFeedBack, 0, 0.95, .01, nil),
				NewSliderNode("Mix", preset.UpdateParameterKind, preset.FBMix, 0, 1, .01, nil),
				NewSliderNode("Tone", preset.UpdateParameterKind, preset.FBTone, 200, 8000, 1, formatHertz),
			),
			NewNode("Low pass filter",
				NewOnOffNode(preset.LPFOnOff),
				NewSliderNode("Cutoff", preset.UpdateParameterKind, preset.LPFCutoff, 20, 20000, 1, formatHertz),
				NewSliderNode("Resonance", preset.UpdateParameterKind, preset.LPFResonance, 0.01, 10, .01, nil),
			),
			NewNode("Unison",
				NewOnOffNode(preset.UnisonOnOff),
				NewSliderNode("Voices", preset.UpdateParameterKind, preset.UnisonVoices, 1, 16, 1, formatVoice),
				NewSliderNode("Pan spread", preset.UpdateParameterKind, preset.UnisonPanSpread, 0, 1, .01, nil),
				NewSliderNode("Phase spread", preset.UpdateParameterKind, preset.UnisonPhaseSpread, 0, 1, .01, formatCycle),
				NewSliderNode("Detune spread", preset.UpdateParameterKind, preset.UnisonDetuneSpread, 0, 100, .1, formatCent),
				NewSliderNode("Curve gamma", preset.UpdateParameterKind, preset.UnisonCurveGamma, 0.1, 4, .1, nil),
			),
		),
		NewNode("Voices",
			NewSelectorNode("Steal mode", preset.UpdateParameterKind, preset.VoicesStealMode,
				NewSelectorOption("Oldest", "", dsp.PolyStealOldest),
				NewSelectorOption("Lowest pitch", "", dsp.PolyStealLowest),
				NewSelectorOption("Highest pitch", "", dsp.PolyStealHighest),
			),
			NewSliderNode("Active voices", preset.UpdateParameterKind, preset.VoicesActive, 1, preset.MaxVoices, 1, formatVoice),
			NewSliderNode("Pitch glide", preset.UpdateParameterKind, preset.VoicesPitchGlide, 0, 1, .001, formatMillisecond),
			NewSliderNode("Gain", preset.UpdateParameterKind, preset.VoicesGain, 0, 1, .01, nil),
			NewSliderNode("Pitch", preset.UpdateParameterKind, preset.VoicesPitch, -48, 48, .01, formatSemiTon),
		),
		NewNode("Visualizer",
			NewFeatureNode("Spectrum", 0), // todo implement spectrum analyzer
			NewFeatureNode("Oscilloscope", FeatureOscilloscope),
		),
		NewNode("Presets",
			NewPresetsNodes(presets)...,
		),
		NewNode("Settings",
			NewSliderNode("Master gain", settings.SettingUpdateKind, settings.MasterGain, 0, 3, .01, nil),
			NewSliderNode("Pitch bend range", settings.SettingUpdateKind, settings.PitchBendRange, 1, 24, 1, formatSemiTon),
		),
	)

	AttachOscPreviews(tree, "Osc 01", "Osc 02", "Osc 03", "Noise", "Sub")
	AttachPreviewToParent(tree, "Status")
	AttachNameIfSubNodeVal(tree, "Status", 1, " + ", "OFF", "Cutoff", "Pitch")

	return tree
}
