package tree

import (
	"synth/dsp"
	"synth/msg"
	"synth/preset"

	"github.com/rs/zerolog"
)

type Tree struct {
	Node
	parameters map[uint8]ValueNode
}

func NewTree(logger zerolog.Logger) *Tree {
	t := &Tree{
		Node:       buildTree(logger),
		parameters: make(map[uint8]ValueNode),
	}

	t.SetRoot(t)

	return t
}

func (t *Tree) SetParam(key uint8, val float32) {
	if pn, ok := t.parameters[key]; ok {
		pn.SetVal(val)
	}
}

func (t *Tree) HandleMessage(msg msg.Message) {
	if msg.Kind == preset.ParamUpdateKind {
		t.SetParam(msg.Key, msg.ValF)
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

func (t *Tree) LoadPreset(p *preset.Preset) {
	for key, param := range p.Params {
		t.SetParam(key, param.GetBase())
	}
}

func (t *Tree) GetPreset() *preset.Preset {
	p := preset.NewPreset()

	for key, pn := range t.parameters {
		p.Params[key] = dsp.NewParam(pn.Val())
	}

	return p
}

func (t *Tree) Query(f func(n Node) bool) []Node {
	var result []Node
	t.queryNodes(t.Node, f, &result)
	return result
}

func (t *Tree) queryNodes(n Node, f func(n Node) bool, result *[]Node) {
	if f(n) {
		*result = append(*result, n)
	}

	for _, c := range n.Children() {
		t.queryNodes(c, f, result)
	}
}
