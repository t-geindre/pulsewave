package preset

import (
	"fmt"
)

type Tree struct {
	Node
	parameters map[uint8]ValueNode
}

func NewTree() *Tree {
	t := &Tree{
		Node:       buildTree(),
		parameters: make(map[uint8]ValueNode),
	}

	return t
}

func (t *Tree) SetParam(key uint8, val float32) {
	if pn, ok := t.parameters[key]; ok {
		pn.SetVal(val)
	}
}

func (t *Tree) AttachUpdater(publish func(key uint8, val float32)) {
	t.attachNodes(t.Node, publish)
}

func (t *Tree) attachNodes(n Node, f func(key uint8, val float32)) {
	if pn, ok := n.(ValueNode); ok {
		t.parameters[pn.Key()] = pn
		pn.Attach(f)
	}

	for _, c := range n.Children() {
		t.attachNodes(c, f)
	}
}

func buildTree() Node {
	return NewListNode("",
		NewListNode("Oscillators",
			NewListNode("Oscillator 1",
				waveFormNode(Osc0Shape),
				NewSliderNode("Detune", Osc0Detune, -100, 100, .1, formatSemiTon),
				NewSliderNode("Gain", Osc0Gain, 0, 1, .01, nil),
				NewSliderNode("Phase", Osc0Phase, 0, 1, .01, formatCycle),
				NewSliderNode("Pulse width", Osc0Pw, 0.01, 0.5, .01, nil),
			),
			NewListNode("Oscillator 2",
				waveFormNode(Osc1Shape),
				NewSliderNode("Detune", Osc1Detune, -100, 100, .1, formatSemiTon),
				NewSliderNode("Gain", Osc1Gain, 0, 1, .01, nil),
				NewSliderNode("Phase", Osc1Phase, 0, 1, .01, formatCycle),
				NewSliderNode("Pulse width", Osc1Pw, 0.01, 0.5, .01, nil),
			),
			NewListNode("Oscillator 3",
				waveFormNode(Osc2Shape),
				NewSliderNode("Detune", Osc2Detune, -100, 100, .1, formatSemiTon),
				NewSliderNode("Gain", Osc2Gain, 0, 1, .01, nil),
				NewSliderNode("Phase", Osc2Phase, 0, 1, .01, formatCycle),
				NewSliderNode("Pulse width", Osc2Pw, 0.01, 0.5, .01, nil),
			),
		),
		NewListNode("Modulation",
			adsrNode("Amplitude", AmpEnvAttack, AmpEnvDecay, AmpEnvSustain, AmpEnvRelease),
			NewListNode("Cutoff",
				NewListNode("LFO",
					onOffNode(LpfLfoOnOff),
					NewSliderNode("Amount", LpfLfoAmount, 20, 20000, 1, formatHertz),
					waveFormNode(LpfLfoShape),
					NewSliderNode("Rate", LpfLfoFreq, 0.01, 20, .01, formatLowHertz),
					NewSliderNode("Phase", LpfLfoPhase, 0, 1, .01, formatCycle),
				),
				adsrNodeWithToggle("ADSR", LpfAdsrOnOff, LpfAdsrAttack, LpfAdsrDecay, LpfAdsrSustain, LpfAdsrRelease,
					NewSliderNode("Amount", LpfAdsrAmount, -20000, 20000, 1, formatHertz),
				),
			),
			NewListNode("Resonance"),
			NewListNode("Pitch",
				NewListNode("LFO",
					onOffNode(PitchLfoOnOff),
					NewSliderNode("Amount", PitchLfoAmount, -1000, 1000, 1, formatSemiTon),
					waveFormNode(PitchLfoShape),
					NewSliderNode("Rate", PitchLfoFreq, 0.01, 20, .01, formatLowHertz),
					NewSliderNode("Phase", PitchLfoPhase, 0, 1, .01, formatCycle),
				),
				adsrNodeWithToggle("ADSR", PitchAdsrOnOff, PitchAdsrAttack, PitchAdsrDecay, PitchAdsrSustain, PitchAdsrRelease,
					NewSliderNode("Amount", PitchAdsrAmount, -1000, 1000, 1, formatSemiTon),
				),
			),
		),
		NewListNode("Effects",
			NewListNode("Feedback delay",
				onOffNode(FBOnOff),
				NewSliderNode("Delay", FBDelayParam, 0, 2, .001, formatMillisecond),
				NewSliderNode("Feedback", FBFeedBack, 0, 0.95, .01, nil),
				NewSliderNode("Mix", FBMix, 0, 1, .01, nil),
				NewSliderNode("Tone", FBTone, 200, 8000, 1, formatHertz),
			),
			NewListNode("Low pass filter",
				onOffNode(LPFOnOff),
				NewSliderNode("Cutoff", LPFCutoff, 20, 20000, 1, formatHertz),
				NewSliderNode("Resonance", LPFResonance, 0.1, 10, .01, nil),
			),
			NewListNode("Unison",
				onOffNode(UnisonOnOff),
				NewSliderNode("Voices", UnisonVoices, 1, 16, 1, func(v float32) string {
					return fmt.Sprintf("%.0f voices", v)
				}),
				NewSliderNode("Pan spread", UnisonPanSpread, 0, 1, .01, nil),
				NewSliderNode("Phase spread", UnisonPhaseSpread, 0, 1, .01, formatCycle),
				NewSliderNode("Detune spread", UnisonDetuneSpread, 0, 100, .1, formatCent),
				NewSliderNode("Curve gamma", UnisonCurveGamma, 0.1, 4, .1, nil),
			),
		),
		NewListNode("Visualizer",
			NewListNode("Spectrum"),
			NewListNode("Oscilloscope"),
		),
		NewListNode("Presets",
			NewListNode("Load preset"),
			NewListNode("Save preset"),
			NewListNode("Auto save"),
		),
		NewListNode("Settings",
			NewListNode("General"),
			NewListNode("Master gain"),
			NewListNode("MIDI controller"),
			NewListNode("Pitch bend"),
			NewListNode("About"),
		),
	)
}

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
