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
			),
			NewListNode("Low pass filter"),
			NewListNode("Unison",
				NewListNode("Voices"),
				NewListNode("Pan spread"),
				NewListNode("Phase spread"),
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
