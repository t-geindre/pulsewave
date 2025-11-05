package preset

type Node interface {
	Label() string
	Children() []Node
	Parent() Node
	SetParent(Node)
	Append(Node)
}

type AttachFunc func(uint8, float32)

type ValueNode interface {
	Val() float32
	SetVal(float32)
	Key() uint8
	Attach(AttachFunc)
}

type ParamNode struct {
	val     float32
	key     uint8
	publish AttachFunc
}

func NewParamNode(key uint8) *ParamNode {
	return &ParamNode{
		key: key,
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
