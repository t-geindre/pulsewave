package tree

import (
	"synth/preset"
)

func NewOscillatorNode(label string, shape, detune, gain, phase, pw uint8) Node {
	return NewNode(label,
		NewWaveFormNode(shape),
		NewSliderNode("Detune", preset.PresetUpdateKind, detune, -100, 100, .01, formatSemiTon),
		NewSliderNode("Gain", preset.PresetUpdateKind, gain, 0, 1, .01, nil),
		NewSliderNode("Phase", preset.PresetUpdateKind, phase, 0, 1, .01, formatCycle),
		NewSliderNode("Pulse width", preset.PresetUpdateKind, pw, 0.01, 0.5, .01, nil),
	)
}

func NewWaveFormNode(key uint8) Node {
	return NewSelectorNode("Waveform", preset.PresetUpdateKind, key,
		NewSelectorOption("Sine", "ui/icons/sine_wave", 0),
		NewSelectorOption("Square", "ui/icons/square_wave", 1),
		NewSelectorOption("Sawtooth", "ui/icons/saw_wave", 2),
		NewSelectorOption("Triangle", "ui/icons/triangle_wave", 3),
	)
}

func NewAdsrNode(label string, att, dec, sus, rel uint8, children ...Node) Node {
	n := NewNode(label,
		NewSliderNode("Attack", preset.PresetUpdateKind, att, 0, 10, .001, formatMillisecond),
		NewSliderNode("Decay", preset.PresetUpdateKind, dec, 0, 10, .001, formatMillisecond),
		NewSliderNode("Sustain", preset.PresetUpdateKind, sus, 0, 1, .01, nil),
		NewSliderNode("Release", preset.PresetUpdateKind, rel, 0, 10, .001, formatMillisecond),
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
	return NewSelectorNode("Status", preset.PresetUpdateKind, key,
		NewSelectorOption("OFF", "", 0),
		NewSelectorOption("ON", "", 1),
	)
}

func NewPresetsNodes(presets []string) []Node {
	nodes := make([]Node, len(presets))
	for i, p := range presets {
		nodes[i] = NewValidatingSelectorNode(p, preset.PresetLoadSaveKind, uint8(i),
			NewSelectorOption("Load", "", 0),
			NewSelectorOption("Save", "", 1),
		)
	}

	return nodes
}
