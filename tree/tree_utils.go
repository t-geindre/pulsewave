package tree

import (
	"synth/preset"
)

func waveFormNode(key uint8) Node {
	return NewSelectorNode("Waveform", preset.PresetUpdateKind, key,
		NewSelectorOption("Sine", "ui/icons/sine_wave", 0),
		NewSelectorOption("Square", "ui/icons/square_wave", 1),
		NewSelectorOption("Sawtooth", "ui/icons/saw_wave", 2),
		NewSelectorOption("Triangle", "ui/icons/triangle_wave", 3),
		NewSelectorOption("Noise", "ui/icons/noise_wave", 4),
	)
}

func adsrNode(label string, att, dec, sus, rel uint8, children ...Node) Node {
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

func adsrNodeWithToggle(label string, toggle, att, dec, sus, rel uint8, children ...Node) Node {
	node := adsrNode(label, att, dec, sus, rel, children...)
	node.Prepend(onOffNode(toggle))

	return node
}

func onOffNode(key uint8) Node {
	return NewSelectorNode("ON/OFF", preset.PresetUpdateKind, key,
		NewSelectorOption("OFF", "", 0),
		NewSelectorOption("ON", "", 1),
	)
}

func allPresetsNodes(presets []string) []Node {
	nodes := make([]Node, len(presets))
	for i, p := range presets {
		nodes[i] = NewValidatingSelectorNode(p, preset.LoadSavePresetKind, uint8(i),
			NewSelectorOption("Load", "", 0),
			NewSelectorOption("Save", "", 1),
		)
	}

	return nodes
}
