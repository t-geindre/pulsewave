package tree

import "synth/msg"

type ValueNode interface {
	Node

	Val() float32
	SetVal(float32)
	SetValAndPublish(f float32) // Force publish
	Key() uint8
}

type ParamNode struct {
	kind      msg.Kind
	key       uint8
	val       float32
	messenger *msg.Messenger

	Node
}

func NewValueNode(label string, kind msg.Kind, key uint8) *ParamNode {
	return &ParamNode{
		kind: kind,
		key:  key,
		Node: NewNode(label),
	}
}

func (p *ParamNode) Val() float32 {
	return p.val
}

func (p *ParamNode) SetVal(f float32) {
	if p.val == f {
		return
	}

	p.SetValAndPublish(f)
}

func (p *ParamNode) SetValAndPublish(f float32) {
	p.val = f
	p.messenger.SendMessage(msg.Message{
		Kind: p.kind,
		Key:  p.key,
		ValF: p.val,
	})
}

func (p *ParamNode) HandleMessage(msg msg.Message) {
	if p.kind == msg.Kind && p.key == msg.Key {
		p.SetVal(msg.ValF)
	}
}

func (p *ParamNode) Key() uint8 {
	return p.key
}

func (p *ParamNode) AttachMessenger(m *msg.Messenger) {
	p.messenger = m
	m.RegisterHandler(p)
}

type OptionValidateNode interface {
	Validate()
}
