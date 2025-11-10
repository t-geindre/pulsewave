package preset

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
