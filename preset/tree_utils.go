package preset

import (
	"os"
	"path"
	"slices"
	"strings"

	"github.com/rs/zerolog"
)

func waveFormNode(key uint8) Node {
	return NewSelectorNode("Waveform", key,
		NewSelectorOption("Sine", "ui/icons/sine_wave", 0),
		NewSelectorOption("Square", "ui/icons/square_wave", 1),
		NewSelectorOption("Sawtooth", "ui/icons/saw_wave", 2),
		NewSelectorOption("Triangle", "ui/icons/triangle_wave", 3),
		NewSelectorOption("Noise", "ui/icons/noise_wave", 4),
	)
}

func adsrNode(label string, att, dec, sus, rel uint8, children ...Node) Node {
	n := NewListNode(label,
		NewSliderNode("Attack", att, 0, 10, .001, formatMillisecond),
		NewSliderNode("Decay", dec, 0, 10, .001, formatMillisecond),
		NewSliderNode("Sustain", sus, 0, 1, .01, nil),
		NewSliderNode("Release", rel, 0, 10, .001, formatMillisecond),
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
	return NewSelectorNode("ON/OFF", key,
		NewSelectorOption("OFF", "", 0),
		NewSelectorOption("ON", "", 1),
	)
}

func allPresetsNodes(pth string, logger zerolog.Logger) []Node {
	files, err := os.ReadDir(pth)
	if err != nil {
		logger.Error().Err(err).Str("path", pth).Msg("failed to read presets directory")
	}
	var presets []Node
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		presets = append(presets, NewPresetNode(path.Join(pth, f.Name()), logger))
	}

	slices.SortFunc(presets, func(a, b Node) int {
		return strings.Compare(a.Label(), b.Label())
	})

	return presets
}
