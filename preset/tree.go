package preset

import (
	"fmt"
	"synth/msg"
)

type Tree struct {
	Node
	pInQueue, pOutQueue *msg.Queue
	parameters          map[uint8]ValueNode
}

func NewTree(pInQueue, pOutQueue *msg.Queue) *Tree {
	t := &Tree{
		Node:       buildTree(),
		pInQueue:   pInQueue,
		pOutQueue:  pOutQueue,
		parameters: make(map[uint8]ValueNode),
	}

	t.AttachNodes(t.Node)
	t.PullAll()

	return t
}

func (t *Tree) PullAll() {
	t.pOutQueue.TryWrite(msg.Message{
		Source: AudioSource,
		Kind:   ParamPullAllKind,
	})
}

func (t *Tree) PublishUpdate(key uint8, val float32) {
	t.pOutQueue.TryWrite(msg.Message{
		Source: AudioSource,
		Kind:   ParamUpdateKind,
		Key:    key,
		ValF:   val,
	})
}

func (t *Tree) Update() {
	t.pInQueue.Drain(10, func(m msg.Message) {
		if m.Kind == ParamUpdateKind {
			for key, node := range t.parameters {
				if key == m.Key {
					node.SetVal(m.ValF)
					break
				}
			}
		}
	})
}

func (t *Tree) AttachNodes(n Node) {
	if pn, ok := n.(ValueNode); ok {
		t.parameters[pn.Key()] = pn
		pn.Attach(t.PublishUpdate)
	}

	for _, c := range n.Children() {
		t.AttachNodes(c)
	}
}

func buildTree() Node {
	return NewListNode("",
		NewListNode("Oscillators",
			NewListNode("Oscillator 1",
				waveFormNode(Osc0Shape),
				NewSliderNode("Detune", Osc0Detune, -100, 100, .1, formatSemiTon),
				NewSliderNode("Gain", Osc0Gain, 0, 1, .01, nil),
			),
			NewListNode("Oscillator 2",
				waveFormNode(Osc2Shape),
				NewSliderNode("Detune", Osc1Detune, -100, 100, .1, formatSemiTon),
				NewSliderNode("Gain", Osc1Gain, 0, 1, .01, nil),
			),
			NewListNode("Oscillator 3",
				waveFormNode(Osc2Shape),
				NewSliderNode("Detune", Osc2Detune, -100, 100, .1, formatSemiTon),
				NewSliderNode("Gain", Osc2Gain, 0, 1, .01, nil),
			),
		),
		NewListNode("Modulation",
			adsrNode("Amplitude", AmpEnvAttack, AmpEnvDecay, AmpEnvSustain, AmpEnvRelease),
			NewListNode("Cutoff",
				NewListNode("LFO",
					NewSliderNode("ON/OFF", LpfLfoOnOff, 0, 1, 1, formatOnOff),
					NewSliderNode("Amount", LpfLfoAmount, 20, 20000, 1, formatHertz),
					waveFormNode(LpfLfoShape),
					NewSliderNode("Rate", LpfLfoFreq, 0.01, 20, .01, formatLowHertz),
					NewSliderNode("Phase", LpfLfoPhase, 0, 1, .01, formatCycle),
				),
				adsrNodeWithToggle("ADSR", LpfAdsrOnOff, LpfAdsrAttack, LpfAdsrDecay, LpfAdsrSustain, LpfAdsrRelease,
					NewSliderNode("Amount", LpfAdsrAmount, 20, 20000, 1, formatHertz),
				),
			),
			NewListNode("Resonance"),
			NewListNode("Pitch"),
		),
		NewListNode("Effects",
			NewListNode("Feedback delay",
				NewSliderNode("ON/OFF", FBOnOff, 0, 1, 1, formatOnOff),
				NewSliderNode("Delay", FBDelayParam, 0, 2, .001, formatMillisecond),
				NewSliderNode("Feedback", FBFeedBack, 0, 0.95, .01, nil),
				NewSliderNode("Mix", FBMix, 0, 1, .01, nil),
				NewSliderNode("Tone", FBTone, 200, 8000, 1, formatHertz),
			),
			NewListNode("Low pass filter",
				NewSliderNode("ON/OFF", LPFOnOff, 0, 1, 1, formatOnOff),
				NewSliderNode("Cutoff", LPFCutoff, 20, 20000, 1, formatHertz),
				NewSliderNode("Resonance", LPFResonance, 0.1, 10, .01, nil),
			),
			NewListNode("Unison",
				NewSliderNode("ON/OFF", UnisonOnOff, 0, 1, 1, formatOnOff),
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
			// Todo ON/OFF, add toggle node
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
	node.Prepend(NewSliderNode("ON/OFF", toggle, 0, 1, 1, formatOnOff))

	return node
}
