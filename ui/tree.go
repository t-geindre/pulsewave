package ui

import (
	"fmt"
	"synth/msg"
	"synth/preset"
)

type Tree struct {
	Node
	pInQueue, pOutQueue *msg.Queue
	parameters          map[uint8]Node
}

func NewTree(pInQueue, pOutQueue *msg.Queue) *Tree {
	t := &Tree{
		Node:       buildTree(),
		pInQueue:   pInQueue,
		pOutQueue:  pOutQueue,
		parameters: make(map[uint8]Node),
	}

	t.AttachNodes(t.Node)
	t.PullAll()

	return t
}

func (t *Tree) PullAll() {
	t.pOutQueue.TryWrite(msg.Message{
		Source: preset.AudioSource,
		Kind:   preset.ParamPullAllKind,
	})
}

func (t *Tree) PublishUpdate(key uint8, val float32) {
	t.pOutQueue.TryWrite(msg.Message{
		Source: preset.AudioSource,
		Kind:   preset.ParamUpdateKind,
		Key:    key,
		ValF:   val,
	})
}

func (t *Tree) Update() {
	t.pInQueue.Drain(10, func(m msg.Message) {
		if m.Kind == preset.ParamUpdateKind {
			for key, node := range t.parameters {
				if key == m.Key {
					if pn, ok := node.(*ParameterNode); ok {
						pn.SetVal(m.ValF)
					}
					break
				}
			}
		}
	})
}

func (t *Tree) AttachNodes(n Node) {
	if pn, ok := n.(*ParameterNode); ok {
		t.parameters[pn.key] = pn
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
				NewListNode("Waveform"),
				NewListNode("Octave"),
				NewListNode("Semitone"),
				NewListNode("Detune"),
				NewListNode("Gain"),
			),
			NewListNode("Oscillator 2",
				NewListNode("Waveform"),
				NewListNode("Octave"),
				NewListNode("Semitone"),
				NewListNode("Detune"),
				NewListNode("Gain"),
			),
			NewListNode("Oscillator 3",
				NewListNode("Waveform"),
				NewListNode("Octave"),
				NewListNode("Semitone"),
				NewListNode("Detune"),
				NewListNode("Gain"),
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
				NewParameterNode("Delay", preset.FBDelayParam, 0, 2, .001, func(v float32) string {
					return fmt.Sprintf("%.0f ms", v*1000)
				}),
				NewParameterNode("Feedback", preset.FBFeedBack, 0, 0.95, .01, nil),
				NewParameterNode("Mix", preset.FBMix, 0, 1, .01, nil),
				NewParameterNode("Tone", preset.FBTone, 200, 8000, 1, func(v float32) string {
					return fmt.Sprintf("%.0f Hz", v)
				}),
			),
			NewListNode("Low pass filter"),
			NewListNode("Unison",
				NewParameterNode("Voices", preset.UnisonVoices, 1, 16, 1, nil),
				NewParameterNode("Pan spread", preset.UnisonPanSpread, 0, 1, .01, nil),
				NewParameterNode("Phase spread", preset.UnisonPhaseSpread, 0, 1, .01, func(v float32) string {
					return fmt.Sprintf("%.2f%% cycle", v*100)
				}),
				NewParameterNode("Detune spread", preset.UnisonDetuneSpread, 0, 100, .1, func(v float32) string {
					return fmt.Sprintf("%.2f cent", v)
				}),
				NewParameterNode("Curve gamma", preset.UnisonCurveGamma, 0.1, 2, .1, nil),
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
