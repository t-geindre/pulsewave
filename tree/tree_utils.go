package tree

import (
	"synth/preset"
)

func NewOscillatorNode(label string, shape, detune, gain, phase, pw uint8) Node {
	return NewNode(label,
		NewWaveFormNode(shape),
		NewSliderNode("Detune", preset.UpdateParameterKind, detune, -100, 100, .01, formatSemiTon),
		NewSliderNode("Gain", preset.UpdateParameterKind, gain, 0, 1, .01, nil),
		NewSliderNode("Phase", preset.UpdateParameterKind, phase, 0, 1, .01, formatCycle),
		NewSliderNode("Pulse width", preset.UpdateParameterKind, pw, 0.01, 0.5, .01, nil),
	)
}

func NewWaveFormNode(key uint8) Node {
	return NewSelectorNode("Waveform", preset.UpdateParameterKind, key,
		NewSelectorOption("Sine", "ui/icons/sine_wave", 0),
		NewSelectorOption("Square", "ui/icons/square_wave", 1),
		NewSelectorOption("Sawtooth", "ui/icons/saw_wave", 2),
		NewSelectorOption("Triangle", "ui/icons/triangle_wave", 3),
	)
}

func NewAdsrNode(label string, att, dec, sus, rel uint8, children ...Node) Node {
	n := NewNode(label,
		NewSliderNode("Attack", preset.UpdateParameterKind, att, 0, 10, .001, formatMillisecond),
		NewSliderNode("Decay", preset.UpdateParameterKind, dec, 0, 10, .001, formatMillisecond),
		NewSliderNode("Sustain", preset.UpdateParameterKind, sus, 0, 1, .01, nil),
		NewSliderNode("Release", preset.UpdateParameterKind, rel, 0, 10, .001, formatMillisecond),
	)

	for _, c := range children {
		n.Append(c)
	}

	return n
}

func NewOnOffNode(key uint8) Node {
	return NewSelectorNode("Status", preset.UpdateParameterKind, key,
		NewSelectorOption("OFF", "", 0),
		NewSelectorOption("ON", "", 1),
	)
}

func NewPresetsNodes(presets []string) []Node {
	nodes := make([]Node, len(presets))
	for i, p := range presets {
		nodes[i] = NewValidatingSelectorNode(p, preset.LoadSavePresetKind, uint8(i),
			NewSelectorOption("Load", "", 0),
			NewSelectorOption("Save", "", 1),
		)
	}

	return nodes
}

func NewLfoNode(label string, shape, rate, phase uint8) Node {
	return NewNode(label,
		NewWaveFormNode(shape),
		NewSliderNode("Rate", preset.UpdateParameterKind, rate, 0.01, 20, .01, formatLowHertz),
		NewSliderNode("Phase", preset.UpdateParameterKind, phase, 0, 1, .01, formatCycle),
	)
}

func NewModulationMatrixNode(label string) Node {
	matrix := NewNode(label)

	for i := uint8(0); i < preset.ModSlots; i++ {
		slotNode := NewNode("NONE",
			NewSelectorNode("Source", preset.ModulationUpdateKind, preset.ModKeysSpacing*i+preset.ModParamSrc,
				//NewSelectorOption("Velocity", "", preset.ModSrcVelocity), // TODO
				NewSelectorOption("LFO 1", "", preset.ModSrcLfo0),
				NewSelectorOption("LFO 2", "", preset.ModSrcLfo1),
				NewSelectorOption("LFO 3", "", preset.ModSrcLfo2),
				NewSelectorOption("ADSR 1", "", preset.ModSrcAdsr0),
				NewSelectorOption("ADSR 2", "", preset.ModSrcAdsr1),
				NewSelectorOption("ADSR 3", "", preset.ModSrcAdsr2),
			),
			NewSelectorNode("Destination", preset.ModulationUpdateKind, preset.ModKeysSpacing*i+preset.ModParamDst, // TODO remove already assigned destinations with the same source
				NewSelectorOption("NONE", "", preset.ParamNone),
				NewSelectorOption("Osc 1 > Gain", "", preset.Osc0Gain),
				NewSelectorOption("Osc 1 > Phase", "", preset.Osc0Phase),
				NewSelectorOption("Osc 1 > Detune", "", preset.Osc0Detune),
				NewSelectorOption("Osc 1 > Pw", "", preset.Osc0Pw),
				NewSelectorOption("Osc 1 > Gain", "", preset.Osc0Gain),
				NewSelectorOption("Osc 2 > Gain", "", preset.Osc1Gain),
				NewSelectorOption("Osc 2 > Phase", "", preset.Osc1Phase),
				NewSelectorOption("Osc 2 > Detune", "", preset.Osc1Detune),
				NewSelectorOption("Osc 2 > Pw", "", preset.Osc1Pw),
				NewSelectorOption("Osc 2 > Gain", "", preset.Osc1Gain),
				NewSelectorOption("Osc 3 > Gain", "", preset.Osc2Gain),
				NewSelectorOption("Osc 3 > Phase", "", preset.Osc2Phase),
				NewSelectorOption("Osc 3 > Detune", "", preset.Osc2Detune),
				NewSelectorOption("Osc 3 > Pw", "", preset.Osc2Pw),
				NewSelectorOption("Osc 3 > Gain", "", preset.Osc2Gain),
				NewSelectorOption("Voices > Pitch", "", preset.VoicesPitch),
				NewSelectorOption("LPF > Cutoff", "", preset.LPFCutoff),
				NewSelectorOption("LPF > Resonance", "", preset.LPFResonance),
			),
			NewSliderNode("Amount", preset.ModulationUpdateKind, preset.ModKeysSpacing*i+preset.ModParamAmt, -1000, 1000, .01, formatSemiTon),
			NewSelectorNode("Shape", preset.ModulationUpdateKind, preset.ModKeysSpacing*i+preset.ModParamShp,
				NewSelectorOption("Linear", "", preset.ModShapeLinear),
				NewSelectorOption("Exponential", "", preset.ModShapeExponential),
				NewSelectorOption("Logarithmic", "", preset.ModShapeLogarithmic),
			),
		)

		slotNode.AttachPreview(func() string {
			src := slotNode.QueryAll("Source")[0].(SelectorNode)
			dst := slotNode.QueryAll("Destination")[0].(SelectorNode)

			if dst.Val() != preset.ParamNone {
				slotNode.SetLabel(src.CurrentOption().Label()) // trick for left/right preview display
				return dst.CurrentOption().Label()
			}

			slotNode.SetLabel("NONE")
			return ""
		})

		matrix.Append(slotNode)
	}

	return matrix
}
