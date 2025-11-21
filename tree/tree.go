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
				NewSelectorNode("Type", preset.PresetUpdateKind, preset.NoiseType,
					NewSelectorOption("White", "", dsp.NoiseWhite),
					NewSelectorOption("Pink", "", dsp.NoisePink),
					NewSelectorOption("Brown", "", dsp.NoiseBrown),
					NewSelectorOption("Blue", "", dsp.NoiseBlue),
					NewSelectorOption("Gaussian", "", dsp.NoiseGaussian),
				),
				NewSliderNode("Gain", preset.PresetUpdateKind, preset.NoiseGain, 0, 1, .01, nil),
			),
			NewNode("Sub",
				NewWaveFormNode(preset.SubOscShape),
				NewSliderNode("Gain", preset.PresetUpdateKind, preset.SubOscGain, 0, 1, .01, nil),
				NewSliderNode("Transpose", preset.PresetUpdateKind, preset.SubOscTranspose, -48, 48, 12, formatOctave),
			),
		),
		NewNode("Modulation",
			NewAdsrNode("Amplitude", preset.AmpEnvAttack, preset.AmpEnvDecay, preset.AmpEnvSustain, preset.AmpEnvRelease),
			NewNode("Cutoff",
				NewLfoNode("LFO", preset.LpfLfoOnOff, preset.LpfLfoShape, preset.LpfLfoFreq, preset.LpfLfoPhase, preset.LpfLfoAmount),
				NewAdsrNodeWithToggle("ADSR", preset.LpfAdsrOnOff, preset.LpfAdsrAttack, preset.LpfAdsrDecay, preset.LpfAdsrSustain, preset.LpfAdsrRelease,
					NewSliderNode("Amount", preset.PresetUpdateKind, preset.LpfAdsrAmount, -20000, 20000, 1, formatHertz),
				),
			),
			NewNode("Pitch",
				NewLfoNode("LFO", preset.PitchLfoOnOff, preset.PitchLfoShape, preset.PitchLfoFreq, preset.PitchLfoPhase, preset.PitchLfoAmount),
				NewAdsrNodeWithToggle("ADSR", preset.PitchAdsrOnOff, preset.PitchAdsrAttack, preset.PitchAdsrDecay, preset.PitchAdsrSustain, preset.PitchAdsrRelease,
					NewSliderNode("Amount", preset.PresetUpdateKind, preset.PitchAdsrAmount, -1000, 1000, .01, formatSemiTon),
				),
			),
		),
		NewNode("Effects",
			NewNode("Feedback delay",
				NewOnOffNode(preset.FBOnOff),
				NewSliderNode("Delay", preset.PresetUpdateKind, preset.FBDelayParam, 0, 2, .001, formatMillisecond),
				NewSliderNode("Feedback", preset.PresetUpdateKind, preset.FBFeedBack, 0, 0.95, .01, nil),
				NewSliderNode("Mix", preset.PresetUpdateKind, preset.FBMix, 0, 1, .01, nil),
				NewSliderNode("Tone", preset.PresetUpdateKind, preset.FBTone, 200, 8000, 1, formatHertz),
			),
			NewNode("Low pass filter",
				NewOnOffNode(preset.LPFOnOff),
				NewSliderNode("Cutoff", preset.PresetUpdateKind, preset.LPFCutoff, 20, 20000, 1, formatHertz),
				NewSliderNode("Resonance", preset.PresetUpdateKind, preset.LPFResonance, 0.01, 10, .01, nil),
			),
			NewNode("Unison",
				NewOnOffNode(preset.UnisonOnOff),
				NewSliderNode("Voices", preset.PresetUpdateKind, preset.UnisonVoices, 1, 16, 1, formatVoice),
				NewSliderNode("Pan spread", preset.PresetUpdateKind, preset.UnisonPanSpread, 0, 1, .01, nil),
				NewSliderNode("Phase spread", preset.PresetUpdateKind, preset.UnisonPhaseSpread, 0, 1, .01, formatCycle),
				NewSliderNode("Detune spread", preset.PresetUpdateKind, preset.UnisonDetuneSpread, 0, 100, .1, formatCent),
				NewSliderNode("Curve gamma", preset.PresetUpdateKind, preset.UnisonCurveGamma, 0.1, 4, .1, nil),
			),
		),
		NewNode("Voices",
			NewSelectorNode("Steal mode", preset.PresetUpdateKind, preset.VoicesStealMode,
				NewSelectorOption("Oldest", "", dsp.PolyStealOldest),
				NewSelectorOption("Lowest pitch", "", dsp.PolyStealLowest),
				NewSelectorOption("Highest pitch", "", dsp.PolyStealHighest),
			),
			NewSliderNode("Active voices", preset.PresetUpdateKind, preset.VoicesActive, 1, 16, 1, formatVoice),
			NewSliderNode("Pitch glide", preset.PresetUpdateKind, preset.VoicesPitchGlide, 0, 1, .001, formatMillisecond),
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

	// For now, amp is modulated only by and ADSR
	tree.QueryAll("Amplitude")[0].AttachPreview(func() string {
		return "ADSR"
	})

	return tree
}
