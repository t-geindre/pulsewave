package tree

import "synth/msg"

type AttachFunc func(msg.Kind, uint8, float32)
type PreviewFunc func() (string, string) // Preview, override label

type Node interface {
	Label() string
	SetLabel(string)

	Children() []Node
	Append(Node)
	Prepend(Node)

	Root() Node
	Parent() Node
	SetParent(Node)

	AttachMessenger(*msg.Messenger)
	AttachPreview(PreviewFunc)

	Preview() (string, string) // Preview, override label
	// QueryAll returns all child nodes with the given label in the subtree
	QueryAll(names ...string) []Node

	SetContext(*Context)
	Context() *Context

	Focus()
	Blur()
}

type node struct {
	label    string
	children []Node
	parent   Node
	context  *Context

	preview   PreviewFunc
	prevCache [2]string
	prevDirty bool
}

func NewNode(label string, children ...Node) Node {
	n := &node{
		label:     label,
		children:  []Node{},
		context:   NewContext(),
		prevDirty: true,
	}
	for _, c := range children {
		n.Append(c)
	}
	return n
}

func (n *node) Label() string {
	return n.label
}

func (n *node) SetLabel(l string) {
	n.label = l
}

func (n *node) Children() []Node {
	return n.children
}

func (n *node) Parent() Node {
	return n.parent
}

func (n *node) SetParent(p Node) {
	n.parent = p
	n.SetContext(p.Context())
}

func (n *node) Append(child Node) {
	child.SetParent(n)
	n.children = append(n.children, child)
}

func (n *node) Prepend(child Node) {
	child.SetParent(n)
	n.children = append([]Node{child}, n.children...)
}

func (n *node) AttachMessenger(m *msg.Messenger) {
	for _, c := range n.children {
		c.AttachMessenger(m)
	}
}

func (n *node) QueryAll(names ...string) []Node {
	var results []Node
	for _, name := range names {

		for _, c := range n.children {
			if c.Label() == name {
				results = append(results, c)
			}
			results = append(results, c.QueryAll(name)...)
		}
	}

	return results
}

func (n *node) AttachPreview(p PreviewFunc) {
	n.preview = p
}

func (n *node) Preview() (string, string) {
	if n.prevDirty && n.preview != nil {
		n.prevCache[0], n.prevCache[1] = n.preview()
	}
	return n.prevCache[0], n.prevCache[1]
}

func (n *node) SetContext(c *Context) {
	n.context = c
	for _, child := range n.children {
		child.SetContext(c)
	}
}

func (n *node) Context() *Context {
	return n.context
}

func (n *node) Root() Node {
	root := Node(n)
	for root.Parent() != nil {
		root = root.Parent()
	}
	return root
}

func (n *node) Focus() {
	// No-op
}

func (n *node) Blur() {
	n.prevDirty = true // Invalidate preview cache
}
