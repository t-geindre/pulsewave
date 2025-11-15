package tree

type Node interface {
	Label() string
	Children() []Node
	Parent() Node
	SetParent(Node)
	Append(Node)
	Prepend(Node)
	Remove(Node)
	SetRoot(Node)
	Root() Node
}

type node struct {
	label    string
	children []Node
	parent   Node
	root     Node
}

func NewNode(label string, children ...Node) Node {
	n := &node{
		label:    label,
		children: []Node{},
	}
	for _, c := range children {
		n.Append(c)
	}
	return n
}

func (n *node) Label() string {
	return n.label
}

func (n *node) Children() []Node {
	return n.children
}

func (n *node) Parent() Node {
	return n.parent
}

func (n *node) SetParent(p Node) {
	n.parent = p
}

func (n *node) Append(child Node) {
	child.SetParent(n)
	n.children = append(n.children, child)
}

func (n *node) Prepend(child Node) {
	child.SetParent(n)
	n.children = append([]Node{child}, n.children...)
}

func (n *node) SetRoot(r Node) {
	n.root = r
	for _, c := range n.children {
		c.SetRoot(r)
	}
}

func (n *node) Remove(child Node) {
	for i, c := range n.children {
		if c == child {
			n.children = append(n.children[:i], n.children[i+1:]...)
			child.SetParent(nil)
			return
		}
	}
}

func (n *node) Root() Node {
	return n.root
}
