package tree

type AttachFunc func(uint8, float32)

type ValueNode interface {
	Node

	Val() float32
	SetVal(float32)
	Key() uint8
	Attach(AttachFunc)
}

type ParamNode struct {
	val     float32
	key     uint8
	publish AttachFunc

	Node
}

func NewValueNode(label string, key uint8) *ParamNode {
	return &ParamNode{
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

	p.val = f
	p.publish(p.key, f)
}

func (p *ParamNode) Key() uint8 {
	return p.key
}

func (p *ParamNode) Attach(publish AttachFunc) {
	p.publish = publish
}

type OptionValidateNode interface {
	Validate()
}
