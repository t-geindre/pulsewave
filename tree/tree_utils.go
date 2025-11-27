package tree

import (
	"fmt"
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

func NewAdsrNodeWithToggle(label string, toggle, att, dec, sus, rel uint8, children ...Node) Node {
	node := NewAdsrNode(label, att, dec, sus, rel, children...)
	node.Prepend(NewOnOffNode(toggle))

	return node
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

func NewLfoNode(label string, onOff, shape, rate, phase, amount uint8, min, max, step float32, f FormatFunc) Node {
	return NewNode(label,
		NewOnOffNode(onOff),
		NewSliderNode("Amount", preset.UpdateParameterKind, amount, min, max, step, f),
		NewWaveFormNode(shape),
		NewSliderNode("Rate", preset.UpdateParameterKind, rate, 0.01, 20, .01, formatLowHertz),
		NewSliderNode("Phase", preset.UpdateParameterKind, phase, 0, 1, .01, formatCycle),
	)
}

func NewModulationMatrixNode(label string) Node {
	matrix := NewNode(label)

	for i := uint8(0); i < preset.ModSlots; i++ {
		matrix.Append(
			NewNode(fmt.Sprintf("Slot %02d", i+1),
				NewSelectorNode("Source", preset.ModulationUpdateKind, preset.ModKeysSpacing*i+preset.ModParamSrc,
					NewSelectorOption("Velocity", "", preset.ModSrcVelocity),
					NewSelectorOption("LFO 1", "", preset.ModSrcLfo1),
					NewSelectorOption("LFO 2", "", preset.ModSrcLfo2),
					NewSelectorOption("LFO 3", "", preset.ModSrcLfo3),
					NewSelectorOption("ADSR 1", "", preset.ModSrcAdsr1),
					NewSelectorOption("ADSR 2", "", preset.ModSrcAdsr2),
					NewSelectorOption("ADSR 3", "", preset.ModSrcAdsr3),
				),
				NewSelectorNode("Destination", preset.ModulationUpdateKind, preset.ModKeysSpacing*i+preset.ModParamDst,
					NewSelectorOption("Test", "", preset.None),
				),
				NewSliderNode("Amount", preset.ModulationUpdateKind, preset.ModKeysSpacing*i+preset.ModParamAmt, -1000, 1000, .01, formatSemiTon),
				NewSelectorNode("Shape", preset.ModulationUpdateKind, preset.ModKeysSpacing*i+preset.ModParamShp,
					NewSelectorOption("Linear", "", preset.ModShapeLinear),
					NewSelectorOption("Exponential", "", preset.ModShapeExponential),
					NewSelectorOption("Logarithmic", "", preset.ModShapeLogarithmic),
				),
			),
		)
	}

	return matrix

}
