package preset

type ListNode struct {
	label    string
	children []Node
	parent   Node
}

func NewListNode(label string, children ...Node) *ListNode {
	n := &ListNode{
		label:    label,
		children: []Node{},
	}
	for _, c := range children {
		n.Append(c)
	}
	return n
}

func (n *ListNode) Label() string {
	return n.label
}

func (n *ListNode) Children() []Node {
	return n.children
}

func (n *ListNode) Parent() Node {
	return n.parent
}

func (n *ListNode) SetParent(p Node) {
	n.parent = p
}

func (n *ListNode) Append(child Node) {
	child.SetParent(n)
	n.children = append(n.children, child)
}

func (n *ListNode) Prepend(child Node) {
	child.SetParent(n)
	n.children = append([]Node{child}, n.children...)
}
