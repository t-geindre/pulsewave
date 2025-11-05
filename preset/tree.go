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
				NewSelectorNode("Waveform", Osc0Shape,
					NewSelectorOption("Sine", "ui/icons/sine_wave", 0),
					NewSelectorOption("Square", "ui/icons/square_wave", 1),
					NewSelectorOption("Sawtooth", "ui/icons/saw_wave", 2),
					NewSelectorOption("Triangle", "ui/icons/triangle_wave", 3),
					NewSelectorOption("Noise", "ui/icons/noise_wave", 4),
				),
				NewSliderNode("Detune", Osc0Detune, -100, 100, .1, func(v float32) string {
					return fmt.Sprintf("%.1f st", v)
				}),
				NewSliderNode("Gain", Osc0Gain, 0, 1, .01, nil),
			),
			NewListNode("Oscillator 2",
				NewSelectorNode("Waveform", Osc1Shape,
					NewSelectorOption("Sine", "ui/icons/sine_wave", 0),
					NewSelectorOption("Square", "ui/icons/square_wave", 1),
					NewSelectorOption("Sawtooth", "ui/icons/saw_wave", 2),
					NewSelectorOption("Triangle", "ui/icons/triangle_wave", 3),
					NewSelectorOption("Noise", "ui/icons/noise_wave", 4),
				),
				NewSliderNode("Detune", Osc1Detune, -100, 100, .1, func(v float32) string {
					return fmt.Sprintf("%.1f st", v)
				}),
				NewSliderNode("Gain", Osc1Gain, 0, 1, .01, nil),
			),
			NewListNode("Oscillator 3",
				NewSelectorNode("Waveform", Osc2Shape,
					NewSelectorOption("Sine", "ui/icons/sine_wave", 0),
					NewSelectorOption("Square", "ui/icons/square_wave", 1),
					NewSelectorOption("Sawtooth", "ui/icons/saw_wave", 2),
					NewSelectorOption("Triangle", "ui/icons/triangle_wave", 3),
					NewSelectorOption("Noise", "ui/icons/noise_wave", 4),
				),
				NewSliderNode("Detune", Osc2Detune, -100, 100, .1, func(v float32) string {
					return fmt.Sprintf("%.1f st", v)
				}),
				NewSliderNode("Gain", Osc2Gain, 0, 1, .01, nil),
			),
		),
		NewListNode("Modulation",
			NewListNode("Amplitude"),
			NewListNode("Cutoff"),
			NewListNode("Resonance"),
			NewListNode("Pitch"),
		),
		NewListNode("Effects",
			NewListNode("Feedback delay",
				NewSliderNode("Delay", FBDelayParam, 0, 2, .001, func(v float32) string {
					return fmt.Sprintf("%.0f ms", v*1000)
				}),
				NewSliderNode("Feedback", FBFeedBack, 0, 0.95, .01, nil),
				NewSliderNode("Mix", FBMix, 0, 1, .01, nil),
				NewSliderNode("Tone", FBTone, 200, 8000, 1, func(v float32) string {
					return fmt.Sprintf("%.0f Hz", v)
				}),
			),
			NewListNode("Low pass filter"),
			NewListNode("Unison",
				NewSliderNode("Voices", UnisonVoices, 1, 16, 1, func(v float32) string {
					return fmt.Sprintf("%.0f voices", v)
				}),
				NewSliderNode("Pan spread", UnisonPanSpread, 0, 1, .01, nil),
				NewSliderNode("Phase spread", UnisonPhaseSpread, 0, 1, .01, func(v float32) string {
					return fmt.Sprintf("%.0f%% cycle", v*100)
				}),
				NewSliderNode("Detune spread", UnisonDetuneSpread, 0, 100, .1, func(v float32) string {
					return fmt.Sprintf("%.1f cent", v)
				}),
				NewSliderNode("Curve gamma", UnisonCurveGamma, 0.1, 2, .1, nil),
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
