package ui

import (
	"synth/midi"
	"synth/msg"
	"synth/preset"
)

type Messenger struct {
	tree      *preset.Tree
	ctrls     *MidiControls
	inQ, outQ *msg.Queue
}

func NewMessenger(tree *preset.Tree, ctrls *MidiControls, inQ, outQ *msg.Queue) *Messenger {
	m := &Messenger{
		tree:  tree,
		ctrls: ctrls,
		inQ:   inQ,
		outQ:  outQ,
	}

	tree.AttachUpdater(m.PublishParameterUpdate)

	return m
}

func (m *Messenger) Update() {
	m.inQ.Drain(10, func(msg msg.Message) {
		switch msg.Kind {
		case preset.ParamUpdateKind:
			m.tree.SetParam(msg.Key, msg.ValF)
		case midi.ControlChangeKind:
			m.ctrls.ControlChange(msg.Key, msg.Val8)
		}
	})
}

func (m *Messenger) PullAllParameters() {
	m.outQ.TryWrite(msg.Message{
		Kind: preset.ParamPullAllKind,
	})
}

func (m *Messenger) PublishParameterUpdate(key uint8, val float32) {
	m.outQ.TryWrite(msg.Message{
		Kind: preset.ParamUpdateKind,
		Key:  key,
		ValF: val,
	})
}
